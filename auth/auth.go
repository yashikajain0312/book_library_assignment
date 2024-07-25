package auth

import (
    "net/http"
    "strings"
    "errors"
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go" 
)

func GetUserTypeFromToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil 
	})
	if err != nil {
		return "", err
	}

	// Extract user type from JWT claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["usertype"].(string), nil
	}

	var ErrUserTypeNotFound = errors.New("user type not found in JWT claims")
	return "", ErrUserTypeNotFound
}