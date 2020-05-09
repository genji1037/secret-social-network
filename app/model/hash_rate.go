package model

import "sync"

// HashRateResult represent hash rate calculate result.
type HashRateResult struct {
	AssignmentCount int
	ResultHashes    map[string]float64
	sync.Mutex
}

// Add sums sub calculate result.
func (r *HashRateResult) Add(subResultHashes map[string]float64) {
	r.Lock()
	r.AssignmentCount += len(subResultHashes)
	for k, v := range subResultHashes {
		r.ResultHashes[k] += v
	}
	r.Unlock()
}
