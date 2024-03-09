package main

import "sync"

type SafeMap struct {
	mu  sync.Mutex
	val map[string]float64
}

func (m *SafeMap) Set(key string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.val[key] = value
}

func (m *SafeMap) Value(key string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.val[key]
}
