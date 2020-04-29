package model

import "sync"

type Result struct {
	AssignmentCount int
	ResultHashes    map[string]float64
	sync.Mutex
}

func (r *Result) Add(subResultHashes map[string]float64) {
	r.Lock()
	r.AssignmentCount += len(subResultHashes)
	for k, v := range subResultHashes {
		r.ResultHashes[k] += v
	}
	r.Unlock()
}
