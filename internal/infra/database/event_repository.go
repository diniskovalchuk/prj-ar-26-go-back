package database

import (
	"prj-ar-26-go-back/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

type Event struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      uint64     `db:"room_id"`
	Action      string     `db:"action"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

const EventTableName = "events"

type eventRepository struct {
	coll db.Collection
	sess db.Session
}

type EventRepository interface {
	Save(e domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error) // Для сумісності з вашим PathObject middleware
	Update(e domain.Event) (domain.Event, error)
	Delete(id uint64) error
	// Метод для аналітики адміністратора
	GetEventsForAnalytics(roomId uint64, fromTime time.Time) ([]domain.Event, error)
}

func NewEventRepository(session db.Session) EventRepository {
	return eventRepository{
		sess: session,
		coll: session.Collection(EventTableName),
	}
}

func (r eventRepository) Save(e domain.Event) (domain.Event, error) {
	model := r.mapDomainToModel(e)
	now := time.Now()
	model.CreatedDate = now
	model.UpdatedDate = now

	err := r.coll.InsertReturning(&model)
	if err != nil {
		return domain.Event{}, err
	}
	return r.mapModelToDomain(model), nil
}

func (r eventRepository) Find(id uint64) (interface{}, error) {
	var model Event
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&model)
	if err != nil {
		return nil, err
	}
	return r.mapModelToDomain(model), nil
}

func (r eventRepository) Update(e domain.Event) (domain.Event, error) {
	model := r.mapDomainToModel(e)
	model.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": model.Id, "deleted_date": nil}).Update(&model)
	if err != nil {
		return domain.Event{}, err
	}
	return r.mapModelToDomain(model), nil
}

func (r eventRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r eventRepository) GetEventsForAnalytics(roomId uint64, fromTime time.Time) ([]domain.Event, error) {
	var models []Event
	cond := db.Cond{
		"deleted_date":    nil,
		"created_date >=": fromTime,
	}
	// Якщо roomId == 0, беремо по всьому підприємству. Якщо задано — фільтруємо по кімнаті.
	if roomId > 0 {
		cond["room_id"] = roomId
	}

	err := r.coll.Find(cond).OrderBy("created_date ASC").All(&models)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(models), nil
}

// Мапери
func (r eventRepository) mapDomainToModel(e domain.Event) Event {
	return Event{Id: e.Id, DeviceId: e.DeviceId, RoomId: e.RoomId, Action: e.Action, CreatedDate: e.CreatedDate, UpdatedDate: e.UpdatedDate, DeletedDate: e.DeletedDate}
}
func (r eventRepository) mapModelToDomain(e Event) domain.Event {
	return domain.Event{Id: e.Id, DeviceId: e.DeviceId, RoomId: e.RoomId, Action: e.Action, CreatedDate: e.CreatedDate, UpdatedDate: e.UpdatedDate, DeletedDate: e.DeletedDate}
}
func (r eventRepository) mapModelToDomainCollection(models []Event) []domain.Event {
	events := make([]domain.Event, len(models))
	for i := range models {
		events[i] = r.mapModelToDomain(models[i])
	}
	return events
}
