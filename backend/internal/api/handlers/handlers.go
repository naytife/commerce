package handlers

import (
	"github.com/petrejonn/naytife/internal/db"
)

type Handler struct {
	Repository db.Repository
}

// NewHandler returns a handler with a repository
func NewHandler(repo db.Repository) *Handler {
	return &Handler{Repository: repo}
}
