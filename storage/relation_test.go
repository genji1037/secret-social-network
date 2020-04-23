package storage

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

func TestListRelation(t *testing.T) {
	startAt := time.Now()
	mapper1, err := ListRelation(5000)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("%s\n", time.Now().Sub(startAt))
	startAt = time.Now()
	mapper2, err := ListRelation(10000)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("%s\n", time.Now().Sub(startAt))
	s1 := fmt.Sprintf("%v", mapper1)
	s2 := fmt.Sprintf("%v", mapper2)
	fmt.Println()
	fmt.Println(fmt.Sprintf("%x", md5.Sum([]byte(s1))))
	fmt.Println(fmt.Sprintf("%x", md5.Sum([]byte(s2))))

}
