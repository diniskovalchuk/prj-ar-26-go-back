package resources

import (
	"prj-ar-26-go-back/internal/domain"
	"time"
)

type MeasurementDto struct {
	Id          uint64     `json:"id"`
	DeviceId    uint64     `json:"device_id"`
	RoomId      uint64     `json:"room_id"`
	Value       float64    `json:"value"`
	CreatedDate time.Time  `json:"created_date"`
	UpdatedDate time.Time  `json:"updated_date"`
	DeletedDate *time.Time `json:"deleted_date,omitempty"` // omitempty приховає поле, якщо воно nil
}

// DomainToDto конвертує доменну модель у структуру для HTTP відповіді
func (d MeasurementDto) DomainToDto(meas domain.Measurement) MeasurementDto {
	return MeasurementDto{
		Id:          meas.Id,
		DeviceId:    meas.DeviceId,
		RoomId:      meas.RoomId,
		Value:       meas.Value,
		CreatedDate: meas.CreatedDate,
		UpdatedDate: meas.UpdatedDate,
		DeletedDate: meas.DeletedDate,
	}
}

// DomainToDtoCollection конвертує зріз доменних моделей у зріз DTO
func (d MeasurementDto) DomainToDtoCollection(meas []domain.Measurement) []MeasurementDto {
	measDto := make([]MeasurementDto, len(meas))
	for i := range meas {
		measDto[i] = d.DomainToDto(meas[i])
	}
	return measDto
}
