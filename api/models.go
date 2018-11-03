package api

type (
	// Location для представления координат водителя в запросе
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}
	// Payload запроса с данными водителя
	Payload struct {
		Timestamp int64    `json:"timestamp"`
		DriverID  int      `json:"driver_id"`
		Location  Location `json:"location"`
	}
	// DefaultResponse для возврата ответа по умолчанию
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	// DriverResponse для возврата ответа, когда мы запрашиваем водителя
	DriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Driver  int    `json:"driver"`
	}
	// NearestDriverResponse для возврата ближайших водителей
	NearestDriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Drivers []int  `json:"drivers"`
	}
)
