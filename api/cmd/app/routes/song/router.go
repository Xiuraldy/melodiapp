package song

import (
	"github.com/gin-gonic/gin"

	songapi "melodiapp/internal/adapters/api/song"
	dbadapter "melodiapp/internal/adapters/database/song"
	coresong "melodiapp/internal/core/song"
)

func AddSongRoutes(r *gin.Engine) {
	group := r.Group("/songs")

	repo := dbadapter.NewGormSongRepository()
	service := coresong.NewService(repo)
	handlers := songapi.NewSongHandlers(service)

	group.GET("", handlers.GetAll)
	group.GET(":id", handlers.GetByID)
	group.POST("", handlers.Create)
	group.PUT(":id", handlers.Update)
	group.DELETE(":id", handlers.Delete)
}
