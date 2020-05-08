package model

import "sync"

type HashRateResult struct {
	AssignmentCount int
	ResultHashes    map[string]float64
	sync.Mutex
}

func (r *HashRateResult) Add(subResultHashes map[string]float64) {
	r.Lock()
	r.AssignmentCount += len(subResultHashes)
	for k, v := range subResultHashes {
		r.ResultHashes[k] += v
	}
	r.Unlock()
}
