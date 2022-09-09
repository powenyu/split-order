package routes

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/powenyu/split-order/bot"
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

	router.GET("/", Start)

	return router
}

func Start(c *gin.Context) {
	bot.Start()
	fmt.Println("Hello World")
}
