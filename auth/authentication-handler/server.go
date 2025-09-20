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
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	hydraAdminURL    string
	backendURL       string
	googleProvider   *OAuthProvider
	oauthProviders   = map[string]*OAuthProvider{}
	hydraAdminClient *hydra.APIClient
	logger           *zap.Logger
)

func init() {
	// Initialize zap logger
	var err error
	if os.Getenv("ENVIRONMENT") == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("Failed to initialize zap logger:", err)
	}
	defer logger.Sync()

	godotenv.Load()
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	hydraAdminURL = os.Getenv("HYDRA_ADMIN_URL")
	if hydraAdminURL == "" {
		logger.Fatal("HYDRA_ADMIN_URL not set")
	}
	logger.Info("Configuration loaded", zap.String("hydra_admin_url", hydraAdminURL))

	backendURL = os.Getenv("BACKEND_URL")
	if backendURL == "" {
		logger.Fatal("BACKEND_URL not set")
	}
	logger.Info("Configuration loaded", zap.String("backend_url", backendURL))

	googleProvider = NewOAuthProvider(google.Endpoint, "google", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("LOGIN_HANDLER_REDIRECT_URI"))
	oauthProviders["google"] = googleProvider

	config := hydra.NewConfiguration()
	config.Servers = hydra.ServerConfigurations{{URL: hydraAdminURL}}
	hydraAdminClient = hydra.NewAPIClient(config)

	logger.Info("Authentication handler initialized successfully",
		zap.String("provider", "google"),
		zap.String("client_id", os.Getenv("GOOGLE_CLIENT_ID")),
		zap.String("redirect_uri", os.Getenv("LOGIN_HANDLER_REDIRECT_URI")),
	)
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
	logger.Info("Starting authentication handler server")

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

	logger.Info("Routes registered successfully")
	logger.Info("Server listening on port 8003")

	err := app.Listen(":8003")
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func handleLogin(c *fiber.Ctx) error {
	logger.Info("Handling login request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("user_agent", c.Get("User-Agent")),
		zap.String("remote_ip", c.IP()),
	)

	provider := c.Query("provider", "google")
	oauthProvider, exists := oauthProviders[provider]
	if !exists {
		logger.Error("Unsupported OAuth provider", zap.String("provider", provider))
		return c.Status(http.StatusBadRequest).SendString("Unsupported provider")
	}

	loginChallenge := c.Query("login_challenge")
	if loginChallenge == "" {
		logger.Error("Missing login_challenge parameter")
		return c.Status(http.StatusBadRequest).SendString("Missing login_challenge")
	}

	logger.Info("Processing login challenge",
		zap.String("login_challenge", loginChallenge),
		zap.String("provider", provider),
	)

	// Use request context for external Hydra call with a short timeout
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	logger.Debug("Calling Hydra Admin API for login challenge",
		zap.String("hydra_admin_url", hydraAdminURL),
		zap.String("login_challenge", loginChallenge),
	)

	loginRequest, _, err := hydraAdminClient.OAuth2API.GetOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).Execute()
	if err != nil {
		logger.Error("Failed to get login request from Hydra",
			zap.Error(err),
			zap.String("login_challenge", loginChallenge),
		)
		return c.Status(http.StatusInternalServerError).SendString("Failed to get login request")
	}

	logger.Info("Successfully retrieved login request from Hydra",
		zap.String("subject", loginRequest.GetSubject()),
		zap.String("request_url", loginRequest.GetRequestUrl()),
		zap.Bool("skip", loginRequest.GetSkip()),
	)

	// Extract shop_id and prompt from request_url
	parsedURL, err := url.Parse(loginRequest.RequestUrl)
	if err != nil {
		logger.Error("Invalid request_url from Hydra",
			zap.Error(err),
			zap.String("request_url", loginRequest.RequestUrl),
		)
		return c.Status(http.StatusInternalServerError).SendString("Invalid request_url")
	}

	query := parsedURL.Query()
	shopID := query.Get("shop_id")
	prompt := query.Get("prompt")
	appType := query.Get("app_type")

	logger.Info("Extracted parameters from request_url",
		zap.String("shop_id", shopID),
		zap.String("prompt", prompt),
		zap.String("app_type", appType),
		zap.String("parsed_url", parsedURL.String()),
	)

	if appType == "" {
		logger.Error("Missing app_type parameter in request_url")
		return c.Status(http.StatusBadRequest).SendString("Missing app_type")
	}

	// Only require shop_id for storefront app type
	if appType == "storefront" && shopID == "" {
		logger.Error("Missing shop_id for storefront app", zap.String("app_type", appType))
		return c.Status(http.StatusBadRequest).SendString("Missing shop_id for storefront app")
	}

	// Only set shop_id cookie if appType is "storefront"
	if appType == "storefront" {
		logger.Debug("Setting shop_id cookie for storefront app", zap.String("shop_id", shopID))
		c.Cookie(&fiber.Cookie{
			Name:     "shop_id",
			Value:    shopID,
			Expires:  time.Now().Add(10 * time.Minute),
			Secure:   true,
			SameSite: "None", // Changed to None for cross-site OAuth compatibility
			HTTPOnly: true,
		})
	}

	logger.Debug("Setting authentication cookies",
		zap.String("login_challenge", loginChallenge),
		zap.String("app_type", appType),
	)

	c.Cookie(&fiber.Cookie{
		Name:     "login_challenge",
		Value:    loginChallenge,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None", // Changed to None for cross-site OAuth compatibility
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "app_type",
		Value:    appType,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None", // Changed to None for cross-site OAuth compatibility
		HTTPOnly: true,
	})

	state, err := generateStateOauthCookie()
	if err != nil {
		logger.Error("Failed to generate OAuth state", zap.Error(err))
		return c.Status(http.StatusInternalServerError).SendString("Failed to generate state")
	}

	logger.Info("Generated OAuth state", zap.String("state", state))

	// Generate PKCE parameters
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		logger.Error("Failed to generate PKCE code verifier", zap.Error(err))
		return c.Status(http.StatusInternalServerError).SendString("Failed to generate code verifier")
	}

	codeChallenge := generateCodeChallenge(codeVerifier)

	logger.Info("Generated PKCE parameters",
		zap.String("code_verifier_length", fmt.Sprintf("%d", len(codeVerifier))),
		zap.String("code_challenge", codeChallenge),
		zap.String("code_challenge_method", "S256"),
	)

	// Store the code_verifier in a secure cookie for later use in callback
	logger.Debug("Setting PKCE code_verifier cookie")
	c.Cookie(&fiber.Cookie{
		Name:     "code_verifier",
		Value:    codeVerifier,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})

	// Store the OAuth2 state in a secure cookie for CSRF protection
	logger.Debug("Setting OAuth state cookie")
	c.Cookie(&fiber.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None", // Changed from "Lax" to "None"
		HTTPOnly: true,
	})

	// Create OAuth URL with PKCE parameters
	logger.Debug("Building OAuth authorization URL with PKCE parameters")
	authURL := oauthProvider.OAuth2Config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	authURL += "&shop_id=" + url.QueryEscape(shopID)
	if prompt != "" {
		authURL += "&prompt=" + url.QueryEscape(prompt)
	}

	logger.Info("Redirecting to OAuth provider",
		zap.String("auth_url", authURL),
		zap.String("provider", provider),
		zap.String("shop_id", shopID),
		zap.String("prompt", prompt),
		zap.Bool("pkce_enabled", true),
	)

	return c.Redirect(authURL)
}

