package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"secret-social-network/service"
)

func FindAll(c *gin.Context) {
	name1 := c.Query("name1")
	name2 := c.Query("name2")

	result, err := service.List(name1, name2)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
