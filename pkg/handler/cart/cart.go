package cart

import (
	"context"
	"net/http"
	"strconv"
	"time"

	cartModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/cart"
	"github.com/gin-gonic/gin"
)

// @Summary     Add item to cart
// @Description Adds a menu item to the authenticated user's cart
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Param       item body cart.CartItemParams true "Cart item details"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/add [post]
func (s *CartHandler) addToCart(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var cartID int64
	userID := c.GetInt64("userID")

	cartInfo, err := s.service.GetCartId(ctx, userID)
	if err != nil {
		var cartData cartModel.Cart
		cartData.UserID = userID
		newCart, err := s.service.Create(ctx, &cartData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cartID = newCart.CartID
	} else {
		cartID = cartInfo.CartID
	}

	var cartItemParam cartModel.CartItemParams
	cartItemParam.CartID = cartID
	if err := c.BindJSON(&cartItemParam); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	cartItem := &cartModel.CartItems{
		CartID:       cartItemParam.CartID,
		ItemID:       cartItemParam.ItemID,
		RestaurantID: cartItemParam.RestaurantID,
		Quantity:     cartItemParam.Quantity,
	}

	_, err = s.service.AddItem(ctx, cartItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Items added to cart!"})
}

// @Summary     List cart items
// @Description Returns all items in the authenticated user's cart
// @Tags        Cart
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/list [get]
func (s *CartHandler) getItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	cartInfo, err := s.service.GetCartId(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items, err := s.service.ListItems(ctx, cartInfo.CartID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// @Summary     Remove item from cart
// @Description Removes a specific item from the cart by cart item ID
// @Tags        Cart
// @Produce     json
// @Param       id path int true "Cart Item ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Security    BearerAuth
// @Router      /cart/remove/{id} [delete]
func (s *CartHandler) deleteItemFromCart(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cartItemId := c.Param("id")
	cartItemID, _ := strconv.ParseInt(cartItemId, 10, 64)
	_, err := s.service.DeleteItem(ctx, cartItemID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
