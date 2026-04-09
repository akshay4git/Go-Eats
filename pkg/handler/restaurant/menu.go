package restaurant

import (
	"context"
	"net/http"
	"strconv"
	"time"

	models "github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
	"github.com/gin-gonic/gin"
)

// @Summary     Add a menu item
// @Description Adds a new menu item to a restaurant
// @Tags        Restaurant
// @Accept      json
// @Produce     json
// @Param       menuItem body object{restaurant_id=int,name=string,description=string,price=number,category=string,available=boolean} true "Menu item details"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /restaurant/menu [post]
func (s *RestaurantHandler) addMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var menuItem models.MenuItem
	if err := c.BindJSON(&menuItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	menuObject, err := s.service.AddMenu(ctx, &menuItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else {
		s.service.UpdateMenuPhoto(ctx, menuObject)
	}
	c.JSON(http.StatusCreated, gin.H{"message": "New Menu Added!"})
}

// @Summary     List menu items
// @Description Returns all menu items, or filtered by restaurant_id if provided
// @Tags        Restaurant
// @Produce     json
// @Param       restaurant_id query int false "Restaurant ID (optional)"
// @Success     200 {array}  object
// @Failure     404 {object} map[string]string
// @Router      /restaurant/menu [get]
func (s *RestaurantHandler) listMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	restaurantId := c.Query("restaurant_id")
	if restaurantId == "" {
		results, err := s.service.ListAllMenus(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, results)
		return
	} else {
		restaurantID, err := strconv.ParseInt(restaurantId, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid RestaurantID"})
			return
		}
		results, err := s.service.ListMenus(ctx, restaurantID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No results found"})
			return
		}
		c.JSON(http.StatusOK, results)
		return
	}
}

// @Summary     Delete a menu item
// @Description Deletes a menu item by restaurant ID and menu ID
// @Tags        Restaurant
// @Produce     json
// @Param       restaurant_id path int true "Restaurant ID"
// @Param       menu_id       path int true "Menu ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Router      /restaurant/menu/{restaurant_id}/{menu_id} [delete]
func (s *RestaurantHandler) deleteMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	menuId, err := strconv.ParseInt(c.Param("menu_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid MenuID"})
		return
	}
	restaurantId, err := strconv.ParseInt(c.Param("restaurant_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid RestaurantID"})
		return
	}
	_, err = s.service.DeleteMenu(ctx, menuId, restaurantId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
