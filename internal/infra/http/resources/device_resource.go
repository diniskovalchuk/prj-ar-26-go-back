package resources

import "github.com/diniskovalchuk/prj2/internal/domain"

type DeviceDto struct {
	
}

func (d DeviceDto) DomainToDto(dev domain.Device) DeviceDto {
	return DeviceDto{
		Id               dev.Id
		OrganizationId   dev.OrganizationId
		RoomId           dev.RoomId
		GUID             dev.GUID
		InventoryNumber  dev.InventoryNumber
		SerialNumber     dev.SerialNumber
		Characteristics  dev.Characteristics
		Category         dev.Category
		Units            dev.Units
		PowerConsumption dev.PowerConsumption
		CreatedDate      dev.CreatedDate
		UpdatedDate      dev.UpdatedDate
		DeletedDate      dev.DeletedDate
	}
}

func (d DeviceDto) DomainToDtoCollection(devs []domain.Device) []DeviceDto {
	devsDto := make([]DeviceDto, len(devs))
	for i := range devs {
		devsDto[i] = d.DomainToDto(devs[i])
	}
	return devsDto
}
