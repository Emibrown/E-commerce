package controllers

import (
	"net/http"

	"github.com/Emibrown/E-commerce-API/config"
	"github.com/Emibrown/E-commerce-API/models"
	"github.com/gin-gonic/gin"
)

type OrderRequest struct {
	Items []struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	} `json:"items"`
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Places a new order for the authenticated user
// @Tags         orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body   OrderRequest  true  "Order Data"
// @Success      201  {object} CreateOrderResponse
// @Failure      400,401,500 {object} ErrorResponse
// @Router       /api/orders [post]
func CreateOrder(c *gin.Context) {
	userId := c.GetUint("user_id")

	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var orderItems []models.OrderItem
	for _, item := range req.Items {
		var product models.Product
		if err := config.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Product not found"})
			return
		}

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order := models.Order{
		UserID:   userId,
		Products: orderItems,
		Status:   models.Pending,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create order"})
		return
	}

	config.DB.Preload("Products.Product").First(&order, order.ID)

	// Prepare response payload
	itemPayloads := make([]OrderItemPayload, 0, len(orderItems))
	for _, oi := range orderItems {
		itemPayloads = append(itemPayloads, OrderItemPayload{
			ProductID:   oi.ProductID,
			Name:        oi.Product.Name,
			Description: oi.Product.Description,
			Quantity:    oi.Quantity,
			Price:       oi.Price,
		})
	}
	c.JSON(http.StatusCreated, CreateOrderResponse{
		Data: OrderPayload{
			ID:        order.ID,
			UserID:    order.UserID,
			Status:    string(order.Status),
			Products:  itemPayloads,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		},
	})
}

// GetOrders godoc
// @Summary      Get all orders for the authenticated user
// @Description  Returns a list of orders belonging to the logged-in user
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} GetOrdersResponse
// @Failure      401,500 {object} ErrorResponse
// @Router       /api/orders [get]
func GetOrders(c *gin.Context) {
	userId := c.GetUint("user_id")

	var orders []models.Order
	// Preload both the OrderItems ("Products") and each OrderItem's "Product"
	if err := config.DB.Preload("Products.Product").
		Where("user_id = ?", userId).
		Find(&orders).Error; err != nil {

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch orders"})
		return
	}

	// Convert to payload
	var orderPayloads []OrderPayload
	for _, o := range orders {
		var itemPayloads []OrderItemPayload
		for _, item := range o.Products {
			itemPayloads = append(itemPayloads, OrderItemPayload{
				ProductID:   item.ProductID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Quantity:    item.Quantity,
				Price:       item.Price,
			})
		}

		orderPayloads = append(orderPayloads, OrderPayload{
			ID:        o.ID,
			UserID:    o.UserID,
			Status:    string(o.Status),
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			Products:  itemPayloads,
		})
	}

	c.JSON(http.StatusOK, GetOrdersResponse{Data: orderPayloads})
}

// CancelOrder godoc
// @Summary      Cancel an order
// @Description  Cancels an order if it's still in 'Pending' status
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object} CancelOrderResponse
// @Failure      400,401,404,500 {object} ErrorResponse
// @Router       /api/orders/{id}/cancel [put]
func CancelOrder(c *gin.Context) {
	userId := c.GetUint("user_id")
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.Preload("Products.Product").
		Where("id = ? AND user_id = ?", orderID, userId).
		First(&order).Error; err != nil {

		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		return
	}

	if order.Status != models.Pending {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Cannot cancel an order that is not Pending"})
		return
	}

	order.Status = models.Cancelled
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, CancelOrderResponse{
		Message: "Order cancelled",
	})
}

// UpdateOrderStatus godoc
// @Summary      Update an order status
// @Description  Allows an admin to update order status (e.g., Shipped, Completed, Cancelled)
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        id     path   int      true  "Order ID"
// @Param        status query  string   true  "New Status (Shipped|Completed|Cancelled)"
// @Success      200    {object} UpdateOrderStatusResponse
// @Failure      400,401,403,404,500 {object} ErrorResponse
// @Router       /api/admin/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	newStatus := c.Query("status")

	var order models.Order
	if err := config.DB.Preload("Products").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		return
	}

	// Validate newStatus
	if newStatus != string(models.Shipped) &&
		newStatus != string(models.Completed) &&
		newStatus != string(models.Cancelled) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid order status"})
		return
	}

	order.Status = models.OrderStatus(newStatus)
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update order status"})
		return
	}

	config.DB.Preload("Products.Product").First(&order, order.ID)

	// Build response payload
	var itemPayloads []OrderItemPayload
	for _, item := range order.Products {
		itemPayloads = append(itemPayloads, OrderItemPayload{
			ProductID:   item.ProductID,
			Name:        item.Product.Name,
			Description: item.Product.Description,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	c.JSON(http.StatusOK, UpdateOrderStatusResponse{
		Data: OrderPayload{
			ID:        order.ID,
			UserID:    order.UserID,
			Status:    string(order.Status),
			Products:  itemPayloads,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		},
	})
}
