package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenFilter(t *testing.T) {
	fmt.Println(genFilter(5))
}

func TestTreeDisableLoop(t *testing.T) {
	us, mapper, err := TreeDisableLoop("D1", "D2", 50)
	if err != nil {
		t.Fatal(err.Error())
	}
	jstr, err := json.Marshal(us)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(jstr))
	fmt.Printf("%+v", mapper)
}
