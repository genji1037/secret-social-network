package service

import (
	"fmt"
	"testing"
)

func TestPropagate(t *testing.T) {
	rs, err := propagate()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("%+v", rs)
}

func TestListRelationWithValue(t *testing.T) {
	rs, err := ListRelationWithValue()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("%d", len(rs))
}
