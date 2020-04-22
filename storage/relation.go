package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// todo decouple logic and storage
func ListRelation() (map[string]float64, error) {
	first := 5000
	offset := 0
	relationValueMapper := make(map[string]float64)

	for {
		var q = fmt.Sprintf(`
query data() {
	data(func: has(name), first:%d, offset:%d) {
		name
	    links @facets(point) {
			name
	    }
	}
}`, first, offset)
		offset += first
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
		if len(r.Users) == 0 {
			break
		}

		UserResp{}.Walk(r.Users, 0, func(u UserResp, depth int) {
			for i, child := range u.Links {
				k := fmt.Sprintf("%s-%s", u.Name, child.Name)
				v := u.Point[strconv.Itoa(i)]
				relationValueMapper[k] = v
			}
		})
	}

	return relationValueMapper, nil
}
