package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестирует основные методы
func TestStorage(t *testing.T) {
	s := New()
	driver := &Driver{
		ID: 1,
		LastLocation: Location{
			Latitude:  1,
			Longitude: 1,
		},
	}

	s.Set(driver.ID, driver)
	d, err := s.Get(driver.ID)
	assert.NoError(t, err)
	assert.Equal(t, d, driver)

	err = s.Delete(driver.ID)
	assert.NoError(t, err)

	d, err = s.Get(driver.ID)
	assert.Equal(t, err, errors.New("Driver does not exist"))

	err = s.Delete(driver.ID)
	assert.Equal(t, err, errors.New("Driver does not exist"))
}

// Тестирует метод `Nearest`
func TestNearest(t *testing.T) {

	s := New()
	s.Set(123, &Driver{
		ID: 123,
		LastLocation: Location{
			Latitude:  1,
			Longitude: 1,
		},
	})
	s.Set(666, &Driver{
		ID: 666,
		LastLocation: Location{
			Latitude:  42.875799,
			Longitude: 74.588279,
		},
	})

	drivers := s.Nearest(1, 42.876420, 74.588332)
	assert.Equal(t, len(drivers), 1)
}

// Бенчмарк метода `Nearest`
func BenchmarkNearest(b *testing.B) {
	s := New()
	for i := 0; i < 100; i++ {
		s.Set(i, &Driver{
			ID: i,
			LastLocation: Location{
				Latitude:  float64(i),
				Longitude: float64(i),
			},
		})
	}
	for i := 0; i < b.N; i++ {
		s.Nearest(10, 123, 123)
	}
}
