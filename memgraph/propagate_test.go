package memgraph

import (
	"fmt"
	"secret-social-network/service"
	"strings"
	"testing"
	"time"
)

func TestPropagate(t *testing.T) {
	db := NewMemDB()
	links := []struct {
		uid1  string
		uid2  string
		value float64
	}{
		{uid1: "A1", uid2: "A2", value: 100},
		{uid1: "A2", uid2: "A3", value: 100},
		{uid1: "A3", uid2: "A4", value: 100},
		{uid1: "A4", uid2: "A5", value: 100},
		{uid1: "A5", uid2: "A1", value: 100},
	}
	for _, link := range links {
		db.Link(link.uid1, link.uid2, link.value)
	} // propagate
	result := db.PropagateAll()
	fmt.Printf("%v\n", result)
}

func BenchmarkPropagate(b *testing.B) {
	db := NewMemDB()
	links := []struct {
		uid1  string
		uid2  string
		value float64
	}{
		{uid1: "A1", uid2: "A2", value: 100},
		{uid1: "A2", uid2: "A3", value: 100},
		{uid1: "A3", uid2: "A4", value: 100},
		{uid1: "A4", uid2: "A5", value: 100},
		{uid1: "A5", uid2: "A1", value: 100},
	}
	for _, link := range links {
		db.Link(link.uid1, link.uid2, link.value)
	}

	for i := 0; i < b.N; i++ {
		db.PropagateAll()
	}
}

func TestBenchPropagate(t *testing.T) {
	db := NewMemDB()
	// list relation network
	relations, err := service.ListRelationWithValue()
	if err != nil {
		t.Fatalf("list relation failed: %s", err.Error())
	}

	// prepare mem data
	startAt := time.Now()
	for relation, value := range relations {
		tmp := strings.Split(relation, "-")
		if tmp[0] < tmp[1] { // 双向关系去重
			db.Link(tmp[0], tmp[1], value)
		}
	}
	fmt.Printf("[link] %s\n", time.Now().Sub(startAt))

	// propagate
	db.PropagateAll()
}
