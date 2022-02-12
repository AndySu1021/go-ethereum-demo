package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ethereum-demo/services"
	"net/http"
)

func GetTransactionDetail (c *gin.Context) {
	txRecord, err := services.GetTransactionByHash(c.Param("txHash"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
	}

	c.JSON(http.StatusOK, txRecord)
}
