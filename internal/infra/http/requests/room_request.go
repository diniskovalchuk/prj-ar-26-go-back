package requests

import "github.com/diniskovalchuk/prj2/internal/domain"

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
