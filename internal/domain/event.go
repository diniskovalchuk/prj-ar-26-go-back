package domain

import "time"

type Event struct {
	Id          uint64
	DeviceId    uint64
	RoomId      uint64
	Action      string // "on" або "off"
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

// Структура для повернення аналітики енергії адміністратору
type EnergyReport struct {
	Period      string  `json:"period"`
	RoomId      uint64  `json:"room_id,omitempty"` // 0, якщо звіт по всьому підприємству
	EnergyUsed  float64 `json:"energy_used"`       // Розраховане значення кВт/год чи умовних одиниць
	EventsCount int     `json:"events_count"`
}
