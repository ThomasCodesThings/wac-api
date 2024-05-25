package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/ThomasCodesThings/wac-api/internal/department"
    //"github.com/ThomasCodesThings/wac-api/internal/db_service"
    //"context"
    "time"
    "github.com/gin-contrib/cors"
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
    corsMiddleware := cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders:    []string{""},
        AllowCredentials: false,
        MaxAge: 12 * time.Hour,
    })
    server.Use(corsMiddleware)

    //dbService := db_service.NewMongoService[department.Operation](db_service.MongoServiceConfig{})
    //defer dbService.Disconnect(context.Background())
    //server.Use(func(ctx *gin.Context) {
        //ctx.Set("db_service", dbService)
        //ctx.Next()
    //})
    // Create a new instance of DepartmentAPI
    api := department.NewDepartmentAPI()

    // Register routes
    apiGroup := server.Group("/api")
    api.AddRoutes(apiGroup)

    // Request routings
    server.Run(":" + port)
}
