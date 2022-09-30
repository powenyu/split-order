package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	postgresql "github.com/powenyu/split-order/postgres"
)

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

func HeartBeat(c *gin.Context) {
	c.Status(http.StatusOK)
}
