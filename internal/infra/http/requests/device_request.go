package requests

import (
	"github.com/diniskovalchuk/prj2/internal/domain"
)

type DeviceRequest struct {
	RoomId           *uint64               `db:"room_id,omitempty"`
	InventoryNumber  string                `db:"inventory_number,omitempty"`
	SerialNumber     string                `db:"serial_number,omitempty"`
	Characteristics  string                `db:"characteristics,omitempty"`
	Category         domain.DeviceCategory `db:"category,omitempty"`
	Units            string                `db:"units,omitempty"`
	PowerConsumption float64               `db:"power_consumption,omitempty"`
}

func (r DeviceRequest) ToDomainModel() (interface{}, error) {
	return domain.Device{
		RoomId:           r.RoomId,
		InventoryNumber:  r.InventoryNumber,
		SerialNumber:     r.SerialNumber,
		Characteristics:  r.Characteristics,
		Category:         r.Category,
		Units:            r.Units,
		PowerConsumption: r.PowerConsumption,
	}, nil
}
