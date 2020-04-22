package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// todo decouple logic and storage
func ListRelation() (map[string]float64, error) {
	const q = `
query data() {
	data(func: has(name)) {
		name
	    links @facets(point) {
			name
	    }
	}
}`
	ctx := context.Background()
	resp, err := Dg.NewTxn().Query(ctx, q)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", string(resp.Json))

	type Root struct {
		Users []UserResp `json:"data"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	relationValueMapper := make(map[string]float64)
	UserResp{}.Walk(r.Users, 0, func(u UserResp, depth int) {
		for i, child := range u.Links {
			k := fmt.Sprintf("%s-%s", u.Name, child.Name)
			v := u.Point[strconv.Itoa(i)]
			relationValueMapper[k] = v
		}
	})

	return relationValueMapper, nil
}
