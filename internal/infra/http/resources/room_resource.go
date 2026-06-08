package resources

import "prj-ar-26-go-back/internal/domain"

type RoomDto struct {
	Id             uint64  `json:"id"`
	OrganizationId uint64  `json:"organizationId"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
}

func (d RoomDto) DomainToDto(r domain.Room) RoomDto {
	return RoomDto{
		Id:             r.Id,
		OrganizationId: r.OrganizationId,
		Name:           r.Name,
		Description:    r.Description,
	}
}

func (d RoomDto) DomainToDtoCollection(rooms []domain.Room) []RoomDto {
	roomsDto := make([]RoomDto, len(rooms))
	for i := range rooms {
		roomsDto[i] = d.DomainToDto(rooms[i])
	}
	return roomsDto
}
