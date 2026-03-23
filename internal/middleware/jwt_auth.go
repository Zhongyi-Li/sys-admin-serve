package middleware

import (
	"strings"

	"sys-admin-serve/internal/pkg/jwt"
	"sys-admin-serve/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const authClaimsContextKey = "auth_claims"

func JWTAuth(jwtManager *jwt.Manager, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		tokenString, ok := extractBearerToken(authorization)
		if !ok {
			response.Unauthorized(c, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			log.Warn("jwt validation failed", zap.Error(err))
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(authClaimsContextKey, claims)
		c.Next()
	}
}

func CurrentClaims(c *gin.Context) (*jwt.Claims, bool) {
	value, ok := c.Get(authClaimsContextKey)
	if !ok {
		return nil, false
	}

	claims, ok := value.(*jwt.Claims)
	return claims, ok
}

func extractBearerToken(authorization string) (string, bool) {
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return strings.TrimSpace(parts[1]), true
}
