package main

import (
    "auth-service/db"
    "auth-service/handlers"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    // Инициализация базы данных
    err := db.InitDB("user=youruser dbname=yourdb sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    router := gin.Default()

    // Роуты
    router.POST("/auth/register", handlers.Register)
    router.POST("/auth/login", handlers.Login)
    router.POST("/auth/refresh", handlers.Refresh)

    log.Println("Server is running on http://localhost:8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}