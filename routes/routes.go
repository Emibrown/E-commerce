package routes

import (
	"github.com/Emibrown/E-commerce-API/controllers"
	"github.com/Emibrown/E-commerce-API/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/register-admin", controllers.RegisterAdmin)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		// Orders (User only)
		api.POST("/orders", controllers.CreateOrder)
		api.GET("/orders", controllers.GetOrders)
		api.PUT("/orders/:id/cancel", controllers.CancelOrder)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			// Product management
			admin.POST("/products", controllers.CreateProduct)
			admin.GET("/products", controllers.GetProducts)
			admin.GET("/products/:id", controllers.GetProductByID)
			admin.PUT("/products/:id", controllers.UpdateProduct)
			admin.DELETE("/products/:id", controllers.DeleteProduct)

			// Order status update
			admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)
		}
	}
}
