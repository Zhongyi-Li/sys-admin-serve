package health

import (
	"time"

	"sys-admin-serve/internal/pkg/config"
	"sys-admin-serve/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{config: cfg}
}

func (h *Handler) GetHealth(c *gin.Context) {
	response.Success(c, gin.H{
		"app":       h.config.App.Name,
		"env":       h.config.App.Env,
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
