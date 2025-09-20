package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	hydra "github.com/ory/hydra-client-go/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	hydraAdminURL    string
	backendURL       string
	googleProvider   *OAuthProvider
	oauthProviders   = map[string]*OAuthProvider{}
	hydraAdminClient *hydra.APIClient
)

func init() {
	godotenv.Load()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	hydraAdminURL = os.Getenv("HYDRA_ADMIN_URL")
	if hydraAdminURL == "" {
		log.Fatal("HYDRA_ADMIN_URL not set")
	}
	log.Printf("Hydra Admin URL: %s", hydraAdminURL)

	backendURL = os.Getenv("BACKEND_URL")
	if backendURL == "" {
		log.Fatal("BACKEND_URL not set")
	}
	log.Printf("Backend URL: %s", backendURL)

	googleProvider = NewOAuthProvider(google.Endpoint, "google", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("LOGIN_HANDLER_REDIRECT_URI"))
	oauthProviders["google"] = googleProvider

	config := hydra.NewConfiguration()
	config.Servers = hydra.ServerConfigurations{{URL: hydraAdminURL}}
	hydraAdminClient = hydra.NewAPIClient(config)
}

type OAuthProvider struct {
	OAuth2Config *oauth2.Config
	Name         string
}

func NewOAuthProvider(endpoint oauth2.Endpoint, name, clientID, clientSecret, redirectURL string) *OAuthProvider {
	return &OAuthProvider{
		OAuth2Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     endpoint,
		},
		Name: name,
	}
}

func generateStateOauthCookie() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateCodeVerifier generates a cryptographically random code_verifier
// according to RFC 7636 specification
func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Use base64url encoding without padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge generates a code_challenge from a code_verifier
// using SHA256 method as specified in RFC 7636
func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func main() {
	app := fiber.New(fiber.Config{
		StrictRouting:  true,
		ServerHeader:   "Fiber",
		ReadBufferSize: 16384,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	})

	app.Get("/login", handleLogin)
	app.Get("/callback", handleCallback)
	app.Get("/consent", handleConsent)
	log.Fatal(app.Listen(":8003"))
}

func handleLogin(c *fiber.Ctx) error {
	log.Printf("Handling login request")
	provider := c.Query("provider", "google")
	oauthProvider, exists := oauthProviders[provider]
	if !exists {
		return c.Status(http.StatusBadRequest).SendString("Unsupported provider")
	}

	loginChallenge := c.Query("login_challenge")
	if loginChallenge == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing login_challenge")
	}
	log.Printf("Calling Hydra Admin for login challenge: %s", loginChallenge)
	// Use request context for external Hydra call with a short timeout
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()
	loginRequest, _, err := hydraAdminClient.OAuth2API.GetOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).Execute()
	if err != nil {
		log.Printf("Failed to get login request: %v", err)
		return c.Status(http.StatusInternalServerError).SendString("Failed to get login request")
	}
	// Extract shop_id and prompt from request_url
	parsedURL, err := url.Parse(loginRequest.RequestUrl)
	if err != nil {
		log.Printf("Invalid request_url: %v", err)
		return c.Status(http.StatusInternalServerError).SendString("Invalid request_url")
	}
	query := parsedURL.Query()
	shopID := query.Get("shop_id")
	log.Printf("shop_id extracted from request_url: %s", shopID)

	prompt := query.Get("prompt")
	log.Printf("prompt extracted from request_url: %s", prompt)

	appType := query.Get("app_type")
	log.Printf("app_type extracted from request_url: %s", appType)
	if appType == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing app_type")
	}

	// Only require shop_id for storefront app type
	if appType == "storefront" && shopID == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing shop_id for storefront app")
	}

	// Only set shop_id cookie if appType is "storefront"
	if appType == "storefront" {
		c.Cookie(&fiber.Cookie{
			Name:     "shop_id",
			Value:    shopID,
			Expires:  time.Now().Add(10 * time.Minute),
			Secure:   true,
			SameSite: "Lax",
			HTTPOnly: true,
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "login_challenge",
		Value:    loginChallenge,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "Lax",
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "app_type",
		Value:    appType,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "Lax",
		HTTPOnly: true,
	})

	state, err := generateStateOauthCookie()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to generate state")
	}

	log.Printf("Generated OAuth state: %s", state)

	// Generate PKCE parameters
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to generate code verifier")
	}
	
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Store the code_verifier in a secure cookie for later use in callback
	c.Cookie(&fiber.Cookie{
		Name:     "code_verifier",
		Value:    codeVerifier,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})

	// Store the OAuth2 state in a secure cookie for CSRF protection
	c.Cookie(&fiber.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None", // Changed from "Lax" to "None"
		HTTPOnly: true,
	})

	// Create OAuth URL with PKCE parameters
	authURL := oauthProvider.OAuth2Config.AuthCodeURL(state, 
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	authURL += "&shop_id=" + url.QueryEscape(shopID)
	if prompt != "" {
		authURL += "&prompt=" + url.QueryEscape(prompt)
	}

	log.Printf("Redirecting to auth URL: %s", authURL)

	return c.Redirect(authURL)
}

