package service

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		link(600, 3000, 100, "BA")
		wg.Done()
	}()
	go func() {
		link(600, 3000, 100, "BB")
		wg.Done()
	}()
	go func() {
		link(600, 3000, 100, "BC")
		wg.Done()
	}()
	go func() {
		link(600, 3000, 100, "BD")
		wg.Done()
	}()
	wg.Wait()
}

func link(uNum, rNum int, value float64, uPrefix string) {
	rand.Seed(time.Now().UnixNano())

	uList := make([]string, 0, uNum)
	linked := make(map[string]struct{})
	for i := 0; i < uNum; i++ {
		uList = append(uList, fmt.Sprintf("%s%d", uPrefix, i))
	}

	for i := 0; i < rNum; i++ {
		ui1 := rand.Intn(uNum)
		ui2 := rand.Intn(uNum)
		if ui1 == ui2 {
			i--
			continue
		}
		var UID1, UID2 string
		if ui1 < ui2 {
			UID1 = uList[ui1]
			UID2 = uList[ui2]
		} else {
			UID1 = uList[ui2]
			UID2 = uList[ui1]
		}
		key := UID1 + "-" + UID2
		_, ok := linked[key]
		if ok {
			// already linked
			i--
			continue
		}
		err := Link(UID1, UID2, value)
		if err != nil {
			fmt.Printf("[ERROR] link failed: %s\n", err.Error())
			i--
			continue
		}
		linked[key] = struct{}{}
	}
}
