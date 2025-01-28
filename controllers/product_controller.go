package controllers

import (
	"net/http"
	"strconv"

	"github.com/Emibrown/E-commerce-API/config"
	"github.com/Emibrown/E-commerce-API/models"
	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Adds a new product to the store (admin only)
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body   CreateProductInput         true  "Create Product Input"
// @Success      201   {object} CreateProductResponse
// @Failure      400,401,403,500 {object} ErrorResponse
// @Router       /api/admin/products [post]
func CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Convert CreateProductInput into a models.Product
	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create product"})
		return
	}

	c.JSON(http.StatusCreated, CreateProductResponse{
		Data: ProductPayload{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	})
}

// GetProducts godoc
// @Summary      Get all products
// @Description  Returns a list of all products (admin only endpoint in this example)
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Success      200   {object} GetProductsResponse
// @Failure      500   {object} ErrorResponse
// @Router       /api/admin/products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		// If there's a real DB error (e.g., connection issue) then respond 500
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	var payloads []ProductPayload
	for _, p := range products {
		payloads = append(payloads, ProductPayload{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, GetProductsResponse{Data: payloads})
}

// GetProductByID godoc
// @Summary      Get a product by its ID
// @Description  Returns a single product (admin only endpoint in this example)
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  SingleProductResponse
// @Failure      400,404 {object} ErrorResponse
// @Router       /api/admin/products/{id} [get]
func GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
		return
	}

	payload := ProductPayload{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	c.JSON(http.StatusOK, SingleProductResponse{Data: payload})
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product (admin only)
// @Tags         products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int              true  "Product ID"
// @Param        body body      models.Product   true  "Product Data"
// @Success      200  {object}  UpdateProductResponse
// @Failure      400,404,500 {object} ErrorResponse
// @Router       /api/admin/products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, UpdateProductResponse{
		Data: ProductPayload{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		},
	})
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Deletes an existing product (admin only)
// @Tags         products
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  DeleteProductResponse
// @Failure      400,404,500 {object} ErrorResponse
// @Router       /api/admin/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
		return
	}

	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, DeleteProductResponse{
		Message: "Product deleted",
	})
}
