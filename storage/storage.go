package storage

import "errors"

// DriverStorage структура для хранения водителей
type DriverStorage struct {
	drivers map[int]*Driver
}

// New создает новое хранилище `DriverStorage`
func New() *DriverStorage {
	d := &DriverStorage{}
	d.drivers = make(map[int]*Driver)
	return d
}

// Set добавляет водителя `driver` с ключом `key`
func (d *DriverStorage) Set(key int, driver *Driver) {
	d.drivers[key] = driver
}

// Delete удаляет водителя `driver` с ключом `key`
func (d *DriverStorage) Delete(key int) error {
	_, ok := d.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	delete(d.drivers, key)
	return nil
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
func (d *DriverStorage) Nearest(radius, lat, lon float64) []*Driver {
	var nearest []*Driver
	for _, driver := range d.drivers {
		if Distance(lat, lon, driver.LastLocation.Latitude, driver.LastLocation.Longitude) <= radius {
			nearest = append(nearest, driver)
		}
	}
	return nearest
}
