package client

import (
	"fmt"
	"secret-social-network/app/config"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.LoadConfig("../../config/server.yml")
	if err != nil {
		panic(err.Error())
	}
	m.Run()
}

func TestGetUID(t *testing.T) {
	openID1, openID2 := "4b8e530c-9020-5a1f-69a0-794c9884865a", "98a9ad48-7438-e337-ee05-659b5f754924"
	uid1, uid2, err := GetUID("05d22151483ff000", openID1, openID2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(uid1, uid2)
}
