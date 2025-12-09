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
	"virtual-cuppa-be/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	config.ConnectDatabase()

	userRepo := repositories.NewUserRepository(config.DB)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

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
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/profile", authHandler.GetProfile)

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
		}
	}

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
