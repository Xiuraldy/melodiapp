package serviceuserapi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"melodiapp/database"
	serviceuserports "melodiapp/internal/ports/serviceuser"
	"melodiapp/models"
	"melodiapp/shared"
)

type ServiceUserHandlers struct {
	service serviceuserports.ServiceUserService
}

type assignUsersRequest struct {
	UserIDs []uint `json:"user_ids"`
}

type changeStatusRequest struct {
	Status string `json:"status"`
}

func NewServiceUserHandlers(s serviceuserports.ServiceUserService) *ServiceUserHandlers {
	return &ServiceUserHandlers{service: s}
}

func getCurrentUser(c *gin.Context) (*models.User, error) {
	tokenStr := shared.GetTokenFromRequest(c)
	if tokenStr == "" {
		return nil, fmt.Errorf("invalid token")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(*shared.Payload)
	session, exists := shared.Sessions[claims.Session]
	if !exists || session.ExpiryTime.Before(time.Now()) {
		return nil, fmt.Errorf("You don't have permission")
	}

	var user models.User
	if err := database.DBConn.First(&user, session.Uid).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *ServiceUserHandlers) AssignUsers(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
		return
	}

	serviceIDParam := c.Param("id")
	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}

	var req assignUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := h.service.AssignUsers(uint(serviceID64), req.UserIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"service_id": serviceIDParam, "user_ids": req.UserIDs})
}

func (h *ServiceUserHandlers) ListByService(c *gin.Context) {
	if _, err := getCurrentUser(c); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	serviceIDParam := c.Param("id")
	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}

	items, err := h.service.ListByService(uint(serviceID64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *ServiceUserHandlers) ChangeStatus(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	serviceIDParam := c.Param("id")
	userIDParam := c.Param("userId")

	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}
	userID64, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	if uint(userID64) != user.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
		return
	}

	var req changeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := h.service.ChangeStatus(uint(serviceID64), uint(userID64), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"service_id": serviceIDParam, "user_id": userIDParam, "status": req.Status})
}