func handleCallback(c *fiber.Ctx) error {
	logger.Info("Handling OAuth callback",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("query_string", string(c.Context().URI().QueryString())),
		zap.String("user_agent", c.Get("User-Agent")),
		zap.String("remote_ip", c.IP()),
	)

	// Validate OAuth state for CSRF protection
	state := c.Cookies("oauthstate")
	queryState := c.Query("state")

	logger.Debug("Validating OAuth state",
		zap.String("cookie_state", state),
		zap.String("query_state", queryState),
	)

	// Check if state values are empty - could indicate callback reuse
	if state == "" || queryState == "" {
		logger.Error("OAuth state missing - possible callback reuse or invalid request",
			zap.String("cookie_state", state),
			zap.String("query_state", queryState),
		)
		return c.Status(http.StatusBadRequest).SendString("Missing OAuth state - callback already processed or invalid request")
	}

	if state != queryState {
		logger.Error("OAuth state mismatch - possible CSRF attack",
			zap.String("cookie_state", state),
			zap.String("query_state", queryState),
		)
		return c.Status(http.StatusUnauthorized).SendString("Invalid OAuth state")
	}

	// Retrieve code_verifier from cookie
	codeVerifier := c.Cookies("code_verifier")
	if codeVerifier == "" {
		logger.Error("Missing PKCE code verifier cookie")
		return c.Status(http.StatusBadRequest).SendString("Missing code verifier")
	}

	logger.Info("Retrieved PKCE code verifier",
		zap.String("code_verifier_length", fmt.Sprintf("%d", len(codeVerifier))),
	)

	// Retrieve app_type and shop_id from cookies
	appType := c.Cookies("app_type")
	if appType == "" {
		logger.Error("Missing app_type cookie")
		return c.Status(http.StatusBadRequest).SendString("Missing app_type cookie")
	}

	shopID := c.Cookies("shop_id")
	logger.Info("Retrieved session cookies",
		zap.String("app_type", appType),
		zap.String("shop_id", shopID),
	)

	// Only require shop_id cookie for storefront app type
	if appType == "storefront" && shopID == "" {
		logger.Error("Missing shop_id cookie for storefront app", zap.String("app_type", appType))
		return c.Status(http.StatusBadRequest).SendString("Missing shop_id cookie for storefront app")
	}

	code := c.Query("code")
	if code == "" {
		logger.Error("Missing authorization code parameter")
		return c.Status(http.StatusBadRequest).SendString("Missing authorization code")
	}

	logger.Info("Received authorization code",
		zap.String("code_length", fmt.Sprintf("%d", len(code))),
	)

	provider := c.Query("provider", "google")
	oauthProvider, exists := oauthProviders[provider]
	if !exists {
		logger.Error("Unsupported OAuth provider", zap.String("provider", provider))
		return c.Status(http.StatusBadRequest).SendString("Unsupported provider")
	}

	// Use request context with timeout for token exchange and subsequent userinfo calls
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	logger.Info("Initiating token exchange with PKCE",
		zap.String("provider", provider),
		zap.String("redirect_uri", oauthProvider.OAuth2Config.RedirectURL),
		zap.Strings("scopes", oauthProvider.OAuth2Config.Scopes),
	)

	// Exchange code for token with PKCE
	token, err := oauthProvider.OAuth2Config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		logger.Error("Token exchange failed",
			zap.Error(err),
			zap.String("provider", provider),
			zap.String("code_length", fmt.Sprintf("%d", len(code))),
			zap.String("code_verifier_length", fmt.Sprintf("%d", len(codeVerifier))),
		)
		return c.Status(http.StatusInternalServerError).SendString("Token exchange failed: " + err.Error())
	}

	logger.Info("Token exchange successful",
		zap.String("token_type", token.TokenType),
		zap.Bool("has_access_token", token.AccessToken != ""),
		zap.Bool("has_refresh_token", token.RefreshToken != ""),
		zap.Time("expires_at", token.Expiry),
	)

	// Clear OAuth-related cookies immediately after successful token exchange to prevent reuse
	logger.Debug("Clearing OAuth cookies to prevent code reuse")
	c.Cookie(&fiber.Cookie{
		Name:     "code_verifier",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "oauthstate",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})

	client := oauthProvider.OAuth2Config.Client(ctx, token)

	logger.Debug("Fetching user information from OAuth provider")
	userInfo, err := fetchUserInfo(client)
	if err != nil {
		logger.Error("Failed to fetch user info from OAuth provider",
			zap.Error(err),
			zap.String("provider", provider),
		)
		return c.Status(http.StatusInternalServerError).SendString("Failed to fetch user info: " + err.Error())
	}

	logger.Info("Successfully fetched user information",
		zap.String("user_id", userInfo.ID),
		zap.String("user_email", userInfo.Email),
		zap.String("user_name", userInfo.Name),
		zap.Bool("email_verified", userInfo.VerifiedEmail),
	)

	// Pre-flight: only register if the user does not already exist
	logger.Debug("Checking if user already exists in backend",
		zap.String("email", userInfo.Email),
		zap.String("app_type", appType),
		zap.String("shop_id", shopID),
	)

	exists, err = userExists(client, userInfo.Email, appType, shopID)
	if err != nil {
		logger.Error("Failed to check user existence in backend",
			zap.Error(err),
			zap.String("email", userInfo.Email),
			zap.String("app_type", appType),
			zap.String("shop_id", shopID),
		)
		return c.Status(http.StatusInternalServerError).SendString("Failed to check user existence: " + err.Error())
	}

	if !exists {
		logger.Info("User does not exist, registering new user",
			zap.String("email", userInfo.Email),
			zap.String("app_type", appType),
		)
		if err := registerUser(client, userInfo, appType, shopID); err != nil {
			logger.Error("Failed to register new user in backend",
				zap.Error(err),
				zap.String("email", userInfo.Email),
				zap.String("app_type", appType),
			)
			return c.Status(http.StatusInternalServerError).SendString("Failed to register user: " + err.Error())
		}
		logger.Info("Successfully registered new user", zap.String("email", userInfo.Email))
	} else {
		logger.Info("User already exists, skipping registration", zap.String("email", userInfo.Email))
	}

	loginChallenge := c.Cookies("login_challenge")
	if loginChallenge == "" {
		logger.Error("Missing login_challenge cookie")
		return c.Status(http.StatusBadRequest).SendString("Missing login_challenge")
	}

	logger.Info("Accepting Hydra login request",
		zap.String("login_challenge", loginChallenge),
		zap.String("user_email", userInfo.Email),
		zap.String("shop_id", shopID),
		zap.String("app_type", appType),
	)

	redirectTo, err := acceptHydraLogin(ctx, loginChallenge, userInfo, shopID, appType)
	if err != nil {
		logger.Error("Failed to accept login in Hydra",
			zap.Error(err),
			zap.String("login_challenge", loginChallenge),
		)
		return c.Status(http.StatusInternalServerError).SendString("Failed to accept login in Hydra: " + err.Error())
	}

	logger.Info("Successfully accepted Hydra login, redirecting",
		zap.String("redirect_to", redirectTo.RedirectTo),
		zap.String("user_email", userInfo.Email),
	)

	return c.Redirect(redirectTo.RedirectTo)
}

