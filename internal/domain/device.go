package domain

import "time"

type Device struct {
	Id               uint64
	OrganizationId   uint64
	RoomId           *uint64
	GUID             string
	InventoryNumber  string
	SerialNumber     string
	Characteristics  string
	Category         DeviceCategory
	Units            string
	PowerConsumption float64
	CreatedDate      time.Time
	UpdatedDate      time.Time
	DeletedDate      *time.Time
}

type DeviceCategory string

const (
	Sensor   DeviceCategory = "SENSOR"
	Actuator DeviceCategory = "ACTUATOR"
)
