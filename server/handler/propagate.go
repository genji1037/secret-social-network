package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"secret-social-network/memgraph"
	"secret-social-network/service"
	"strings"
	"time"
)

func Propagate(c *gin.Context) {
	db := memgraph.NewMemDB()
	// list relation network
	relations, err := service.ListRelationWithValue()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	// prepare mem data
	startAt := time.Now()
	for relation, value := range relations {
		tmp := strings.Split(relation, "-")
		db.Link(tmp[0], tmp[1], value)
	}
	fmt.Printf("[link] %s\n", time.Now().Sub(startAt))

	// propagate
	db.PropagateAll()
	c.JSON(200, "ok")
}
