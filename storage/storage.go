package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/alaija/gocodelabru/storage/lru"
	"github.com/dhconnelly/rtreego"
)

// DriverStorage структура для хранения водителей
type DriverStorage struct {
	mu        *sync.RWMutex
	drivers   map[int]*Driver
	locations *rtreego.Rtree
	lruSize   int
}

// Expired возвращает `true` если значение просрочено.
func (d *Driver) Expired() bool {
	if d.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > d.Expiration
}

// Bounds метод необходимы для корректной работы rtree
// Latitude - Y, Longitude - X в координатной системе
func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Latitude, d.LastLocation.Longitude}.ToRect(0.01)
}

// New создает новое хранилище `DriverStorage`
func New(lruSize int) *DriverStorage {
	s := &DriverStorage{}
	s.drivers = make(map[int]*Driver)
	s.locations = rtreego.NewTree(2, 25, 50)
	s.mu = new(sync.RWMutex)
	s.lruSize = lruSize
	return s
}

// Set добавляет водителя `driver` с ключом `key`
func (s *DriverStorage) Set(driver *Driver) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.drivers[driver.ID]
	if !ok {
		d = driver
		cache, err := lru.New(s.lruSize)
		if err != nil {
			return errors.New("could not create LRU")
		}
		d.Locations = cache
		s.locations.Insert(d)
	}
	d.LastLocation = driver.LastLocation
	d.Locations.Add(time.Now().UnixNano(), d.LastLocation)
	d.Expiration = driver.Expiration
	s.drivers[driver.ID] = driver
	return nil
}

// Delete удаляет водителя `driver` с ключом `key`
func (s *DriverStorage) Delete(key int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.drivers[key]
	if !ok {
		return errors.New("does not exist")
	}
	deleted := s.locations.Delete(d)
	if deleted {
		delete(s.drivers, d.ID)
		return nil
	}
	return errors.New("could not remove item")
}

// Get возвращет водителя по ключу `key` или ошибку, если ничего не найдено
func (s *DriverStorage) Get(key int) (*Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.drivers[key]
	if !ok {
		return nil, errors.New("does not exist")
	}
	return d, nil
}

// Nearest возвращет ближйщих водителей по координатам `lat` и `lon` или ошибку, если ничего не найдено
func (s *DriverStorage) Nearest(point rtreego.Point, count int) []*Driver {
	s.mu.Lock()
	defer s.mu.Unlock()
	results := s.locations.NearestNeighbors(count, point)
	var drivers []*Driver
	for _, item := range results {
		if item == nil {
			continue
		}
		drivers = append(drivers, item.(*Driver))
	}
	return drivers
}

// DeleteExpired удаляет все просроченные значения
func (s *DriverStorage) DeleteExpired() {
	now := time.Now().UnixNano()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.drivers {
		if v.Expiration > 0 && now > v.Expiration {
			deleted := s.locations.Delete(v)
			if deleted {
				delete(s.drivers, v.ID)
			}
		}
	}
}
