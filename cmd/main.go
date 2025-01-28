package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Emibrown/E-commerce-API/config"
	"github.com/Emibrown/E-commerce-API/docs"
	"github.com/Emibrown/E-commerce-API/models"
	"github.com/Emibrown/E-commerce-API/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           E-commerce API
// @version         1.0
// @description     This is a sample e-commerce server with Gin and GORM.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost
// @BasePath  /api

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and JWT token.

// @accept json
// @produce json

func main() {
	// Auto-migrate models
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	docs.SwaggerInfo.BasePath = "/"

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server on the specified port
	r.Run(":" + port)
}
