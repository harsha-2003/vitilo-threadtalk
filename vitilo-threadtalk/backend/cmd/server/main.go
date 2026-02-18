package main

import (
    "log"
    "os"
    
    "github.com/gin-gonic/gin"
    "github.com/harsha-2003/vitilo-threadtalk/backend/internal/api/routes"
    "github.com/harsha-2003/vitilo-threadtalk/backend/internal/config"
    "github.com/harsha-2003/vitilo-threadtalk/backend/internal/models"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize database
    db := config.InitDB()
    
    // Auto migrate models
    if err := db.AutoMigrate(
        &models.User{},
        &models.Community{},
        &models.Post{},
        &models.Comment{},
        &models.Vote{},
        &models.CommunityMember{},
    ); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Create uploads directory
    if err := os.MkdirAll("uploads", 0755); err != nil {
        log.Fatal("Failed to create uploads directory:", err)
    }

    // Initialize Gin router
    r := gin.Default()
    
    // Setup CORS
    r.Use(config.CORSMiddleware())
    
    // Setup routes
    routes.SetupRoutes(r, db)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s...", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}