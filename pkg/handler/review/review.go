package review

import (
	"context"
	"net/http"
	"strconv"
	"time"

	reviewModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/review"
	"github.com/gin-gonic/gin"
)

// @Summary     Add a review
// @Description Adds a review for a restaurant by the authenticated user
// @Tags        Reviews
// @Accept      json
// @Produce     json
// @Param       restaurant_id path int                  true "Restaurant ID"
// @Param       review body review.ReviewParams true "Review details"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /review/{restaurant_id} [post]
func (s *ReviewProtectedHandler) addReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetInt64("userID")
	restaurantId, err := strconv.ParseInt(c.Param("restaurant_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid RestaurantID"})
		return
	}

	var reviewParam reviewModel.ReviewParams
	var review reviewModel.Review
	if err := c.BindJSON(&reviewParam); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := s.validate.Struct(reviewParam); err != nil {
		validationError := reviewModel.ReviewValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError})
		return
	}

	review.UserID = userID
	review.RestaurantID = restaurantId
	review.Rating = reviewParam.Rating
	review.Comment = reviewParam.Comment

	_, err = s.service.Add(ctx, &review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Review Added!"})
}

// @Summary     List reviews
// @Description Returns all reviews for a specific restaurant
// @Tags        Reviews
// @Produce     json
// @Param       restaurant_id path int true "Restaurant ID"
// @Success     200 {array}  object
// @Failure     404 {object} map[string]string
// @Security    BearerAuth
// @Router      /review/{restaurant_id} [get]
func (s *ReviewProtectedHandler) listReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	restaurantId, err := strconv.ParseInt(c.Param("restaurant_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Invalid RestaurantID"})
		return
	}
	results, err := s.service.ListReviews(ctx, restaurantId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No results found"})
		return
	}
	c.JSON(http.StatusOK, results)
}

// @Summary     Delete a review
// @Description Deletes a review by ID for the authenticated user
// @Tags        Reviews
// @Produce     json
// @Param       review_id path int true "Review ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Security    BearerAuth
// @Router      /review/{review_id} [delete]
func (s *ReviewProtectedHandler) deleteReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	reviewId := c.Param("review_id")
	userID := c.GetInt64("userID")
	reviewID, _ := strconv.ParseInt(reviewId, 10, 64)

	_, err := s.service.DeleteReview(ctx, reviewID, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
