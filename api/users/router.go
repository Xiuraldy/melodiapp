package users

import "github.com/gin-gonic/gin"

func AddUserRoutes(r *gin.Engine) {
	group := r.Group("/users")

	group.GET("", GetAllUsers)
	group.GET("/me", GetMe)
	group.POST("", CreateUser)
	group.GET("/:id", GetUserById)
	group.DELETE("/:id", DeleteUser)
	group.PUT("/:id", EditUser)
}
