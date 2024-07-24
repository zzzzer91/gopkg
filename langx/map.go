package langx

import (
	"strings"
	"sync"
)

func ConvertSliceToMap[T any, K comparable](slice []T, keySelector func(T) K) map[K]T {
	result := make(map[K]T)
	for _, v := range slice {
		key := keySelector(v)
		result[key] = v
	}
	return result
}

func MergeTwoMap[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V, len(m1)+len(m2))
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result
}

func MergeMapBToA[K comparable, V any](mA, mB map[K]V) map[K]V {
	if mA == nil {
		mA = make(map[K]V, len(mB))
	}
	for k, v := range mB {
		mA[k] = v
	}
	return mA
}

type LazyMap[K comparable, V any] struct {
	m    map[K]*V
	lock sync.RWMutex

	constructor func(k K) *V
}

func NewLazyMap[K comparable, V any](constructor func(k K) *V) *LazyMap[K, V] {
	return &LazyMap[K, V]{
		m:           make(map[K]*V),
		constructor: constructor,
	}
}

func (m *LazyMap[K, V]) Get(k K) *V {
	m.lock.RLock()
	if v, ok := m.m[k]; ok {
		m.lock.RUnlock()
		return v
	}
	m.lock.RUnlock()

	m.lock.Lock()
	defer m.lock.Unlock()
	if v, ok := m.m[k]; ok {
		return v
	}
	v := m.constructor(k)
	if v != nil {
		m.m[k] = v
	}
	return v
}

// CaseInsensitiveMap is a case-insensitive map.
type CaseInsensitiveMap[V any] struct {
	data map[string]V
}

func NewCaseInsensitiveMap[V any]() *CaseInsensitiveMap[V] {
	return &CaseInsensitiveMap[V]{
		data: make(map[string]V),
	}
}

func (m *CaseInsensitiveMap[V]) Set(key string, value V) {
	lowerKey := strings.ToLower(key)
	m.data[lowerKey] = value
}

func (m *CaseInsensitiveMap[V]) Get(key string) (V, bool) {
	if m == nil {
		var v V
		return v, false
	}
	lowerKey := strings.ToLower(key)
	value, exists := m.data[lowerKey]
	return value, exists
}

func (m *CaseInsensitiveMap[V]) Delete(key string) {
	if m == nil {
		return
	}
	lowerKey := strings.ToLower(key)
	delete(m.data, lowerKey)
}

func (m *CaseInsensitiveMap[V]) Exists(key string) bool {
	if m == nil {
		return false
	}
	lowerKey := strings.ToLower(key)
	_, exists := m.data[lowerKey]
	return exists
}

func (m *CaseInsensitiveMap[V]) ToMap() map[string]V {
	if m == nil {
		return nil
	}
	newM := make(map[string]V, len(m.data))
	for k, v := range m.data {
		newM[k] = v
	}
	return newM
}

// CaseInsensitiveSet is a case-insensitive set.
type CaseInsensitiveSet struct {
	data *CaseInsensitiveMap[struct{}]
}

func NewCaseInsensitiveSet() *CaseInsensitiveSet {
	return &CaseInsensitiveSet{
		data: NewCaseInsensitiveMap[struct{}](),
	}
}

func (s *CaseInsensitiveSet) Add(value string) {
	s.data.Set(value, struct{}{})
}

func (s *CaseInsensitiveSet) Delete(value string) {
	s.data.Delete(value)
}

func (s *CaseInsensitiveSet) Exists(value string) bool {
	return s.data.Exists(value)
}
