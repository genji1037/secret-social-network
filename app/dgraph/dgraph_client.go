package dgraph

import (
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

var Dg *dgo.Dgraph

func Open() error {
	conn, err := grpc.Dial("127.0.0.1:9081", grpc.WithInsecure())
	if err != nil {
		return err
	}

	dc := api.NewDgraphClient(conn)
	Dg = dgo.NewDgraphClient(dc)
	return nil
}
