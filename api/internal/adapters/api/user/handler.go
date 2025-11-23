package userapi

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	userports "melodiapp/internal/ports/user"
	"melodiapp/models"
	"melodiapp/shared"
)

type UserHandlers struct {
	service userports.UserService
}

func NewUserHandlers(service userports.UserService) *UserHandlers {
	return &UserHandlers{service: service}
}

func (h *UserHandlers) GetAllUsers(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	_ = userData // authentication already validated, we don't use userData directly here

	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandlers) GetMe(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	if userData.ExpiryTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	user, err := h.service.GetUserByUintID(userData.Uid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) GetUserById(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	if userData.ExpiryTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	// ensure the session user still exists
	if _, err := h.service.GetUserByUintID(userData.Uid); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	created, err := h.service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *UserHandlers) EditUser(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	_, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")
	var updatedData models.User
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	updated, err := h.service.UpdateUser(id, &updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": updated})
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"id": id})
}
