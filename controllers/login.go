package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

// Define a secret key for signing JWT tokens
var jwtSecret = []byte("secret")

type User struct {
    Username string
    Password string
    UserType string
}

var users = map[string]User{
    "user1": {"user1", "password1", "regular"},
    "user2": {"user2", "password2", "admin"},
}

func LoginHandler(c *gin.Context) {
    var requestBody struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    user, ok := users[requestBody.Username]
    if !ok || user.Password != requestBody.Password {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = user.Username
    claims["usertype"] = user.UserType
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() 

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
