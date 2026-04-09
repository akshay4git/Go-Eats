package cart

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PlaceOrderRequest struct {
	Address string `json:"address" validate:"required"`
}

// @Summary     Place a new order
// @Description Places a new order from the authenticated user's cart
// @Tags        Orders
// @Accept      json
// @Produce     json
// @Param       order body PlaceOrderRequest true "Delivery address"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/order/new [post]
func (s *CartHandler) PlaceNewOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	cartInfo, err := s.service.GetCartId(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	newOrder, err := s.service.PlaceOrder(ctx, cartInfo.CartID, userID, req.Address)
	err = s.service.RemoveItemsFromCart(ctx, cartInfo.CartID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.service.NewOrderPlacedNotification(userID, newOrder.OrderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Notification failed!"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order placed!"})
}

// @Summary     List orders
// @Description Returns all orders for the authenticated user
// @Tags        Orders
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/orders [get]
func (s *CartHandler) getOrderList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	orders, err := s.service.OrderList(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// @Summary     List order items
// @Description Returns all items for a specific order
// @Tags        Orders
// @Produce     json
// @Param       id path int true "Order ID"
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/orders/{id} [get]
func (s *CartHandler) getOrderItemsList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	orders, err := s.service.OrderItemsList(ctx, userID, orderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// @Summary     Get delivery info for an order
// @Description Returns delivery information for a specific order
// @Tags        Orders
// @Produce     json
// @Param       id path int true "Order ID"
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/orders/deliveries/{id} [get]
func (s *CartHandler) getDeliveriesList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	orderID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	deliveries, err := s.service.DeliveryInformation(ctx, orderID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"delivery_info": deliveries})
}