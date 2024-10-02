package handlers

import (
    "auth-service/db"
    "auth-service/models"
    "auth-service/utils"
    "github.com/gin-gonic/gin"
    "net/http"
)

func Register(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    user.ID = uuid.New()
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
        return
    }
    user.RefreshTokenHash = hashedPassword

    _, err = db.DB.NamedExec(`INSERT INTO users (id, email, refresh_token_hash, last_ip) VALUES (:id, :email, :refresh_token_hash, :last_ip)`, &user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
    var loginInfo models.User
    if err := c.ShouldBindJSON(&loginInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var user models.User
    err := db.DB.Get(&user, "SELECT * FROM users WHERE email=$1", loginInfo.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if !utils.CheckPasswordHash(loginInfo.Password, user.RefreshTokenHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    accessToken, refreshToken, err := utils.GenerateTokens(c.ClientIP())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func Refresh(c *gin.Context) {
    refreshToken := c.PostForm("refresh_token")
    
    claims := &utils.Claims{}
    token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
        return utils.SecretKey, nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        return
    }

    accessToken, newRefreshToken, err := utils.GenerateTokens(claims.IP)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new tokens"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": newRefreshToken})
}