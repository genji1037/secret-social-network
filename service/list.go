package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"secret-social-network/dgraph"
	"strconv"
	"strings"
)

// List list.
func List(uid1, uid2 string) ([]dgraph.UserResp, map[int]float64, error) {
	pointBal, err := dgraph.QueryPoint(uid1, uid2)
	if err != nil {
		return nil, nil, err
	}

	return TreeDisableLoop(uid1, uid2, pointBal)
}

// 拉取关系树（禁止循环）
func TreeDisableLoop(uid1, uid2 string, valueBal float64) ([]dgraph.UserResp, map[int]float64, error) {
	var q = `
query All($name: string){
	u0 as data(func: anyofterms(name, $name)) {
		n:name 
		%s
	}
}
	`
	valueBal = valueBal * valueRatio
	depth := 0
	depthValueMapper := make(map[int]float64)
	for valueBal > 10 {
		depth++
		q = embedLink(depth, q, valueBal)
		depthValueMapper[depth] = valueBal
		valueBal = valueBal * valueRatio
	}
	// remove last 2 unused variable
	q = strings.ReplaceAll(q, fmt.Sprintf("u%d as ", depth), "")

	q = fmt.Sprintf(q, "")

	ctx := context.Background()
	variables := make(map[string]string)
	variables["$name"] = uid1 + " " + uid2

	resp, err := dgraph.Dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, nil, fmt.Errorf("%s-%s%s: %s", uid1, uid2, q, err.Error())
	}

	result := struct {
		Data []dgraph.UserResp
	}{}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, nil, err
	}

	return result.Data, depthValueMapper, nil
}

func embedLink(depth int, input string, point float64) string {
	point = math.Floor(point)
	// TODO: 同一层可能会出现多个相同节点，去掉相同节点，减少计算量
	return fmt.Sprintf(input, fmt.Sprintf(`		%s l:links @facets(gt(point, %.0f)) %s {
			n:name
			%s
		}`, fmt.Sprintf("u%d as", depth), point, genFilter(depth), "%s"))
}

func genFilter(depth int) string {
	filter := "@filter(not (uid("
	first := true
	for i := 0; i < depth; i++ {
		if first {
			first = false
		} else {
			filter += ","
		}
		filter += "u"
		filter += strconv.Itoa(i)
	}
	filter += ")))"
	return filter
}
