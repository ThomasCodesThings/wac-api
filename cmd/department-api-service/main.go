package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/ThomasCodesThings/wac-api/api"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("API_ENVIRONMENT")
    if !strings.EqualFold(environment, "production") { // case insensitive comparison
        gin.SetMode(gin.DebugMode)
    }
    server := gin.New()
    server.Use(gin.Recovery())
    // request routings
    server.GET("/openapi", api.HandleOpenApi)
    server.Run(":" + port)
}