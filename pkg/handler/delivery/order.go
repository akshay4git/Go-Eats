package delivery

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Ayocodes24/GO-Eats/pkg/database/models/delivery"
	"github.com/gin-gonic/gin"
)

// @Summary     Update order status
// @Description Delivery person updates an order's delivery status
// @Tags        Delivery
// @Accept      json
// @Produce     json
// @Param       order body delivery.DeliveryOrderPlacementParams true "Order ID and new status"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Security    BearerAuth
// @Router      /delivery/update-order [post]
func (s *DeliveryHandler) updateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var deliveryOrder delivery.DeliveryOrderPlacementParams
	if err := c.BindJSON(&deliveryOrder); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetInt64("userID")
	_, err := s.service.OrderPlacement(ctx, userID, deliveryOrder.OrderID, deliveryOrder.Status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order Updated!"})
}

// @Summary     List deliveries for an order
// @Description Returns all delivery records for a specific order
// @Tags        Delivery
// @Produce     json
// @Param       order_id path int true "Order ID"
// @Success     200 {object} map[string]interface{}
// @Failure     404 {object} map[string]string
// @Security    BearerAuth
// @Router      /delivery/deliveries/{order_id} [get]
func (s *DeliveryHandler) deliveryListing(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderId := c.Param("order_id")
	orderID, _ := strconv.ParseInt(orderId, 10, 64)
	userID := c.GetInt64("userID")

	deliveries, err := s.service.DeliveryListing(ctx, orderID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}