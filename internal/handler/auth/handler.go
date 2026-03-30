package auth

import (
	"errors"

	authdto "sys-admin-serve/internal/dto/auth"
	"sys-admin-serve/internal/middleware"
	"sys-admin-serve/internal/response"
	serviceauth "sys-admin-serve/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *serviceauth.Service
}

func NewHandler(authService *serviceauth.Service) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid register request")
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, serviceauth.ErrInvalidRegister):
			response.BadRequest(c, "invalid register request")
		case errors.Is(err, serviceauth.ErrUsernameExists):
			response.BadRequest(c, "username already exists")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid login request")
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, serviceauth.ErrInvalidCredentials):
			response.Unauthorized(c, "invalid username or password")
		case errors.Is(err, serviceauth.ErrUserDisabled):
			response.Forbidden(c, "user is disabled")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) Me(c *gin.Context) {
	claims, ok := middleware.CurrentClaims(c)
	if !ok {
		response.Unauthorized(c, "missing auth claims")
		return
	}

	user, err := h.authService.GetCurrentUser(c.Request.Context(), claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, serviceauth.ErrUserNotFound):
			response.NotFound(c, "user not found")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, user)
}

func (h *Handler) Menus(c *gin.Context) {
	claims, ok := middleware.CurrentClaims(c)
	if !ok {
		response.Unauthorized(c, "missing auth claims")
		return
	}

	menus, err := h.authService.GetCurrentUserMenus(c.Request.Context(), claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, serviceauth.ErrUserNotFound):
			response.NotFound(c, "user not found")
		default:
			response.InternalError(c)
		}
		return
	}

	response.Success(c, menus)
}
