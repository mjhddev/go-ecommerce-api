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

	productRepository := repositories.NewProductRepository(configs.DB)
	productService := services.NewProductService(productRepository, categoryRepository)
	productHandler := handlers.NewProductHandler(productService)

	cartRepository := repositories.NewCartRepository(configs.DB)
	cartService := services.NewCartService(cartRepository, productRepository)
	cartHandler := handlers.NewCartHandler(cartService)

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
		category.GET("", categoryHandler.GetAll)
		category.GET("/:id", categoryHandler.GetByID)

		category.POST("", middleware.RoleMiddleware("admin"), categoryHandler.Create)
		category.PUT("/:id", middleware.RoleMiddleware("admin"), categoryHandler.Update)
		category.DELETE("/:id", middleware.RoleMiddleware("admin"), categoryHandler.Delete)
	}

	product := api.Group("/products")
	product.Use(middleware.AuthMiddleware())
	{
		product.GET("", productHandler.GetAll)
		product.GET("/:id", productHandler.GetByID)

		product.POST("", middleware.RoleMiddleware("admin"), productHandler.Create)
		product.PUT("/:id", middleware.RoleMiddleware("admin"), productHandler.Update)
		product.DELETE("/:id", middleware.RoleMiddleware("admin"), productHandler.Delete)
	}

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("", cartHandler.AddToCart)
		cart.GET("", cartHandler.GetCart)
		cart.PUT("/:id", cartHandler.Update)
		cart.DELETE("/:id", cartHandler.Delete)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "API is running",
		})
	})

	return router
}
