package dgraph

import (
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"log"
)

var Dg *dgo.Dgraph

func init() {
	conn, err := grpc.Dial("127.0.0.1:9081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	dc := api.NewDgraphClient(conn)
	Dg = dgo.NewDgraphClient(dc)
}
