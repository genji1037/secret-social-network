package handler

import (
	"context"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"secret-social-network/app/dgraph"
)

func Init(c *gin.Context) {

	op := api.Operation{DropAll: true}
	ctx := context.Background()
	if err := dgraph.Dg.Alter(ctx, &op); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	op = api.Operation{}
	op.Schema = `
		name: string @index(term) .
	`

	err := dgraph.Dg.Alter(ctx, &op)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}

	c.JSON(http.StatusOK, "ok")
}