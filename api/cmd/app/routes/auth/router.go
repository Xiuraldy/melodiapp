package auth

import (
	"github.com/gin-gonic/gin"

	authapi "melodiapp/internal/adapters/api/auth"
	dbadapter "melodiapp/internal/adapters/database/user"
	authcore "melodiapp/internal/core/auth"
)

func AddAuthRoutes(r *gin.Engine) {
	group := r.Group("/auth")

	userRepo := dbadapter.NewGormUserRepository()
	service := authcore.NewService(userRepo)
	handlers := authapi.NewAuthHandlers(service)

	group.POST("/register", handlers.Register)
	group.POST("/login", handlers.Login)
	group.DELETE("/logout", handlers.Logout)
}
