package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/protos/api"
	"sync"
)

const (
	FacetsConsensusPrefix = "cons_"
	FacetsTotalKey        = "cons_total"
	FacetsValueSuffix     = "|value"
)

type User struct{}

var lock sync.Mutex

func (User) LinkOrAdd(appID, u1, u2 string, value float64) error {
	// apply lock FIXME: 给每个u1 u2 pair上细粒度的锁.consider database layer lock.
	lock.Lock()
	defer lock.Unlock()

	appValue, totalValue, err := QueryValue(appID, u1, u2)
	if err != nil {
		return fmt.Errorf("query value failed: %s", err.Error())
	}
	err = User{}.Link(appID, u1, u2, appValue+value, totalValue+value)
	if err != nil {
		return fmt.Errorf("link user failed: %s", err.Error())
	}
	return nil
}

// Link links two user with value. if two users already linked, over writes.
func (User) Link(appId, u1, u2 string, appValue, totalValue float64) error {
	appFacetsKey := FacetsConsensusPrefix + appId
	query := fmt.Sprintf(`
query {
	u1 as q1(func: eq(name, "%s")) {
		uid
	}
	u2 as q2(func: eq(name, "%s")) {
		uid
	}
}`, u1, u2)
	mus := []*api.Mutation{
		{
			Cond: `@if((not eq(len(u1),0)) and eq(len(u2),0))`, // u1 exist, u2 not exists.
			SetNquads: []byte(`
							_:u2 <name> "` + u2 + `" .
							uid(u1) <` + FacetsTotalKey + `> _:u2 (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							_:u2 <` + FacetsTotalKey + `> uid(u1) (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							uid(u1) <` + appFacetsKey + `> _:u2 (value=` + fmt.Sprintf("%.0f", appValue) + `) .
							_:u2 <` + appFacetsKey + `> uid(u1) (value=` + fmt.Sprintf("%.0f", appValue) + `) .`),
		},
		{
			Cond: `@if(eq(len(u1),0) and (not eq(len(u2), 0)))`, // u1 not exists, u2 exists.
			SetNquads: []byte(`
							_:u1 <name> "` + u1 + `" .
							_:u1 <` + FacetsTotalKey + `> uid(u2) (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							uid(u2) <` + FacetsTotalKey + `> _:u1 (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							_:u1 <` + appFacetsKey + `> uid(u2) (value=` + fmt.Sprintf("%.0f", appValue) + `) .
							uid(u2) <` + appFacetsKey + `> _:u1 (value=` + fmt.Sprintf("%.0f", appValue) + `) .`),
		},
		{
			Cond: `@if(eq(len(u1),0) and eq(len(u2),0))`, // both not exists,
			SetNquads: []byte(`
							_:u1 <name> "` + u1 + `" .
							_:u2 <name> "` + u2 + `" .
							_:u2 <` + FacetsTotalKey + `> _:u1 (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							_:u1 <` + FacetsTotalKey + `> _:u2 (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							_:u2 <` + appFacetsKey + `> _:u1 (value=` + fmt.Sprintf("%.0f", appValue) + `) .
							_:u1 <` + appFacetsKey + `> _:u2 (value=` + fmt.Sprintf("%.0f", appValue) + `) .`),
		},
		{
			Cond: `@if((not eq(len(u1),0)) and (not eq(len(u2),0)))`, // both exists
			SetNquads: []byte(`
							uid(u1) <` + FacetsTotalKey + `> uid(u2) (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							uid(u2) <` + FacetsTotalKey + `> uid(u1) (value=` + fmt.Sprintf("%.0f", totalValue) + `) .
							uid(u1) <` + appFacetsKey + `> uid(u2) (value=` + fmt.Sprintf("%.0f", appValue) + `) .
							uid(u2) <` + appFacetsKey + `> uid(u1) (value=` + fmt.Sprintf("%.0f", appValue) + `) .`),
		},
	}
	req := &api.Request{
		Query:     query,
		Mutations: mus,
		CommitNow: true,
	}
	if _, err := Dg.NewTxn().Do(context.Background(), req); err != nil {
		return err
	}

	return nil
}

// QueryValue return first params as
func QueryValue(appID, u1, u2 string) (appValue float64, totalValue float64, err error) {
	appIDFacets := FacetsConsensusPrefix + appID
	variables := make(map[string]string)

	variables["$name1"] = u1
	variables["$name2"] = u2

	ctx := context.Background()

	var q = `
		query Point($name1: string, $name2: string){
			users(func: eq(name, $name1)) {
				uid
				name
				` + FacetsTotalKey + ` @filter(eq(name, $name2)) @facets {
					uid
					name
				}
				` + appIDFacets + ` @filter(eq(name, $name2)) @facets {
					uid
					name
				}
			}
		}
`
	fmt.Println(q)
	resp, err := Dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		return
	}

	type Root struct {
		Users []map[string]interface{} `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return
	}

	if len(r.Users) == 0 {
		return
	}

	totalFacetsMap, ok := r.Users[0][FacetsTotalKey+FacetsValueSuffix].(map[string]interface{})
	if !ok {
		return
	}
	total, ok := totalFacetsMap["0"].(float64)
	if !ok {
		return
	}
	totalValue = total

	appFacetsMap, ok := r.Users[0][appIDFacets+FacetsValueSuffix].(map[string]interface{})
	if !ok {
		return
	}
	app, ok := appFacetsMap["0"].(float64)
	if !ok {
		return
	}
	appValue = app
	return
}
