package onePass

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BatchPay(c *gin.Context) {

}

func UserTrade(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func QueryUserAmount(c *gin.Context) {

}
