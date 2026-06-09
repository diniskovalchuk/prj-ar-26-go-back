package resources

import (
	"prj-ar-26-go-back/internal/domain"
	"time"
)

type EventDto struct {
	Id          uint64     `json:"id"`
	DeviceId    uint64     `json:"device_id"`
	RoomId      uint64     `json:"room_id"`
	Action      string     `json:"action"` // "on" або "off"
	CreatedDate time.Time  `json:"created_date"`
	UpdatedDate time.Time  `json:"updated_date"`
	DeletedDate *time.Time `json:"deleted_date,omitempty"` // Ховаємо з JSON, якщо nil
}

// DomainToDto конвертує одну доменну модель у DTO формат для клієнта
func (d EventDto) DomainToDto(event domain.Event) EventDto {
	return EventDto{
		Id:          event.Id,
		DeviceId:    event.DeviceId,
		RoomId:      event.RoomId,
		Action:      event.Action,
		CreatedDate: event.CreatedDate,
		UpdatedDate: event.UpdatedDate,
		DeletedDate: event.DeletedDate,
	}
}

// DomainToDtoCollection конвертує список доменних моделей у список DTO
func (d EventDto) DomainToDtoCollection(events []domain.Event) []EventDto {
	eventsDto := make([]EventDto, len(events))
	for i := range events {
		eventsDto[i] = d.DomainToDto(events[i])
	}
	return eventsDto
}
