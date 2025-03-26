package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Permission struct {
	Name   string `json:"name"`
	Module string `json:"module"`
}

type JWTClaims struct {
	ID          uint         `json:"id"`
	Username    string       `json:"username"`
	Roles       []string     `json:"roles"`
	Permissions []Permission `json:"permissions"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userID uint, username string, roles []string, permissions []Permission) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey     string
	expirationStr string
}

func NewJWTService(secretKey string, expirationStr string) JWTService {
	return &jwtService{
		secretKey:     secretKey,
		expirationStr: expirationStr,
	}
}

func (s *jwtService) GenerateToken(userID uint, username string, roles []string, permissions []Permission) (string, error) {
	// Parse expiration duration
	expiration, err := time.ParseDuration(s.expirationStr)
	if err != nil {
		return "", err
	}

	claims := &JWTClaims{
		ID:          userID,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
