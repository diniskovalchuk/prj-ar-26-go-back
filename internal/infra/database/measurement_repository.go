package database

import (
	"time"

	"prj-ar-26-go-back/internal/domain"

	"github.com/upper/db/v4"
)

type Measurement struct {
	Id          uint64     `db:"id,omitempty"`
	DeviceId    uint64     `db:"device_id"`
	RoomId      uint64     `db:"room_id"`
	Value       float64    `db:"value"`
	CreatedDate time.Time  `db:"created_date"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date"`
}

const MeasurementTableName = "measurements"

type measurementRepository struct {
	coll db.Collection
	sess db.Session
}

type MeasurementRepository interface {
	Save(m domain.Measurement) (domain.Measurement, error)
	FindList(roomId uint64) ([]domain.Measurement, error) // Повертає список вимірювань для конкретної кімнати
	Find(id uint64) (interface{}, error)
	Update(m domain.Measurement) (domain.Measurement, error)
	Delete(id uint64) error

	// Новий метод для виконання вимог ТЗ адміністратора
	FindListByDeviceAndTime(deviceId uint64, fromTime time.Time) ([]domain.Measurement, error)
}

func NewMeasurementRepository(session db.Session) MeasurementRepository {
	return measurementRepository{
		sess: session,
		coll: session.Collection(MeasurementTableName),
	}
}

func (r measurementRepository) Save(m domain.Measurement) (domain.Measurement, error) {
	meas := r.mapDomainToModel(m)
	now := time.Now()
	meas.CreatedDate = now
	meas.UpdatedDate = now

	err := r.coll.InsertReturning(&meas)
	if err != nil {
		return domain.Measurement{}, err
	}

	m = r.mapModelToDomain(meas)
	return m, nil
}

// FindList тепер коректно шукає вимірювання за roomId
func (r measurementRepository) FindList(roomId uint64) ([]domain.Measurement, error) {
	var meas []Measurement

	err := r.coll.
		Find(db.Cond{
			"room_id":      roomId,
			"deleted_date": nil,
		}).
		All(&meas)
	if err != nil {
		return nil, err
	}

	measurements := r.mapModelToDomainCollection(meas)
	return measurements, nil
}

func (r measurementRepository) Find(id uint64) (interface{}, error) { // І тут
	var meas Measurement
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&meas)
	if err != nil {
		return nil, err
	}
	return r.mapModelToDomain(meas), nil
}
func (r measurementRepository) Update(m domain.Measurement) (domain.Measurement, error) {
	meas := r.mapDomainToModel(m)
	meas.UpdatedDate = time.Now()

	err := r.coll.
		Find(db.Cond{"id": meas.Id, "deleted_date": nil}).
		Update(&meas)
	if err != nil {
		return domain.Measurement{}, err
	}

	m = r.mapModelToDomain(meas)
	return m, nil
}

func (r measurementRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

// FindListByDeviceAndTime реалізує фільтрацію для адміна (знаходить записи, створені після fromTime)
func (r measurementRepository) FindListByDeviceAndTime(deviceId uint64, fromTime time.Time) ([]domain.Measurement, error) {
	var meas []Measurement

	err := r.coll.
		Find(db.Cond{
			"device_id":       deviceId,
			"deleted_date":    nil,
			"created_date >=": fromTime, // Використання оператора порівняння в upper/db
		}).
		All(&meas)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(meas), nil
}

func (r measurementRepository) mapDomainToModel(m domain.Measurement) Measurement {
	return Measurement{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r measurementRepository) mapModelToDomain(m Measurement) domain.Measurement {
	return domain.Measurement{
		Id:          m.Id,
		DeviceId:    m.DeviceId,
		RoomId:      m.RoomId,
		Value:       m.Value,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r measurementRepository) mapModelToDomainCollection(meas []Measurement) []domain.Measurement {
	measurements := make([]domain.Measurement, len(meas))
	for i := range meas {
		measurements[i] = r.mapModelToDomain(meas[i])
	}
	return measurements
}
