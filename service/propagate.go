package service

import (
	"fmt"
	"secret-social-network/storage"
	"secret-social-network/util"
	"strings"
	"sync"
	"time"
)

type result struct {
	assignmentCount int
	resultHashes    map[string]float64
	sync.Mutex
}

func (r *result) add(subResultHashes map[string]float64) {
	r.Lock()
	r.assignmentCount += len(subResultHashes)
	for k, v := range subResultHashes {
		r.resultHashes[k] += v
	}
	r.Unlock()
}

func propagate() (map[string]float64, error) {
	res := result{
		assignmentCount: 0,
		resultHashes:    make(map[string]float64),
	}
	startAt := time.Now()
	originValues, err := listRelationWithValue()
	if err != nil {
		return nil, err
	}

	type workInput struct {
		k string
		v float64
	}
	agent := util.NewAgent(4, 1024, func(in interface{}) interface{} {
		kv := in.(workInput)
		// 从一条关系开始拉取关系树
		tree, depthHashMapper := listRelationTree(strings.Split(kv.k, "-"), kv.v)
		if len(tree) > 0 {
			// 给每个节点分配算力
			subResultHashes := assignHash(tree, depthHashMapper)
			return subResultHashes
		}
		return make(map[string]float64)
	})

	wg := sync.WaitGroup{}
	for k, v := range originValues {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := agent.Do(workInput{
				k: k,
				v: v,
			})
			out := <-c
			res.add(out.(map[string]float64))
		}()
	}
	wg.Wait()
	fmt.Printf("[propagate] %d relations %d assignment cost %s\n", len(originValues), res.assignmentCount, time.Now().Sub(startAt))
	return res.resultHashes, nil
}

func listRelationWithValue() (map[string]float64, error) {
	startAt := time.Now()
	rs, err := storage.ListRelation()
	fmt.Printf("[listRelationWithValue] %s\n", time.Now().Sub(startAt))
	return rs, err
}

func listRelationTree(us []string, value float64) ([]storage.UserResp, map[int]float64) {
	startAt := time.Now()
	if len(us) != 2 {
		fmt.Printf("[ERROR] list relation tree by us[%+v] failed: \n", us)
		return nil, nil
	}
	tree, depthHashMapper, err := TreeDisableLoop(us[0], us[1], value)
	if err != nil {
		fmt.Printf("[ERROR] List failed: %s\n", err.Error())
		return nil, nil
	}
	fmt.Printf("[tree] %s\n", time.Now().Sub(startAt))
	return tree, depthHashMapper
}

// 分配算力
func assignHash(tree []storage.UserResp, depthHashMapper map[int]float64) map[string]float64 {
	startAt := time.Now()
	// 根据树的深度 原始价值 分配算力到values 和 resultHashes里
	resultHashes := make(map[string]float64)
	storage.UserResp{}.Walk(tree, 0, func(u storage.UserResp, depth int) {
		if depth > 0 {
			hash := depthHashMapper[depth]
			resultHashes[u.Name] = hash
		}
	})
	fmt.Printf("[assi] %s\n", time.Now().Sub(startAt))
	return resultHashes
}
