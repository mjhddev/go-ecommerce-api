package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjhddev/go-ecommerce-api/configs"
	"github.com/mjhddev/go-ecommerce-api/internal/handlers"
	"github.com/mjhddev/go-ecommerce-api/internal/middleware"
	"github.com/mjhddev/go-ecommerce-api/internal/repositories"
	"github.com/mjhddev/go-ecommerce-api/internal/response"
	"github.com/mjhddev/go-ecommerce-api/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/uploads", "./uploads")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	orderRepository := repositories.NewOrderRepository(configs.DB)
	orderService := services.NewOrderService(
		configs.DB,
		orderRepository,
		cartRepository,
		productRepository,
	)

	orderHandler := handlers.NewOrderHandler(orderService)

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
		product.PUT("/:id/image", middleware.RoleMiddleware("admin"), productHandler.UploadImage)

	}

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("", cartHandler.AddToCart)
		cart.GET("", cartHandler.GetCart)
		cart.PUT("/:id", cartHandler.Update)
		cart.DELETE("/:id", cartHandler.Delete)
	}

	order := api.Group("/orders")
	order.Use(middleware.AuthMiddleware())
	{
		order.GET("", orderHandler.GetOrders)
		order.GET("/:id", orderHandler.GetOrderByID)
		order.POST("/checkout", orderHandler.Checkout)
	}

	router.GET("/health", func(c *gin.Context) {
		response.Success(
			c,
			http.StatusOK,
			"API is running",
			nil,
		)
	})

	return router
}
