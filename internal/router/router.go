package router

import (
	"stockfoilo_test/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	}))

	router.POST("/upload", controller.UploadVideo)
	router.POST("/modify", controller.ModifyVideo)
	router.GET("/video", controller.GetVideoInfoList)
}
