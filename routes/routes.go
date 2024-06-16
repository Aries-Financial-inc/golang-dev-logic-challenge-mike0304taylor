package routes

import (
	_ "golang-dev-logic-challenge/docs"
	"golang-dev-logic-challenge/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/analysis", handlers.AnalysisHandler)
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
