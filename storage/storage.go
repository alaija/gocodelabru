package storage

import (
	"errors"

	"github.com/dhconnelly/rtreego"
)

// DriverStorage структура для хранения водителей
type DriverStorage struct {
	drivers   map[int]*Driver
	locations *rtreego.Rtree
}

// Bounds метод необходимы для корректной работы rtree
// Latitude - Y, Longitude - X в координатной системе
func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Latitude, d.LastLocation.Longitude}.ToRect(0.01)
}

// New создает новое хранилище `DriverStorage`
func New() *DriverStorage {
	d := &DriverStorage{}
	d.drivers = make(map[int]*Driver)
	d.locations = rtreego.NewTree(2, 25, 50)
	return d
}

// Set добавляет водителя `driver` с ключом `key`
func (d *DriverStorage) Set(key int, driver *Driver) {
	_, ok := d.drivers[key]
	if !ok {
		d.locations.Insert(driver)
	}
	d.drivers[key] = driver
}

// Delete удаляет водителя `driver` с ключом `key`
func (d *DriverStorage) Delete(key int) error {
	driver, ok := d.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	if d.locations.Delete(driver) {
		delete(d.drivers, key)
		return nil
	}
	return errors.New("could not remove driver")
}

// Get возвращет водителя по ключу `key` или ошибку, если ничего не найдено
func (d *DriverStorage) Get(key int) (*Driver, error) {
	driver, ok := d.drivers[key]
	if !ok {
		return nil, errors.New("Driver does not exist")
	}
	return driver, nil
}

// Nearest возвращет ближйщих водителей по координатам `lat` и `lon` или ошибку, если ничего не найдено
func (d *DriverStorage) Nearest(count int, lat, lon float64) []*Driver {
	point := rtreego.Point{lat, lon}
	results := d.locations.NearestNeighbors(count, point)
	var nearest []*Driver
	for _, item := range results {
		if item == nil {
			continue
		}
		nearest = append(nearest, item.(*Driver))
	}
	return nearest
}
