package dgraph

import "testing"

func TestUser_Link(t *testing.T) {
	err := User{}.Link("TT1", "TT4", 200)
	if err != nil {
		t.Fatal(err.Error())
	}
}
