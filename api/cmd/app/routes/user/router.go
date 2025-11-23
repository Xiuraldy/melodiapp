package users

import (
	"github.com/gin-gonic/gin"

	userapi "melodiapp/internal/adapters/api/user"
	dbadapter "melodiapp/internal/adapters/database/user"
	coreuser "melodiapp/internal/core/user"
)

func AddUserRoutes(r *gin.Engine) {
	group := r.Group("/users")

	repo := dbadapter.NewGormUserRepository()
	service := coreuser.NewService(repo)
	handlers := userapi.NewUserHandlers(service)

	group.GET("", handlers.GetAllUsers)
	group.GET("/me", handlers.GetMe)
	group.POST("", handlers.CreateUser)
	group.GET("/:id", handlers.GetUserById)
	group.DELETE("/:id", handlers.DeleteUser)
	group.PUT("/:id", handlers.EditUser)
}
