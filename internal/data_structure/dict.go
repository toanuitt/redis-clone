package data_structure

import (
	"log"
	"redis-clone/internal/config"
	"time"
)

type Obj struct {
	Value          interface{}
	LastAccessTime uint32
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

func (d *Dict) Get(key string) *Obj {
	v := d.dictStore[key]
	if v != nil {
		v.LastAccessTime = now()
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

func (d *Dict) evict() {
	switch config.EvictionPolicy {
	case "allkeys-random":
		d.evictRandom()
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
