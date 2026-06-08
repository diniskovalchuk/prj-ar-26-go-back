package resources

import "prj-ar-26-go-back/internal/domain"

type DeviceDto struct {
	Id               uint64                `json:"id"`
	OrganizationId   uint64                `json:"organizationId"`
	RoomId           *uint64               `json:"roomId"`
	GUID             string                `json:"GUID"`
	InventoryNumber  string                `json:"inventoryNumber"`
	SerialNumber     string                `json:"serialNumber"`
	Characteristics  string                `json:"characteristics"`
	Category         domain.DeviceCategory `json:"category"`
	Units            string                `json:"units"`
	PowerConsumption float64               `json:"powerConsumption"`
}

func (d DeviceDto) DomainToDto(dev domain.Device) DeviceDto {
	return DeviceDto{
		Id:               dev.Id,
		OrganizationId:   dev.OrganizationId,
		RoomId:           dev.RoomId,
		GUID:             dev.GUID,
		InventoryNumber:  dev.InventoryNumber,
		SerialNumber:     dev.SerialNumber,
		Characteristics:  dev.Characteristics,
		Category:         dev.Category,
		Units:            dev.Units,
		PowerConsumption: dev.PowerConsumption,
	}
}

func (d DeviceDto) DomainToDtoCollection(devs []domain.Device) []DeviceDto {
	devsDto := make([]DeviceDto, len(devs))
	for i := range devs {
		devsDto[i] = d.DomainToDto(devs[i])
	}
	return devsDto
}