func fetchUserInfo(client *http.Client) (*GoogleUserInfo, error) {
	logger.Debug("Fetching user info from Google API")
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logger.Error("Failed to call Google userinfo API", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	logger.Debug("Google userinfo API response",
		zap.Int("status_code", resp.StatusCode),
		zap.String("content_type", resp.Header.Get("Content-Type")),
	)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Google userinfo API returned error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status),
		)
		return nil, fmt.Errorf("userinfo API returned status %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		logger.Error("Failed to decode Google userinfo response", zap.Error(err))
		return nil, err
	}

	logger.Debug("Successfully decoded user info from Google",
		zap.String("user_id", userInfo.ID),
		zap.String("email", userInfo.Email),
		zap.String("name", userInfo.Name),
	)

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
		logger.Error("Unsupported app_type for user existence check", zap.String("app_type", appType))
		return false, fmt.Errorf("unsupported app_type: %s", appType)
	}

	logger.Debug("Checking user existence",
		zap.String("url", urlStr),
		zap.String("email", email),
		zap.String("app_type", appType),
		zap.String("shop_id", shopID),
	)

	resp, err := client.Get(urlStr)
	if err != nil {
		logger.Error("Failed to call backend user existence API",
			zap.Error(err),
			zap.String("url", urlStr),
		)
		return false, err
	}
	defer resp.Body.Close()

	logger.Debug("Backend user existence API response",
		zap.Int("status_code", resp.StatusCode),
		zap.String("status", resp.Status),
		zap.String("url", urlStr),
	)

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info("User exists in backend", zap.String("email", email))
		return true, nil
	case http.StatusNotFound:
		logger.Info("User does not exist in backend", zap.String("email", email))
		return false, nil
	default:
		logger.Error("Unexpected status from backend user existence check",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status),
			zap.String("url", urlStr),
		)
		return false, fmt.Errorf("unexpected status %d from user existence check", resp.StatusCode)
	}
}

