package models

import (
	"sync"
	"fmt"
)

type MetrValue float64
type MetrType string

const (
	TGauge   MetrType = "gauge"
	TCounter MetrType = "counter"
)


type MetricStorage interface {
	// set or update
	// can return error
	Update(name string, mType MetrType, value MetrValue) error

	// return value and bool if exist on not
	Get(name string, mType MetrType) (MetrValue, bool)
}

type MemStorage struct {
	mu sync.RWMutex // Mutex для безопасного доступа в многопоточной среде

	gauges   map[string]MetrValue
	counters map[string]MetrValue
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]MetrValue),
		counters: make(map[string]MetrValue),
	}
}


func (s *MemStorage) Update(name string, mType MetrType, value MetrValue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch mType {
	case TGauge:
		s.gauges[name] = value
		return nil
	case TCounter:
		s.counters[name] += value
		return nil
	default:
		return fmt.Errorf("unknown metric type: %s", mType)
	}
}

// Get извлекает значение метрики.
func (s *MemStorage) Get(name string, mType MetrType) (MetrValue, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch mType {
	case TGauge:
		val, ok := s.gauges[name]
		return val, ok
	case TCounter:
		val, ok := s.counters[name]
		return val, ok
	default:
		//not found because of type
		return 0, false
	}
}