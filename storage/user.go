package storage

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/protos/api"
)

type User struct {
	Uid   string `json:"uid"`
	Name  string `json:"name,omitempty"`
	Links []User `json:"links,omitempty"`

	Point float64 `json:"links|point,omitempty"`
}

func (u *User) Save() error {
	ctx := context.Background()
	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(u)
	if err != nil {
		return err
	}

	mu.SetJson = pb

	_, err = Dg.NewTxn().Mutate(ctx, mu)

	err = u.QueryByName()

	return err
}

func (u *User) QueryByName() error {
	const q = `
		query Me($name: string){
			data(func: eq(name, $name)) {
				uid
				name
			}
		}
	`
	variables := make(map[string]string)
	variables["$name"] = u.Name
	ctx := context.Background()
	resp, err := Dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		return err
	}

	type Root struct {
		Users []*User `json:"data"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return err
	}

	if len(r.Users) == 0 {
		return nil
	}

	u.Uid = r.Users[0].Uid
	u.Name = r.Users[0].Name
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
