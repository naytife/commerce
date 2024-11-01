package handlers

import (
	"github.com/petrejonn/naytife/internal/db"
)

type Handler struct {
	Repository db.Repository
}

func NewHandler(repo db.Repository) *Handler {
	return &Handler{Repository: repo}
}