func handleCallback(c *fiber.Ctx) error {
	state := c.Cookies("oauthstate")
	if state != c.Query("state") {
		return c.Status(http.StatusUnauthorized).SendString("Invalid OAuth state")
	}

	// Retrieve code_verifier from cookie
	codeVerifier := c.Cookies("code_verifier")
	if codeVerifier == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing code verifier")
	}

	// Retrieve app_type and shop_id from cookies
	appType := c.Cookies("app_type")
	if appType == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing app_type cookie")
	}
	shopID := c.Cookies("shop_id")
	log.Printf("shop_id from cookie: %s", shopID)

	// Only require shop_id cookie for storefront app type
	if appType == "storefront" && shopID == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing shop_id cookie for storefront app")
	}
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing authorization code")
	}

	provider := c.Query("provider", "google")
	oauthProvider, exists := oauthProviders[provider]
	if !exists {
		return c.Status(http.StatusBadRequest).SendString("Unsupported provider")
	}

	// Use request context with timeout for token exchange and subsequent userinfo calls
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	// Exchange code for token with PKCE
	token, err := oauthProvider.OAuth2Config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Token exchange failed: " + err.Error())
	}

	client := oauthProvider.OAuth2Config.Client(ctx, token)
	userInfo, err := fetchUserInfo(client)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to fetch user info: " + err.Error())
	}

	// Pre-flight: only register if the user does not already exist
	exists, err = userExists(client, userInfo.Email, appType, shopID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to check user existence: " + err.Error())
	}
	if !exists {
		if err := registerUser(client, userInfo, appType, shopID); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Failed to register user: " + err.Error())
		}
	}

	loginChallenge := c.Cookies("login_challenge")
	if loginChallenge == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing login_challenge")
	}
	redirectTo, err := acceptHydraLogin(ctx, loginChallenge, userInfo, shopID, appType)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to accept login in Hydra: " + err.Error())
	}

	return c.Redirect(redirectTo.RedirectTo)
}

func fetchUserInfo(client *http.Client) (*GoogleUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func userExists(client *http.Client, email, appType, shopID string) (bool, error) {
	var urlStr string
	if appType == "dashboard" {
		// Call the admin endpoint
		urlStr = fmt.Sprintf("%s/v1/userinfo?email=%s", backendURL, url.QueryEscape(email))
	} else if appType == "storefront" {
		// Call the customer endpoint only if appType is "storefront"
		urlStr = fmt.Sprintf("%s/v1/customerinfo?subdomain=%s&email=%s", backendURL, url.QueryEscape(shopID), url.QueryEscape(email))
	} else {
		return false, fmt.Errorf("unsupported app_type: %s", appType)
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status %d from user existence check", resp.StatusCode)
	}
}

func registerUser(client *http.Client, userInfo *GoogleUserInfo, appType, shopID string) error {
	var urlStr string
	var payload map[string]interface{}

	if appType == "storefront" {
		// Call the register customer endpoint
		urlStr = fmt.Sprintf("%s/v1/auth/register-customer", backendURL)
		// Convert shopID to int64 for customer registration
		shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid shop_id: %s", shopID)
		}
		payload = map[string]interface{}{
			"shop_id":          shopIDInt,
			"email":            userInfo.Email,
			"name":             userInfo.Name,
			"picture":          userInfo.Picture,
			"locale":           userInfo.Locale,
			"verified_email":   userInfo.VerifiedEmail,
			"auth_provider":    "google",
			"auth_provider_id": userInfo.ID,
		}
	} else {
		// Call the register user endpoint for admin
		urlStr = fmt.Sprintf("%s/v1/auth/register", backendURL)
		payload = map[string]interface{}{
			"email":          userInfo.Email,
			"name":           userInfo.Name,
			"id":             userInfo.ID,
			"sub":            userInfo.Sub,
			"provider":       "google",
			"picture":        userInfo.Picture,
			"locale":         userInfo.Locale,
			"verified_email": userInfo.VerifiedEmail,
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := client.Post(urlStr, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register user: status %d", resp.StatusCode)
	}
	return nil
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func acceptHydraLogin(ctx context.Context, loginChallenge string, user *GoogleUserInfo, shopID, appType string) (*hydra.OAuth2RedirectTo, error) {
	// Create a new map to hold the context data
	contextData := make(map[string]interface{})

	// Only include shop_id if appType is "storefront"
	if appType == "storefront" {
		contextData["shop_id"] = shopID
	}

	// Accept the login request
	acceptRequest := hydra.AcceptOAuth2LoginRequest{
		Subject:     *hydra.PtrString(user.Email),
		Remember:    hydra.PtrBool(true),
		RememberFor: hydra.PtrInt64(3600),
		Context:     contextData, // Assign the map here
	}
	// Note: ctx should be a request context provided by caller with timeout/cancellation.
	redirectTo, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2LoginRequest(ctx).
		LoginChallenge(loginChallenge).
		AcceptOAuth2LoginRequest(acceptRequest).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to accept Hydra login request: %w", err)
	}

	return redirectTo, nil
}

func handleConsent(c *fiber.Ctx) error {
	consentChallenge := c.Query("consent_challenge")
	if consentChallenge == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing consent_challenge")
	}
	// Use request context with timeout for calls to Hydra
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	consentRequest, _, err := hydraAdminClient.OAuth2API.GetOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).Execute()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get consent request in Hydra: " + err.Error())
	}
	redirectTo, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2ConsentRequest(ctx).
		ConsentChallenge(consentChallenge).
		AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
			GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
			GrantScope:               consentRequest.RequestedScope,
			Remember:                 hydra.PtrBool(true),
			RememberFor:              hydra.PtrInt64(3600),
		}).Execute()
	if err != nil {
		return fmt.Errorf("failed to accept Hydra consent request: %w", err)
	}
	return c.Redirect(redirectTo.RedirectTo)
}
