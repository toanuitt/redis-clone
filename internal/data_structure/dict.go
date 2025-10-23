package data_structure

import (
	"log"
	"math/rand"
	"redis-clone/internal/config"
	"time"
)

type Obj struct {
	Value          interface{}
	LastAccessTime uint32
	AccessCount    uint8
	LastDecayTime  uint32
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func CreateDict() *Dict {
	res := Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
	return &res
}

func (d *Dict) GetexpiredDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) GetDictStore() map[string]*Obj {
	return d.dictStore
}

func now() uint32 {
	return uint32(time.Now().Unix())
}

func (d *Dict) NewObject(key string, value interface{}, ttlMs int64) *Obj {
	obj := &Obj{
		Value:          value,
		LastAccessTime: now(),
		AccessCount:    1, // Initial counter value
		LastDecayTime:  now(),
	}
	if ttlMs > 0 {
		d.SetExpired(key, ttlMs)
	}
	return obj
}

func (d *Dict) GetExpired(key string) (uint64, bool) {
	exp, exists := d.expiredDictStore[key]
	return exp, exists
}

func (d *Dict) SetExpired(key string, ttlMs int64) {
	d.expiredDictStore[key] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func (d *Dict) HasExpired(key string) bool {
	exp, exist := d.expiredDictStore[key]
	if !exist {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func (obj *Obj) updateLFUCounter() {
	currentTime := now()
	timeDiff := currentTime - obj.LastDecayTime

	if timeDiff > 0 {
		numDecays := timeDiff / 60
		if numDecays > 0 {
			if obj.AccessCount > uint8(numDecays) {
				obj.AccessCount -= uint8(numDecays)
			} else {
				obj.AccessCount = 0
			}
			obj.LastDecayTime = currentTime
		}
	}
	if obj.AccessCount < 255 {
		baseProbability := 1.0 / float64(obj.AccessCount+1)
		if rand.Float64() < baseProbability {
			obj.AccessCount++
		}
	}
}

func (d *Dict) Get(key string) *Obj {
	v := d.dictStore[key]
	if v != nil {
		v.LastAccessTime = now()
		v.updateLFUCounter()
		if d.HasExpired(key) {
			d.Delete(key)
			return nil
		}
	}
	return v
}

func (d *Dict) evictRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("trigger random eviction")
	for k := range d.dictStore {
		d.Delete(k)
		evictCount--
		if evictCount == 0 {
			break
		}
	}
}

func (d *Dict) populatePool() {
	remain := config.EpoolLruSampleSize
	for k := range d.dictStore {
		ePool.Push(k, d.dictStore[k].LastAccessTime)
		remain--
		if remain == 0 {
			break
		}
	}
	log.Printf("Epool")
	for _, item := range ePool.pool {
		log.Println(item.key, item.lastAccessTime)
	}
}

func (d *Dict) evictLru() {
	d.populatePool()
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("trigger LRU eviction")
	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item != nil {
			d.Delete(item.key)
		}
	}
}

func (d *Dict) populatePoolLfu() {
	lfuEvictionPool.Clear()
	remain := config.EpoolLfuSampleSize

	for k, obj := range d.dictStore {
		idleTime := now() - obj.LastAccessTime
		freq := uint64(obj.AccessCount)
		if idleTime > 0 {
			idleFactor := idleTime / 60
			if idleFactor > 0 && freq > 0 {
				freq = freq / (1 + uint64(idleFactor))
			}
		}

		lfuEvictionPool.Push(k, freq)
		remain--
		if remain == 0 {
			break
		}
	}
}

func (d *Dict) evictLfu() {
	d.populatePoolLfu()
	evictCount := int64(config.EvictionRatio * float64(config.MaxKeyNumber))
	log.Print("trigger LFU eviction")

	for i := 0; i < int(evictCount) && len(lfuEvictionPool.pool) > 0; i++ {
		item := lfuEvictionPool.Pop()
		if item != nil {
			d.Delete(item.key)
		}
	}
}

func (d *Dict) evict() {
	switch config.EvictionPolicy {
	case "allkeys-random":
		d.evictRandom()
	case "allkeys-lru":
		d.evictLru()
	case "allkeys-lfu":
		d.evictLfu()
	}
}

func (d *Dict) Set(key string, obj *Obj) {
	if len(d.dictStore) == config.MaxKeyNumber {
		d.evict()
	}
	v := d.dictStore[key]
	if v == nil {
		HashKeySpaceStat.Key++
	}
	d.dictStore[key] = obj
}

func (d *Dict) Delete(key string) bool {
	log.Printf("Delete key %s", key)
	if _, exists := d.dictStore[key]; exists {
		delete(d.dictStore, key)
		delete(d.expiredDictStore, key)
		HashKeySpaceStat.Key--
		return true
	}
	return false

}
