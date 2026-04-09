package restaurant

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	restaurantModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
	"github.com/gin-gonic/gin"
)

// @Summary     Add a restaurant
// @Description Creates a new restaurant with photo upload (multipart/form-data)
// @Tags        Restaurant
// @Accept      mpfd
// @Produce     json
// @Param       file        formData file   true  "Restaurant photo"
// @Param       name        formData string true  "Restaurant name"
// @Param       description formData string true  "Description"
// @Param       address     formData string true  "Address"
// @Param       city        formData string true  "City"
// @Param       state       formData string true  "State"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /restaurant/ [post]
func (s *RestaurantHandler) addRestaurant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	_ = ctx

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	originalFileName := fileHeader.Filename
	newFileName := generateFileName(originalFileName)
	_, err = s.Serve.Storage.Upload(newFileName, file)
	if err != nil {
		slog.Error("Error", "addRestaurant", err.Error())
	}

	uploadedFile := filepath.Join(os.Getenv("STORAGE_DIRECTORY"), newFileName)

	var restaurant restaurantModel.Restaurant
	restaurant.Name = c.PostForm("name")
	restaurant.Description = c.PostForm("description")
	restaurant.Address = c.PostForm("address")
	restaurant.City = c.PostForm("city")
	restaurant.State = c.PostForm("state")
	restaurant.Photo = uploadedFile

	_, err = s.service.Add(ctx, &restaurant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
}

// @Summary     List all restaurants
// @Description Returns a list of all restaurants
// @Tags        Restaurant
// @Produce     json
// @Success     200 {array}  object
// @Failure     404 {object} map[string]string
// @Router      /restaurant/ [get]
func (s *RestaurantHandler) listRestaurants(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	results, err := s.service.ListRestaurants(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if results == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No restaurants found"})
		return
	}
	c.JSON(http.StatusOK, results)
}

// @Summary     Get restaurant by ID
// @Description Returns a single restaurant by its ID
// @Tags        Restaurant
// @Produce     json
// @Param       id path int true "Restaurant ID"
// @Success     200 {array} object
// @Failure     404 {object} map[string]string
// @Router      /restaurant/{id} [get]
func (s *RestaurantHandler) listRestaurantById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	restaurantId := c.Param("id")
	restaurantID, _ := strconv.ParseInt(restaurantId, 10, 64)

	result, err := s.service.ListRestaurantById(ctx, restaurantID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// @Summary     Delete a restaurant
// @Description Deletes a restaurant by ID
// @Tags        Restaurant
// @Produce     json
// @Param       id path int true "Restaurant ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Router      /restaurant/{id} [delete]
func (s *RestaurantHandler) deleteRestaurant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	restaurantId := c.Param("id")
	restaurantID, _ := strconv.ParseInt(restaurantId, 10, 64)

	_, err := s.service.DeleteRestaurant(ctx, restaurantID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
