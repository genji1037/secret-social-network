package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"secret-social-network/storage"
	"strconv"
	"strings"
)

// List list.
func List(uid1, uid2 string) ([]storage.UserResp, map[int]float64, error) {
	pointBal, err := storage.QueryPoint(uid1, uid2)
	if err != nil {
		return nil, nil, err
	}

	return TreeDisableLoop(uid1, uid2, pointBal)
}

// 拉取关系树（禁止循环）
func TreeDisableLoop(uid1, uid2 string, valueBal float64) ([]storage.UserResp, map[int]float64, error) {
	var q = `
query All($name: string){
	u0 as data(func: anyofterms(name, $name)) {
		uid
		name 
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
	if depth > 1 {
		q = strings.ReplaceAll(q, fmt.Sprintf("u%d as ", depth-1), "")
	}

	q = fmt.Sprintf(q, "")

	ctx := context.Background()
	variables := make(map[string]string)
	variables["$name"] = uid1 + " " + uid2
	resp, err := storage.Dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, nil, err
	}

	result := struct {
		Data []storage.UserResp
	}{}
	if err := json.Unmarshal(resp.Json, &result); err != nil {
		return nil, nil, err
	}

	return result.Data, depthValueMapper, nil
}

func embedLink(depth int, input string, point float64) string {
	point = math.Floor(point)
	// TODO: 同一层可能会出现多个相同节点，去掉相同节点，减少计算量
	return fmt.Sprintf(input, fmt.Sprintf(`
		%s links @facets(gt(point, %.0f)) %s {
			name
			uid
			%s
		}`, fmt.Sprintf("u%d as", depth), point, genFilter(depth), "%s"))
}

func genFilter(depth int) string {
	if depth > 1 {
		depth--
	}
	filter := "@filter(not ("
	first := true
	for i := 0; i < depth; i++ {
		if first {
			first = false
		} else {
			filter += " or "
		}
		filter += "uid(u"
		filter += strconv.Itoa(i)
		filter += ")"
	}
	filter += "))"
	return filter
}
