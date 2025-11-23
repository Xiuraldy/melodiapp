package service

import (
	"github.com/gin-gonic/gin"

	serviceapi "melodiapp/internal/adapters/api/service"
	serviceuserapi "melodiapp/internal/adapters/api/serviceuser"
	servicesongapi "melodiapp/internal/adapters/api/servicesong"
	dbadapter "melodiapp/internal/adapters/database/service"
	dbserviceuser "melodiapp/internal/adapters/database/serviceuser"
	dbservicesong "melodiapp/internal/adapters/database/servicesong"
	coreservice "melodiapp/internal/core/service"
	coreserviceuser "melodiapp/internal/core/serviceuser"
	coreservicesong "melodiapp/internal/core/servicesong"
)

func AddServiceRoutes(r *gin.Engine) {
	group := r.Group("/services")

	serviceRepo := dbadapter.NewGormServiceRepository()
	serviceUsecase := coreservice.NewServiceUsecase(serviceRepo)
	serviceHandlers := serviceapi.NewServiceHandlers(serviceUsecase)

	serviceUserRepo := dbserviceuser.NewGormServiceUserRepository()
	serviceUserUsecase := coreserviceuser.NewService(serviceUserRepo)
	serviceUserHandlers := serviceuserapi.NewServiceUserHandlers(serviceUserUsecase)

	serviceSongRepo := dbservicesong.NewGormServiceSongRepository()
	serviceSongUsecase := coreservicesong.NewService(serviceSongRepo)
	serviceSongHandlers := servicesongapi.NewServiceSongHandlers(serviceSongUsecase)

	group.GET("", serviceHandlers.GetAll)
	group.GET(":id", serviceHandlers.GetByID)
	group.POST("", serviceHandlers.Create)
	group.PUT(":id", serviceHandlers.Update)
	group.DELETE(":id", serviceHandlers.Delete)

	group.POST(":id/users", serviceUserHandlers.AssignUsers)
	group.GET(":id/users", serviceUserHandlers.ListByService)
	group.PATCH(":id/users/:userId/status", serviceUserHandlers.ChangeStatus)

	group.POST(":id/songs", serviceSongHandlers.AssignSongs)
	group.GET(":id/songs", serviceSongHandlers.ListByService)
	group.DELETE(":id/songs/:songId", serviceSongHandlers.Remove)
}
