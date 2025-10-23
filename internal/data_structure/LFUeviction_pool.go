package data_structure

import (
	"redis-clone/internal/config"
	"sort"
)

type LFUEvictionCandidate struct {
	key  string
	freq uint64
}

type LFUEvictionPool struct {
	pool []*LFUEvictionCandidate
}

type ByFreq []*LFUEvictionCandidate

func (a ByFreq) Len() int {
	return len(a)
}

func (a ByFreq) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByFreq) Less(i, j int) bool {
	return a[i].freq < a[j].freq
}

func (p *LFUEvictionPool) Push(key string, freq uint64) {
	candidate := &LFUEvictionCandidate{
		key:  key,
		freq: freq,
	}
	p.pool = append(p.pool, candidate)
	p.Sort()
	if len(p.pool) > config.EpoolLfuSampleSize {
		// remove the largest freq (keep smallest freq on left)
		p.pool = p.pool[:config.EpoolLfuSampleSize]
	}
}

func (p *LFUEvictionPool) Sort() {
	sort.Sort(ByFreq(p.pool))
}

func (p *LFUEvictionPool) Pop() *LFUEvictionCandidate {
	if len(p.pool) == 0 {
		return nil
	}
	candidate := p.pool[0]
	p.pool = p.pool[1:]
	return candidate
}

func (p *LFUEvictionPool) Clear() {
	p.pool = p.pool[:0]
}

func NewLFUEvictionPool(size int) *LFUEvictionPool {
	return &LFUEvictionPool{
		pool: make([]*LFUEvictionCandidate, 0, size),
	}
}

var lfuEvictionPool = NewLFUEvictionPool(0)
