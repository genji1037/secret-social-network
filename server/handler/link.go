package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"secret-social-network/service"
)

func Link(c *gin.Context) {
	req := struct {
		UID1  string  `json:"uid1" binding:"required"`
		UID2  string  `json:"uid2" binding:"required"`
		Point float64 `json:"point" binding:"required"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	if err := service.Link(req.UID1, req.UID2, req.Point); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	c.JSON(http.StatusOK, "ok")

}
