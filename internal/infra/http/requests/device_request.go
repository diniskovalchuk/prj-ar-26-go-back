package requests


type DeviceRequest struct {
	Id               uint64 `db:"id,omitempty"`
	OrganizationId   uint64 `db:"organization_id,omitempty"`
	RoomId           *uint64 `db:"room_id,omitempty"`
	GUID             string `db:"guid,omitempty"`
	InventoryNumber  string `db:"inventory_number,omitempty"`
	SerialNumber     string `db:"serial_number,omitempty"`
	Characteristics  string `db:"characteristics,omitempty"`
	Category         DeviceCategory `db:"category,omitempty"`
	Units            string `db:"units,omitempty"`
	PowerConsumption float64 `db:"power_consumption,omitempty"`
	CreatedDate      time.Time `db:"created_date,omitempty"`
	UpdatedDate      time. Time `db:"updated_date,omitempty"`
	DeletedDate      *time.Time `db:"deleted_date,omitempty"`
}

func (r DeviceRequest) ToDomainModel() (interface{}, error) {
	return domain.Device{
		Id               r.Id
		OrganizationId   r.OrganizationId
		RoomId           r.RoomId
		GUID             r.GUID
		InventoryNumber  r.InventoryNumber
		SerialNumber     r.SerialNumber
		Characteristics  r.Characteristics
		Category         r.Category
		Units            r.Units
		PowerConsumption r.PowerConsumption
		CreatedDate      r.CreatedDate
		UpdatedDate      r.UpdatedDate
		DeletedDate      r.DeletedDate
	}, nil
}
