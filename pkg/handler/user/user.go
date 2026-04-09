package user

import (
	"context"
	"net/http"
	"strconv"
	"time"

	userModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/user"
	userService "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
)

// @Summary     Register a new user
// @Description Creates a new user account
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       user body object{name=string,email=string,password=string} true "User registration details"
// @Success     201 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /user/ [post]
func (s *UserHandler) addUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var u userModel.User
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := s.validate.Struct(u); err != nil {
		validationError := userModel.UserValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError})
		return
	}
	_, err := s.service.Add(ctx, &u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// @Summary     Delete a user
// @Description Deletes a user by ID
// @Tags        User
// @Produce     json
// @Param       id path int true "User ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Security    BearerAuth
// @Router      /user/{id} [delete]
func (s *UserHandler) deleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	idParam := c.Param("id")
	userID, _ := strconv.ParseInt(idParam, 10, 64)
	_, err := s.service.Delete(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary     Login
// @Description Authenticates a user and returns a JWT token
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       credentials body object{email=string,password=string} true "Login credentials"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /user/login [post]
func (s *UserHandler) loginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var creds userModel.LoginUser
	if err := c.BindJSON(&creds); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	login := userService.ValidateAccount(
		s.service.Login,
		s.service.UserExist,
		s.service.ValidatePassword,
	)
	token, err := login(ctx, &userModel.LoginUser{
		Email:    creds.Email,
		Password: creds.Password,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
