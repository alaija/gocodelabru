package storage

import (
	"github.com/alaija/gocodelabru/storage/lru"
)

type (
	// Location используется для хранения координат водителя
	Location struct {
		Latitude  float64
		Longitude float64
	}
	// Driver модель хранения водителя
	Driver struct {
		ID           int
		LastLocation Location
		Expiration   int64
		Locations    *lru.LRU
	}
)
