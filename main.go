package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"virtual-cuppa-be/config"
	"virtual-cuppa-be/handlers"
	"virtual-cuppa-be/middleware"
	"virtual-cuppa-be/repositories"
	"virtual-cuppa-be/scheduler"
	"virtual-cuppa-be/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	config.ConnectDatabase()

	userRepo := repositories.NewUserRepository(config.DB)
	orgRepo := repositories.NewOrganisationRepository(config.DB)
	tagRepo := repositories.NewTagRepository(config.DB)
	matchRepo := repositories.NewMatchRepository(config.DB)
	matchHistoryRepo := repositories.NewMatchHistoryRepository(config.DB)
	matchFeedbackRepo := repositories.NewMatchFeedbackRepository(config.DB)
	emailService := services.NewEmailService()
	authService := services.NewAuthService(userRepo, emailService)
	userService := services.NewUserService(userRepo, orgRepo, tagRepo, emailService)
	orgService := services.NewOrganisationService(orgRepo)
	matchService := services.NewMatchService(matchRepo, matchHistoryRepo, matchFeedbackRepo, userRepo, emailService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	orgHandler := handlers.NewOrganisationHandler(orgService, userService)

	// Start match scheduler
	matchScheduler := scheduler.NewMatchScheduler(matchService, orgRepo)
	matchScheduler.Start()

	matchHandler := handlers.NewMatchHandler(matchService, matchScheduler)
	feedbackHandler := handlers.NewMatchFeedbackHandler(matchService)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/request-code", authHandler.RequestCode)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/profile", authHandler.GetProfile)
		api.GET("/organisation", orgHandler.GetOrganisation)
		
		// Match endpoints for all authenticated users
		api.GET("/matches/current", matchHandler.GetCurrentMatch)
		api.GET("/matches/history", matchHandler.GetMatchHistory)
		api.PATCH("/matches/:id/accept", matchHandler.AcceptMatch)
		api.PATCH("/matches/:id/reject", matchHandler.RejectMatch)
		api.GET("/matches/:id/availabilities", matchHandler.GetMatchAvailabilities)
		
		// Feedback endpoints
		api.POST("/matches/:id/feedback", feedbackHandler.SubmitFeedback)
		api.GET("/matches/:id/feedbacks", feedbackHandler.GetMatchFeedbacks)
		api.GET("/matches/pending-feedback", feedbackHandler.GetPendingFeedback)

		admin := api.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Welcome to admin dashboard",
				})
			})
		admin.POST("/import-csv", userHandler.ImportCSV)
		admin.POST("/confirm-user", userHandler.ConfirmUser)
		admin.GET("/users", userHandler.GetOrganisationUsers)
		admin.POST("/users", userHandler.CreateUser)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
		admin.PATCH("/users/:userId/tags", userHandler.UpdateTags)
		admin.GET("/organisation", orgHandler.GetOrganisation)
		admin.PUT("/organisation", orgHandler.UpsertOrganisation)
		
		// Match endpoints for admins
		admin.POST("/matches/generate", matchHandler.GenerateMatches)
		admin.POST("/matches/trigger-scheduler", matchHandler.TriggerScheduler)
		admin.GET("/matches", matchHandler.GetOrganisationMatches)
		admin.GET("/matches/:id/feedbacks", feedbackHandler.AdminGetMatchFeedbacks)
		}
	}

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
