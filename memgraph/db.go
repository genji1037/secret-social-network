package memgraph

import (
	"fmt"
	"secret-social-network/model"
	"secret-social-network/util"
	"sync"
	"time"
)

type Link struct {
	UID1  string
	UID2  string
	Value float64
}

type db struct {
	nodes map[string]*Node
	links []Link
}

func NewMemDB() *db {
	return &db{
		nodes: make(map[string]*Node),
		links: make([]Link, 0),
	}
}

func (db *db) Link(uid1, uid2 string, value float64) {
	db.links = append(db.links, Link{
		UID1:  uid1,
		UID2:  uid2,
		Value: value,
	})
	u1, ok1 := db.nodes[uid1]
	u2, ok2 := db.nodes[uid2]
	if !ok1 {
		u1 = &Node{
			Name:      uid1,
			Consensus: make([]*Node, 0, 1),
			Values:    make([]float64, 0, 1),
		}
		db.nodes[uid1] = u1
	}
	if !ok2 {
		u2 = &Node{
			Name:      uid2,
			Consensus: make([]*Node, 0, 1),
			Values:    make([]float64, 0, 1),
		}
		db.nodes[uid2] = u2
	}
	u1.Consensus = append(u1.Consensus, u2)
	u1.Values = append(u1.Values, value)
	u2.Consensus = append(u2.Consensus, u1)
	u2.Values = append(u2.Values, value)
}

func (db *db) PropagateAll() map[string]float64 {
	startAt := time.Now()
	result := model.Result{
		ResultHashes: make(map[string]float64),
	}

	agent := util.NewAgent(4, 0, func(in interface{}) interface{} {
		link := in.(Link)
		subResult, err := db.Propagate(link.UID1, link.UID2, link.Value)
		if err != nil {
			fmt.Printf("propagate failed: %s\n", err.Error())
		}
		return subResult
	})

	wg := sync.WaitGroup{}
	done := 0
	allWork := len(db.links)
	go watchProgress(&done, allWork)
	for i := range db.links {
		wg.Add(1)
		c := agent.Do(db.links[i])
		go func(c chan interface{}) {
			defer wg.Done()
			out := <-c
			result.Add(out.(map[string]float64))
			done++
		}(c)
	}

	wg.Wait()
	elapsed := time.Now().Sub(startAt)
	fmt.Printf("[propagate] %d nodes %d links %d assign cost %s\n", len(db.nodes), len(db.links), result.AssignmentCount, elapsed)
	return result.ResultHashes
}

func watchProgress(inProgress *int, all int) {
	for {
		time.Sleep(20 * time.Second)
		fmt.Printf("%d/%d\n", *inProgress, all)
	}
}

func (db *db) Propagate(uid1, uid2 string, value float64) (map[string]float64, error) {
	//startAt := time.Now()

	result := make(map[string]float64)

	assignHashRate([]*Node{db.nodes[uid1], db.nodes[uid2]}, value, result)

	//fmt.Printf("[prop] %s\n", time.Now().Sub(startAt))
	return result, nil
}

func assignHashRate(nodes []*Node, value float64, result map[string]float64) {
	subValue := value * AttenuationRate
	if subValue < MinValue { // 价值衰减到低于最小值 10
		return
	}
	subNodes := make([]*Node, 0)
	for _, node := range nodes {
		_, ok := result[node.Name]
		if !ok {
			result[node.Name] += value
			for i, subNode := range node.Consensus {
				if subValue <= node.Values[i] { // 传播过来的价值小于等于共识价值
					subNodes = append(subNodes, subNode)
				}
			}
		}
	}
	assignHashRate(subNodes, subValue, result)
}
