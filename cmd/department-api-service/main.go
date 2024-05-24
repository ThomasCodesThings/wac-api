package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/ThomasCodesThings/wac-api/internal/department"
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
    } else {
        gin.SetMode(gin.ReleaseMode)
    }

    server := gin.New()
    server.Use(gin.Recovery())

    // Create a new instance of DepartmentAPI
    api := department.NewDepartmentAPI()

    // Register routes
    apiGroup := server.Group("/api")
    api.AddRoutes(apiGroup)

    // Request routings
    server.Run(":" + port)
}
