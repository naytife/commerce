package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/gql/public/generated"
	"github.com/petrejonn/naytife/internal/gql/public/resolver"
)

func NewHandler(repo db.Repository) fiber.Handler {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			Repository: repo,
		},
	}))

	return adaptor.HTTPHandler(h)
}
func NewPlaygroundHandler(endpoint string) fiber.Handler {
	h := playground.AltairHandler("Naytife Storefront Playground", endpoint)

	return adaptor.HTTPHandler(h)
}
