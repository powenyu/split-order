package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/powenyu/split-order/bot"
	postgresql "github.com/powenyu/split-order/postgres"
)

func Start(c *gin.Context) {
	bot.Start()
}

func Dbtest(c *gin.Context) {
	sql := `SELECT 1`
	var t int
	err := postgresql.Pool.QueryRow(c, sql).Scan(&t)
	if err != nil {
		fmt.Println("error : ", err)
		return
	}
	fmt.Println("t : ", t)
}
