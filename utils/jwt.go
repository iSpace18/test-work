package utils

import (
    "github.com/dgrijalva/jwt-go"
    "time"
)

var secretKey = []byte("your_secret_key_here")

type Claims struct {
    IP string `json:"ip"`
    jwt.StandardClaims
}

func GenerateTokens(ip string) (string, string, error) {
    accessTokenExpiry := time.Now().Add(time.Minute * 15)
    refreshTokenExpiry := time.Now().Add(time.Hour * 24 * 7)

    accessClaims := Claims{
        IP: ip,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: accessTokenExpiry.Unix(),
        },
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

    refreshClaims := Claims{
        IP: ip,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: refreshTokenExpiry.Unix(),
        },
    }
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

    signedAccessToken, err := accessToken.SignedString(secretKey)
    if err != nil {
        return "", "", err
    }

    signedRefreshToken, err := refreshToken.SignedString(secretKey)
    if err != nil {
        return "", "", err
    }

    return signedAccessToken, signedRefreshToken, nil
}