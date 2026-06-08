package requests

import "prj-ar-26-go-back/internal/domain"

type RoomRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

func (r RoomRequest) ToDomainModel() (interface{}, error) {
	return domain.Room{
		Name:        r.Name,
		Description: r.Description,
	}, nil
}
