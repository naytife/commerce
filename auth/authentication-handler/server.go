package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	hydra "github.com/ory/hydra-client-go/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	hydraAdminURL    string
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

func main() {
	app := fiber.New(fiber.Config{
		StrictRouting:  true,
		ServerHeader:   "Fiber",
		ReadBufferSize: 8192,
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
	log.Fatal(app.Listen(":3000"))
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
	c.Cookie(&fiber.Cookie{
		Name:     "login_challenge",
		Value:    loginChallenge,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "Lax",
		HTTPOnly: true,
	})

	state, err := generateStateOauthCookie()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to generate state")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		SameSite: "None",
		HTTPOnly: true,
	})

	authURL := oauthProvider.OAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(authURL)
}

func handleCallback(c *fiber.Ctx) error {
	state := c.Cookies("oauthstate")
	if state != c.Query("state") {
		return c.Status(http.StatusUnauthorized).SendString("Invalid OAuth state")
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

	token, err := oauthProvider.OAuth2Config.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Token exchange failed: " + err.Error())
	}

	client := oauthProvider.OAuth2Config.Client(context.Background(), token)
	userInfo, err := fetchUserInfo(client)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to fetch user info: " + err.Error())
	}

	if err := registerUser(client, userInfo); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to register user: " + err.Error())
	}

	loginChallenge := c.Cookies("login_challenge")
	if loginChallenge == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing login_challenge")
	}
	redirectTo, err := acceptHydraLogin(loginChallenge, userInfo)
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

func registerUser(client *http.Client, userInfo *GoogleUserInfo) error {
	// Marshal userInfo to JSON
	// TODO: send provider
	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		return err // Handle JSON marshalling error
	}

	// Use bytes.NewReader to create an io.Reader
	resp, err := client.Post("http://host.minikube.internal:8000/api/v1/auth/register", "application/json", bytes.NewReader(userInfoJSON))
	if err != nil {
		return err // Handle HTTP POST error
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send user info: status code %d", resp.StatusCode)
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

func acceptHydraLogin(loginChallenge string, user *GoogleUserInfo) (*hydra.OAuth2RedirectTo, error) {
	// Accept the login request
	acceptRequest := hydra.AcceptOAuth2LoginRequest{
		Subject:     *hydra.PtrString(user.Email),
		Remember:    hydra.PtrBool(true),
		RememberFor: hydra.PtrInt64(3600),
	}

	redirectTo, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).
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
	consentRequest, _, err := hydraAdminClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).Execute()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get consent request in Hydra: " + err.Error())
	}
	redirectTo, _, err := hydraAdminClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).
		ConsentChallenge(consentChallenge).
		AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
			GrantAccessTokenAudience: consentRequest.RequestedAccessTokenAudience, // List any audiences for the token
			GrantScope:               consentRequest.RequestedScope,
			Remember:                 hydra.PtrBool(true),
			RememberFor:              hydra.PtrInt64(3600),
		}).Execute()
	if err != nil {
		return fmt.Errorf("failed to accept Hydra consent request: %w", err)
	}
	return c.Redirect(redirectTo.RedirectTo)

}
