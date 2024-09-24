package middleware

import (
	"context"
	"net/http"

	"github.com/petrejonn/naytife/internal/db"
)

// Add the key for the shop_id in the context
type contextKey string

const shopIDKey = contextKey("shop_id")

func ShopIDMiddleware(repo db.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the host from the HTTP request
			host := r.Host

			// Query the database to get the shop_id using the domain (host)
			shopID, err := repo.GetShopIDByDomain(r.Context(), host)
			if err != nil {
				// Handle error (e.g., log or return unauthorized response)
				http.Error(w, "Invalid shop", http.StatusUnauthorized)
				return
			}

			// Store the shopID in the context
			ctx := context.WithValue(r.Context(), shopIDKey, shopID)

			// Pass the request along with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
