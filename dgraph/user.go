package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/protos/api"
)

type User struct {
	Uid   string `json:"uid"`
	Name  string `json:"name,omitempty"`
	Links []User `json:"links,omitempty"`

	Point float64 `json:"links|point,omitempty"`
}

// Link links two user with point
func (User) Link(u1, u2 string, point float64) error {
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
			Cond: `@if((not eq(len(u1),0)) and eq(len(u2),0))`,
			SetNquads: []byte(`
							_:u2 <name> "` + u2 + `" .
							uid(u1) <links> _:u2 (point=` + fmt.Sprintf("%.0f", point) + `) .
							_:u2 <links> uid(u1) (point=` + fmt.Sprintf("%.0f", point) + `) .`),
		},
		{
			Cond: `@if(eq(len(u1),0) and (not eq(len(u2), 0)))`, // if u2 not exists
			SetNquads: []byte(`
							_:u1 <name> "` + u1 + `" .
							_:u1 <links> uid(u2) (point=` + fmt.Sprintf("%.0f", point) + `) .
							uid(u2) <links> _:u1 (point=` + fmt.Sprintf("%.0f", point) + `) .`),
		},
		{
			Cond: `@if(eq(len(u1),0) and eq(len(u2),0))`,
			SetNquads: []byte(`
							_:u1 <name> "` + u1 + `" .
							_:u2 <name> "` + u2 + `" .
							_:u2 <links> _:u1 (point=` + fmt.Sprintf("%.0f", point) + `) .
							_:u1 <links> _:u2 (point=` + fmt.Sprintf("%.0f", point) + `) .`),
		},
		{
			Cond: `@if((not eq(len(u1),0)) and (not eq(len(u2),0)))`,
			SetNquads: []byte(`
							uid(u1) <links> uid(u2) (point=` + fmt.Sprintf("%.0f", point) + `) .
							uid(u2) <links> uid(u1) (point=` + fmt.Sprintf("%.0f", point) + `) .`),
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

func QueryPoint(name1, name2 string) (float64, error) {
	variables := make(map[string]string)

	variables["$name1"] = name1
	variables["$name2"] = name2

	ctx := context.Background()

	const q = `
		query Point($name1: string, $name2: string){
			point(func: eq(name, $name1)) {
				uid
				name
				links @filter(eq(name, $name2)) @facets {
					uid
					name
				}
			}
		}
`
	resp, err := Dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		return 0, err
	}

	type Root struct {
		Point []*UserResp `json:"point"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return 0, err
	}

	if len(r.Point) == 0 {
		return 0, nil
	}

	return r.Point[0].Point["0"], nil
}
