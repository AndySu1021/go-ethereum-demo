package routes

import (
	"github.com/gin-gonic/gin"
	"go-ethereum-demo/controllers"
)

func SetRouter() *gin.Engine  {
	router := gin.Default()
	router.TrustedPlatform = "True-Client-IP"
	routes := router.Group("api")
	{
		routes.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		blocksGroup := routes.Group("blocks")
		{
			blocksGroup.GET("", controllers.GetBlocksList)
			blocksGroup.GET("/:id", controllers.GetBlockDetail)
		}

		transactionGroup := routes.Group("transaction")
		{
			transactionGroup.GET("/:txHash", controllers.GetTransactionDetail)
		}
	}

	return router
}
