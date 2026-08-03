package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/configs"
	"github.com/mjhddev/go-ecommerce-api/internal/handlers"
	"github.com/mjhddev/go-ecommerce-api/internal/middleware"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
	"github.com/mjhddev/go-ecommerce-api/internal/services"
)

func SetupRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Dependency Injection
	userRepository := repositories.NewUserRepository(configs.DB)
	userService := services.NewUserService(userRepository)
	authHandler := handlers.NewAuthHandler(userService)

	categoryRepository := repositories.NewCategoryRepository(configs.DB)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Routes

	api := router.Group("/api/v1")

	profile := api.Group("/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.GET("", authHandler.Profile)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	category := api.Group("/categories")
	category.Use(middleware.AuthMiddleware())
	{
		category.POST("", categoryHandler.Create)
		category.GET("", categoryHandler.GetAll)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "API is running",
		})
	})

	return router
}
