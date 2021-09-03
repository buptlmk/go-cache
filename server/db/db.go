package db

import (
	"go-cache/internal/hash"
	"math"
	"sync"
	"sync/atomic"
)

// 采用分段锁的概念，降低锁的竞争
type Dict struct {
	m   []*shared
	num int32
}

type DB struct {
	Data *Dict
}

type shared struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

func NewDict(size int) *Dict {
	size = compute(size)
	d := &Dict{}
	for i := 0; i < size; i++ {
		d.m = append(d.m, &shared{
			data: make(map[string]interface{}, 16),
		})
	}
	return d
}

func NewDB(size ...int) *DB {
	var d *Dict
	if len(size) == 0 {
		d = NewDict(16)
	} else {
		d = NewDict(size[len(size)-1])
	}
	return &DB{
		Data: d,
	}
}

func (d *Dict) spread(hashCode uint32) uint32 {
	if d == nil {
		panic("d is nil")
	}

	num := uint32(len(d.m))

	// 限制num 必须为 2的次方数
	return (num - 1) & hashCode // == hashcode%num
}

func (d *Dict) getShared(index uint32) *shared {

	if d == nil {
		panic("d is nil")
	}

	return d.m[index]
}

func (d *Dict) addCount() int32 {
	return atomic.AddInt32(&d.num, 1)
}

func (d *Dict) decreaseCount() int32 {
	return atomic.AddInt32(&d.num, -1)
}

func (d *Dict) Get(key string) (val interface{}, exist bool) {
	if d == nil {
		panic("d is nil")
	}
	hashCode := hash.FNV(key)
	index := d.spread(hashCode)
	s := d.getShared(index)

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, exist = s.data[key]
	return
}

func (d *Dict) Set(key string, val interface{}) (result int) {
	if d == nil {
		panic("d is nil")
	}
	hashCode := hash.FNV(key)
	index := d.spread(hashCode)
	s := d.getShared(index)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.data[key]; ok {
		s.data[key] = val
		return 0
	}
	s.data[key] = val
	d.addCount()

	return 1
}

func (d *Dict) Del(key string) (result int) {
	if d == nil {
		panic("d is nil")
	}
	hashCode := hash.FNV(key)
	index := d.spread(hashCode)
	s := d.getShared(index)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return 1
	}
	d.decreaseCount()

	return 0
}

func (d *Dict) Exec(key string, fn func(s *shared) interface{}) (result interface{}) {
	if d == nil {
		panic("d is nil")
	}
	hashCode := hash.FNV(key)
	index := d.spread(hashCode)
	s := d.getShared(index)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	return fn(s)
}

func compute(n int) int {

	if n <= 16 {
		return 16
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}
