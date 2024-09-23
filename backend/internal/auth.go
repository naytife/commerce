// auth.go
package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/petrejonn/naytife/config"
)

var (
	enforcer *casbin.Enforcer
)

// InitializeCasbin initializes Casbin with a model and policy.
func InitializeCasbin() error {
	m, err := model.NewModelFromString(`
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [role_definition]
    g = _, _

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
    `)
	if err != nil {
		return err
	}

	a := fileadapter.NewAdapter("path/to/policy.csv") // Define your policy file
	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		return err
	}

	return enforcer.LoadPolicy()
}

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope             string `json:"scope"`
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	ProfilePictureURL string `json:"picture"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// ParseAuth0Token parses and validates the Auth0 token and returns the custom claims.
// ParseAuth0Token parses and validates the Auth0 token and returns the custom claims.
func ParseAuth0Token(ctx context.Context, tokenString string) (*CustomClaims, error) {
	// Load environment variables
	env, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Prepare the issuer URL
	issuerURL, err := url.Parse("https://" + env.AUTH0_DOMAIN + "/")
	if err != nil {
		return nil, err
	}

	// Create a caching provider for the JWKS (JSON Web Key Set)
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	// Setup JWT validator
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{env.AUTH0_AUDIENCE}, // This will verify the 'aud' claim
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	// Validate the JWT token
	token, err := jwtValidator.ValidateToken(ctx, tokenString)
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		return nil, err
	}

	// Extract custom claims from the token
	claims, ok := token.(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("invalid token claims structure")
	}

	// Retrieve custom claims (like scope, sub, etc.)
	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return nil, errors.New("unable to parse custom claims")
	}

	// Return the custom claims
	return customClaims, nil
}

// JWTMiddleware validates JWT tokens from Auth0
func JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// host, _, err := net.SplitHostPort(r.Host)
			// if err != nil {
			// 	host = r.Host
			// 	// http.Error(w, "Invalid host", http.StatusInternalServerError)
			// }
			// Extract the token from the Authorization header
			// authHeader := r.Header.Get("Authorization")
			// if authHeader == "" {
			// 	http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			// 	return
			// }

			// tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate the token using ParseAuth0Token
			// claims, err := ParseAuth0Token(r.Context(), tokenString)
			// if err != nil {
			// 	http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			// 	return
			// }

			// Store the claims in the request context
			// ctx := context.WithValue(r.Context(), "userClaims", claims)

			// Set the host in the request context
			ctx := context.WithValue(r.Context(), "shopHost", r.Host)

			// Pass the request along with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CasbinMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the user claims from the context
			claims, ok := r.Context().Value("userClaims").(*CustomClaims)
			if !ok {
				http.Error(w, "Unable to retrieve user claims", http.StatusForbidden)
				return
			}

			// Extract the necessary information (e.g., the user's role or scope)
			userRole := claims.Scope // Or use another field, depending on how roles are modeled

			// Define the object (the resource) and action (the method)
			// For example, obj can be the URL path or endpoint being accessed
			obj := r.URL.Path
			act := r.Method

			// Enforce the Casbin policy
			ok, err := enforcer.Enforce(userRole, obj, act)
			if err != nil {
				http.Error(w, "Error occurred while enforcing policy", http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			// Pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
