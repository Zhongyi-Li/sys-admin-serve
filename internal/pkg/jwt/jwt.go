package jwt

import (
	"fmt"
	"time"

	"sys-admin-serve/internal/pkg/config"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint64   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwtv5.RegisteredClaims
}

type Manager struct {
	secret      []byte
	expireHours int
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{
		secret:      []byte(cfg.Secret),
		expireHours: cfg.ExpireHours,
	}
}

func (m *Manager) GenerateToken(userID uint64, username string, roles []string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(m.expireHours) * time.Hour)
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwtv5.RegisteredClaims{
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(expiresAt),
			NotBefore: jwtv5.NewNumericDate(now),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return signed, expiresAt, nil
}

func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (any, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
