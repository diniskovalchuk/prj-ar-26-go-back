package requests

import (
	"errors"
	"prj-ar-26-go-back/internal/domain"
)

type MeasurementRequest struct {
	DeviceId uint64  `json:"device_id"`
	RoomId   uint64  `json:"room_id"`
	Value    float64 `json:"value"`
}

// ToDomainModel тепер повертає конкретний тип domain.Measurement замість interface{}
func (r MeasurementRequest) ToDomainModel() (interface{}, error) {
	return domain.Measurement{
		DeviceId: r.DeviceId,
		RoomId:   r.RoomId,
		Value:    r.Value,
	}, nil
}

// Додатково: метод для базової валідації вхідних даних від пристроїв
func (r MeasurementRequest) Validate() error {
	if r.DeviceId == 0 {
		return errors.New("device_id is required")
	}
	if r.RoomId == 0 {
		return errors.New("room_id is required")
	}
	// Значення Value може бути 0 (наприклад, 0 градусів), тому перевіряємо лише критичні ID
	return nil
}
