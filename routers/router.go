package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/huiming23344/balanceapi/routers/api/onePass"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	onePathApi := r.Group("/onePass")
	onePathApi.Use()
	{
		onePathApi.POST("/batchPay", onePass.BatchPay)
		onePathApi.POST("/userTrade", onePass.UserTrade)
		onePathApi.POST("/queryUserAmount", onePass.QueryUserAmount)
	}

	return r
}
