package service

import (
	"fmt"
	"net/http"
	"social-network-service/internal/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	secret          = "secret-key"
	durationInHours = 1
)

type JwtService struct {
}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (s JwtService) GenerateToken(userId model.UserId) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * durationInHours).Unix(),
		"sub": userId,
	})

	tokenStr, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("failed to generate token for user %v: %w", userId, err)
	}

	return tokenStr, nil
}

func (s JwtService) GetUserId(c *gin.Context) (model.UserId, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("header Authorization not found")
	}

	token, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found {
		return "", fmt.Errorf("prefix Bearer is not present")
	}

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("failed to get claims: %w", err)
	}

	userIdStr := claims["sub"].(string)
	userId := model.UserId(userIdStr)

	return userId, nil
}

func (s JwtService) GetUserIdFromRequest(r *http.Request) (model.UserId, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("header Authorization not found")
	}

	token, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found {
		return "", fmt.Errorf("prefix Bearer is not present")
	}

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("failed to get claims: %w", err)
	}

	userIdStr := claims["sub"].(string)
	userId := model.UserId(userIdStr)

	return userId, nil
}

func (s JwtService) CheckAccess(c *gin.Context) (model.UserId, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", model.NewUnauthenticatedError("no authorization header provided", nil)
	}

	token, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found {
		return "", model.NewUnauthenticatedError("no Bearer prefix", nil)
	}

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", model.NewUnauthenticatedError("failed to check JWT token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("failed to get claims: %w", err)
	}

	// For some reason 'exp' is parsed as float64.
	exp := int64(claims["exp"].(float64))
	now := time.Now().Unix()

	if now > exp {
		return "", model.NewUnauthenticatedError("token expired", nil)
	}

	userIdStr := claims["sub"].(string)
	userId := model.UserId(userIdStr)

	return userId, nil
}

func (s JwtService) CheckAccessFromRequest(r *http.Request) (model.UserId, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", model.NewUnauthenticatedError("no authorization header provided", nil)
	}

	token, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found {
		return "", model.NewUnauthenticatedError("no Bearer prefix", nil)
	}

	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", model.NewUnauthenticatedError("failed to check JWT token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("failed to get claims: %w", err)
	}

	// For some reason 'exp' is parsed as float64.
	exp := int64(claims["exp"].(float64))
	now := time.Now().Unix()

	if now > exp {
		return "", model.NewUnauthenticatedError("token expired", nil)
	}

	userIdStr := claims["sub"].(string)
	userId := model.UserId(userIdStr)

	return userId, nil
}
