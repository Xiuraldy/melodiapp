package authapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authports "melodiapp/internal/ports/auth"
	"melodiapp/models"
	"melodiapp/shared"
)

type AuthHandlers struct {
	service authports.AuthService
}

func NewAuthHandlers(service authports.AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

func (h *AuthHandlers) Register(c *gin.Context) {
	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Register(input)
	if err != nil {
		switch err.Error() {
		case "Incomplete fields":
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case "Invalid email format":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "Email already exists":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	if tokenStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if err := h.service.Logout(tokenStr); err != nil {
		if err.Error() == "You don't have permission" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandlers) Login(c *gin.Context) {
	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	token, err := h.service.Login(input)
	if err != nil {
		switch err.Error() {
		case "Incomplete fields":
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case "Invalid email format":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "Invalid credentials":
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
