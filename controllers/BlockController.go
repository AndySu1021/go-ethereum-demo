package controllers

import (
	"github.com/gin-gonic/gin"
	"go-ethereum-demo/requests"
	"go-ethereum-demo/services"
	"net/http"
)

func GetBlocksList (c *gin.Context) {
	var params requests.GetBlocksListParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error":err.Error()})
		return
	}

	blocksList, err := services.GetBlocksList(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocks": blocksList})
	return
}

func GetBlockDetail (c *gin.Context) {
	block, err := services.GetBlockById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, block)
}
