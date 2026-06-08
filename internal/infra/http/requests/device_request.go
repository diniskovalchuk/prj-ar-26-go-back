package requests

import (
	"prj-ar-26-go-back/internal/domain"

	"github.com/google/uuid"
)

type DeviceRequest struct {
	OrganizationId   uint64                `json:"organization_id,omitempty"`
	RoomId           *uint64               `json:"room_id,omitempty"`
	InventoryNumber  string                `json:"inventory_number,omitempty"`
	SerialNumber     string                `json:"serial_number,omitempty"`
	Characteristics  string                `json:"characteristics,omitempty"`
	Category         domain.DeviceCategory `json:"category,omitempty"`
	Units            string                `json:"units,omitempty"`
	PowerConsumption float64               `json:"power_consumption,omitempty"`
}

func (r DeviceRequest) ToDomainModel() (interface{}, error) {
	return domain.Device{
		OrganizationId:   r.OrganizationId,
		InventoryNumber:  r.InventoryNumber,
		SerialNumber:     r.SerialNumber,
		Characteristics:  r.Characteristics,
		Category:         r.Category,
		Units:            r.Units,
		PowerConsumption: r.PowerConsumption,
		RoomId:           r.RoomId,
		GUID:             uuid.New().String(),
	}, nil
}
