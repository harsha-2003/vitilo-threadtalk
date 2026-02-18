package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/api/handlers"
	"github.com/harsha-2003/vitilo-threadtalk/backend/internal/api/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	postHandler := handlers.NewPostHandler(db)
	communityHandler := handlers.NewCommunityHandler(db)
	voteHandler := handlers.NewVoteHandler(db)
	commentHandler := handlers.NewCommentHandler(db)

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Vitilo ThreadTalk API is running",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("/auth/me", authHandler.GetCurrentUser)

			// Community routes
			protected.GET("/communities", communityHandler.GetCommunities)
			protected.POST("/communities", communityHandler.CreateCommunity)
			protected.GET("/communities/user/joined", communityHandler.GetUserCommunities) // BEFORE :id route
			protected.GET("/communities/:id", communityHandler.GetCommunity)
			protected.POST("/communities/:id/join", communityHandler.JoinCommunity)
			protected.POST("/communities/:id/leave", communityHandler.LeaveCommunity)

			// Post routes
			protected.GET("/posts", postHandler.GetPosts)
			protected.POST("/posts", postHandler.CreatePost)
			protected.POST("/posts/upload", postHandler.UploadImage) // BEFORE :id route
			protected.GET("/posts/:id", postHandler.GetPost)
			protected.DELETE("/posts/:id", postHandler.DeletePost)
			protected.GET("/posts/:id/comments", commentHandler.GetComments) // Changed from :post_id to :id

			// Vote routes
			protected.POST("/posts/:id/vote", voteHandler.VotePost)
			protected.POST("/comments/:id/vote", voteHandler.VoteComment)

			// Comment routes
			protected.POST("/comments", commentHandler.CreateComment)
			protected.DELETE("/comments/:id", commentHandler.DeleteComment)
		}
	}
}
