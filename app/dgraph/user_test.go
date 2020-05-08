package dgraph

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	err := Open()
	if err != nil {
		panic(err.Error())
	}
	m.Run()
}

func TestUser_Link(t *testing.T) {
	err := User{}.LinkOrAdd("app2", "QQ8", "QQ9", 100)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestQueryValue(t *testing.T) {
	appValue, totalValue, err := QueryValue("app3", "QQ8", "QQ9")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("app: %f\ntotal: %f\n", appValue, totalValue)
}
