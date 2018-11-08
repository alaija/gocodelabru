package storage

import (
	"testing"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/stretchr/testify/assert"
)

// Тестирует основные методы
func TestStorage(t *testing.T) {
	s := New(10)
	driver := &Driver{
		ID: 1,
		LastLocation: Location{
			Latitude:  1,
			Longitude: 1,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	}

	s.Set(driver)
	d, err := s.Get(driver.ID)
	assert.NoError(t, err)
	assert.Equal(t, d, driver)

	err = s.Delete(driver.ID)
	assert.NoError(t, err)

	d, err = s.Get(driver.ID)
	assert.Error(t, err)

	err = s.Delete(driver.ID)
	assert.Error(t, err)
}

// Тестирует метод `Nearest`
func TestNearest(t *testing.T) {
	s := New(10)
	s.Set(&Driver{
		ID: 123,
		LastLocation: Location{
			Latitude:  42.875799,
			Longitude: 74.588279,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 321,
		LastLocation: Location{
			Latitude:  42.875508,
			Longitude: 74.588107,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 666,
		LastLocation: Location{
			Latitude:  42.876106,
			Longitude: 74.588204,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 2319,
		LastLocation: Location{
			Latitude:  42.874942,
			Longitude: 74.585908,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 991,
		LastLocation: Location{
			Latitude:  42.875744,
			Longitude: 74.584503,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})

	drivers := s.Nearest(rtreego.Point{42.876420, 74.588332}, 3)
	assert.Equal(t, len(drivers), 3)
	assert.Equal(t, drivers[0].ID, 123)
	assert.Equal(t, drivers[1].ID, 321)
	assert.Equal(t, drivers[2].ID, 666)
}

func TestExpire(t *testing.T) {
	s := New(10)
	driver := &Driver{
		ID: 123,
		LastLocation: Location{
			Latitude:  42.876420,
			Longitude: 74.588332,
		},
		Expiration: time.Now().Add(-15).UnixNano(),
	}
	driver2 := &Driver{
		ID: 13,
		LastLocation: Location{
			Latitude:  42.876420,
			Longitude: 74.588332,
		},
		Expiration: time.Now().Add(time.Second).UnixNano(),
	}

	driver3 := &Driver{
		ID: 1,
		LastLocation: Location{
			Latitude:  42.876420,
			Longitude: 74.588332,
		},
		Expiration: 0,
	}
	s.Set(driver)
	s.Set(driver2)
	s.Set(driver3)
	assert.True(t, driver.Expired())
	assert.False(t, driver2.Expired())
	assert.False(t, driver3.Expired())
	s.DeleteExpired()
	d, err := s.Get(123)
	assert.Error(t, err)
	assert.NotEqual(t, d, driver)
}
