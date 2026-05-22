package database

import (
	"time"

	"github.com/diniskovalchuk/prj2/internal/domain"
	"github.com/upper/db/v4"
)

const roomTableName = "rooms"

type room struct {
	Id          uint64     `db:"id,omitempty"`
	OrganizationId uint64  `db:"organization_id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

type roomRepository struct {
	coll db.Collection
	sess db.Session
}

type RoomRepository interface {
	Save(r domain.Room) (domain.Room, error)
	FindList(oId uint64) ([]domain.Room, error)
	Find(id uint64) (domain.Room, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

func NewRoomRepository(session db.Session) RoomRepository {
	return roomRepository{
		sess: session,
		coll: session.Collection(roomTableName),
	}
}

func (r roomRepository) Save(rm domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(rm)
	now := time.Now()
	room.CreatedDate = now
	room.UpdatedDate = now

	err := r.coll.InsertReturning(&room)
	if err != nil {
		return domain.Room{}, err
	}

	rm = r.mapModelToDomain(room)
	return rm, nil
}

func (r roomRepository) FindList(oId uint64) ([]domain.Room, error) {
	var rooms []room

	err := r.coll.
		Find(db.Cond{
			"organization_id": oId,
			"deleted_date": nil,
		}).
		All(&rooms)
	if err != nil {
		return nil, err
	}

	roomsList := r.mapModelToDomainCollection(rooms)
	return roomsList, nil
}

func (r roomRepository) Find(id uint64) (domain.Room, error) {
	var room room

	err := r.coll.
		Find(db.Cond{"id": id, "deleted_date": nil}).
		One(&room)
	if err != nil {
		return domain.Room{}, err
	}

	rm := r.mapModelToDomain(room)
	return rm, nil
}

func (r roomRepository) Update(rm domain.Room) (domain.Room, error) {
	room := r.mapDomainToModel(rm)
	room.UpdatedDate = time.Now()

	err := r.coll.
		Find(db.Cond{"id": rm.Id, "deleted_date": nil}).
		Update(&room)
	if err != nil {
		return domain.Room{}, err
	}

	rm = r.mapModelToDomain(room)
	return rm, nil
}

func (r roomRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r roomRepository) mapDomainToModel(rm domain.Room) room {
	return room{
		Id:          rm.Id,
		OrganizationId: rm.OrganizationId,
		Name:        rm.Name,
		Description: rm.Description,
		CreatedDate: rm.CreatedDate,
		UpdatedDate: rm.UpdatedDate,
		DeletedDate: rm.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomain(rm room) domain.Room {
	return domain.Room{
		Id:          rm.Id,
		OrganizationId: rm.OrganizationId,
		Name:        rm.Name,
		Description: rm.Description,
		CreatedDate: rm.CreatedDate,
		UpdatedDate: rm.UpdatedDate,
		DeletedDate: rm.DeletedDate,
	}
}

func (r roomRepository) mapModelToDomainCollection(rooms []room) []domain.Room {
	roomsList := make([]domain.Room, len(rooms))
	for i := range rooms {
		roomsList[i] = r.mapModelToDomain(rooms[i])
	}
	return roomsList
}
