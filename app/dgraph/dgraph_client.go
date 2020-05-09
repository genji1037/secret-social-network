package dgraph

import (
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"secret-social-network/app/config"
)

// Dg is dGraph client.
var Dg *dgo.Dgraph

// Open connecting to dGraph.
func Open(cfg config.DGraph) error {
	conn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	dc := api.NewDgraphClient(conn)
	Dg = dgo.NewDgraphClient(dc)
	// TODO: ⭐ doing something like ping, in order to confirm correctly connected to dGraph.
	return nil
}
