package dgraph

import (
	"testing"
)

//func TestMain(m *testing.M) {
//	err := Open()
//	if err != nil {
//		panic(err.Error())
//	}
//	m.Run()
//}

func TestUser_Link(t *testing.T) {
	err := LinkOrAdd("app2", "QQ8", "QQ9", 100)
	if err != nil {
		t.Fatal(err.Error())
	}
}
