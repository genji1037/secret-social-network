package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/protos/api"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	facetsPrefix   = "cons|"
	facetsTotalKey = "total"
)

// LinkOrAdd links two user with its value,
// or add link value when two user are already linked.
func LinkOrAdd(appID, u1, u2 string, value float64) error {

	commitCount := 0

	variables := make(map[string]string)

	variables["$name1"] = u1
	variables["$name2"] = u2

	ctx := context.Background()

	var q = `
		query Point($name1: string, $name2: string){
			users(func: eq(name, $name1)) {
				uid
				name
				cons @filter(eq(name, $name2)) @facets {
					uid
					name
				}
			}
		}
`
top:
	txn := Dg.NewTxn()
	defer txn.Discard(ctx)
	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return err
	}

	type Root struct {
		Users []map[string]interface{} `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return err
	}

	if len(r.Users) == 0 {
		return err
	}

	totalFacetsMap, ok := r.Users[0][facetsPrefix+facetsTotalKey].(map[string]interface{})
	if !ok {
		return err
	}
	total, ok := totalFacetsMap["0"].(float64)
	if !ok {
		return err
	}
	totalValue := total + value

	appFacetsMap, ok := r.Users[0][facetsPrefix+appID].(map[string]interface{})
	if !ok {
		return fmt.Errorf("users[0] is not a map")
	}
	app, ok := appFacetsMap["0"].(float64)
	if !ok {
		return fmt.Errorf("appFacetsMap[0] is not float64")
	}
	appValue := app + value

	query := fmt.Sprintf(`
query {
	u1 as q1(func: eq(name, "%s")) {
		uid
	}
	u2 as q2(func: eq(name, "%s")) {
		uid
	}
}`, u1, u2)
	appV := fmt.Sprintf("%f", appValue)
	totalV := fmt.Sprintf("%f", totalValue)
	facets := facetsTotalKey + `=` + totalV + `, ` + appID + `=` + appV
	mus := []*api.Mutation{
		{
			Cond: `@if((not eq(len(u1),0)) and eq(len(u2),0))`, // u1 exist, u2 not exists.
			SetNquads: []byte(`
							_:u2 <name> "` + u2 + `" .
							uid(u1) <cons> _:u2 (` + facets + `) .
							_:u2 <cons> uid(u1) (` + facets + `) .`),
		},
		{
			Cond: `@if(eq(len(u1),0) and (not eq(len(u2), 0)))`, // u1 not exists, u2 exists.
			SetNquads: []byte(`
							_:u1 <name> "` + u1 + `" .
							_:u1 <cons> uid(u2) (` + facets + `) .
							uid(u2) <cons> _:u1 (` + facets + `) .`),
		},
		{
			Cond: `@if(eq(len(u1),0) and eq(len(u2),0))`, // both not exists,
			SetNquads: []byte(`
							_:u1 <name> "` + u1 + `" .
							_:u2 <name> "` + u2 + `" .
							_:u2 <cons> _:u1 (` + facets + `) .
							_:u1 <cons> _:u2 (` + facets + `) .`),
		},
		{
			Cond: `@if((not eq(len(u1),0)) and (not eq(len(u2),0)))`, // both exists
			SetNquads: []byte(`
							uid(u1) <cons> uid(u2) (` + facets + `) .
							uid(u2) <cons> uid(u1) (` + facets + `) .`),
		},
	}
	req := &api.Request{
		Query:     query,
		Mutations: mus,
	}
	time.Sleep(5 * time.Second)
	if _, err := txn.Do(ctx, req); err != nil {
		return err
	}
	err = txn.Commit(ctx)
	if err != nil && commitCount < 3 {
		commitCount++
		log.WithFields(log.Fields{
			"commit_count": commitCount,
			"reason":       err.Error(),
		}).Warnf("link or add commit failed, retrying...")
		goto top
	}
	return err
}
