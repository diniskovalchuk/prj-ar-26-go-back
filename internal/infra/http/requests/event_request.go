package requests

import "prj-ar-26-go-back/internal/domain"

type EventRequest struct {
	DeviceId uint64 `json:"device_id"`
	RoomId   uint64 `json:"room_id"`
	Action   string `json:"action"` // "on" або "off"
}

func (r EventRequest) ToDomainModel() (interface{}, error) {
	return domain.Event{
		DeviceId: r.DeviceId,
		RoomId:   r.RoomId,
		Action:   r.Action,
	}, nil
}
