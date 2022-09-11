package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/powenyu/split-order/controllers/v1"
)

// InitRouter init router
func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	if err := config.Validate(); err != nil {
		panic(err)
	}
	router.Use(cors.New(config))

	router.GET("/", v1.Start)

	apiv1 := router.Group("/api/v1")
	apiv1.GET("/test", v1.Dbtest)
	return router
}