func registerUser(client *http.Client, userInfo *GoogleUserInfo, appType, shopID string) error {
	var urlStr string
	var payload map[string]interface{}

	logger.Info("Preparing user registration",
		zap.String("app_type", appType),
		zap.String("user_email", userInfo.Email),
		zap.String("shop_id", shopID),
	)

	if appType == "storefront" {
		// Call the register customer endpoint
		urlStr = fmt.Sprintf("%s/v1/auth/register-customer", backendURL)
		// Convert shopID to int64 for customer registration
		shopIDInt, err := strconv.ParseInt(shopID, 10, 64)
		if err != nil {
			logger.Error("Invalid shop_id for customer registration",
				zap.Error(err),
				zap.String("shop_id", shopID),
			)
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
		logger.Info("Prepared customer registration payload",
			zap.Int64("shop_id", shopIDInt),
			zap.String("email", userInfo.Email),
		)
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
		logger.Info("Prepared admin user registration payload",
			zap.String("email", userInfo.Email),
			zap.String("provider_id", userInfo.ID),
		)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal registration payload", zap.Error(err))
		return err
	}

	logger.Debug("Sending registration request to backend",
		zap.String("url", urlStr),
		zap.String("payload_size", fmt.Sprintf("%d bytes", len(body))),
	)

	resp, err := client.Post(urlStr, "application/json", bytes.NewReader(body))
	if err != nil {
		logger.Error("Failed to send registration request to backend",
			zap.Error(err),
			zap.String("url", urlStr),
		)
		return err
	}
	defer resp.Body.Close()

	logger.Debug("Backend registration API response",
		zap.Int("status_code", resp.StatusCode),
		zap.String("status", resp.Status),
		zap.String("content_type", resp.Header.Get("Content-Type")),
	)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Backend registration failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status", resp.Status),
			zap.String("email", userInfo.Email),
		)
		return fmt.Errorf("failed to register user: status %d", resp.StatusCode)
	}

	logger.Info("Successfully registered user in backend",
		zap.String("email", userInfo.Email),
		zap.String("app_type", appType),
	)
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
	logger.Debug("Preparing Hydra login acceptance",
		zap.String("login_challenge", loginChallenge),
		zap.String("subject", user.Email),
		zap.String("app_type", appType),
		zap.String("shop_id", shopID),
	)

	// Create a new map to hold the context data
	contextData := make(map[string]interface{})

	// Only include shop_id if appType is "storefront"
	if appType == "storefront" {
		contextData["shop_id"] = shopID
		logger.Debug("Added shop_id to Hydra context", zap.String("shop_id", shopID))
	}

	// Accept the login request
	acceptRequest := hydra.AcceptOAuth2LoginRequest{
		Subject:     *hydra.PtrString(user.Email),
		Remember:    hydra.PtrBool(true),
		RememberFor: hydra.PtrInt64(3600),
		Context:     contextData, // Assign the map here
	}

	logger.Debug("Sending accept login request to Hydra",
		zap.String("subject", user.Email),
		zap.Bool("remember", true),
		zap.Int64("remember_for", 3600),
	)

	// Note: ctx should be a request context provided by caller with timeout/cancellation.
	redirectTo, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2LoginRequest(ctx).
		LoginChallenge(loginChallenge).
		AcceptOAuth2LoginRequest(acceptRequest).
		Execute()
	if err != nil {
		logger.Error("Hydra accept login request failed",
			zap.Error(err),
			zap.String("login_challenge", loginChallenge),
			zap.String("subject", user.Email),
		)
		return nil, fmt.Errorf("failed to accept Hydra login request: %w", err)
	}

	logger.Info("Successfully accepted Hydra login request",
		zap.String("redirect_to", redirectTo.RedirectTo),
		zap.String("subject", user.Email),
	)

	return redirectTo, nil
}

func handleConsent(c *fiber.Ctx) error {
	logger.Info("Handling consent request",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("user_agent", c.Get("User-Agent")),
	)

	consentChallenge := c.Query("consent_challenge")
	if consentChallenge == "" {
		logger.Error("Missing consent_challenge parameter")
		return c.Status(http.StatusBadRequest).SendString("Missing consent_challenge")
	}

	logger.Info("Processing consent challenge", zap.String("consent_challenge", consentChallenge))

	// Use request context with timeout for calls to Hydra
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	logger.Debug("Calling Hydra Admin API for consent challenge")
	consentRequest, _, err := hydraAdminClient.OAuth2API.GetOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).Execute()
	if err != nil {
		logger.Error("Failed to get consent request from Hydra",
			zap.Error(err),
			zap.String("consent_challenge", consentChallenge),
		)
		return c.Status(http.StatusInternalServerError).SendString("Failed to get consent request in Hydra: " + err.Error())
	}

	logger.Info("Successfully retrieved consent request",
		zap.String("subject", consentRequest.GetSubject()),
		zap.Strings("requested_scopes", consentRequest.GetRequestedScope()),
		zap.Strings("requested_audience", consentRequest.GetRequestedAccessTokenAudience()),
	)

	logger.Debug("Accepting consent request")
	redirectTo, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2ConsentRequest(ctx).
		ConsentChallenge(consentChallenge).
		AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
			GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience,
			GrantScope:               consentRequest.RequestedScope,
			Remember:                 hydra.PtrBool(true),
			RememberFor:              hydra.PtrInt64(3600),
		}).Execute()
	if err != nil {
		logger.Error("Failed to accept consent request in Hydra",
			zap.Error(err),
			zap.String("consent_challenge", consentChallenge),
		)
		return fmt.Errorf("failed to accept Hydra consent request: %w", err)
	}

	logger.Info("Successfully accepted consent request, redirecting",
		zap.String("redirect_to", redirectTo.RedirectTo),
		zap.String("subject", consentRequest.GetSubject()),
	)

	return c.Redirect(redirectTo.RedirectTo)
}
