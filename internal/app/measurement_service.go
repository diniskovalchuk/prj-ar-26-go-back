package app

import (
	"errors"
	"log"
	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/database"
	"time"
)

type MeasurementService interface {
	Save(m domain.Measurement) (domain.Measurement, error)
	Find(id uint64) (interface{}, error)
	Update(m domain.Measurement) (domain.Measurement, error)
	Delete(id uint64) error

	// Розширений метод для виконання вимоги ТЗ (перегляд за пристроєм та періодом)
	GetByDeviceAndPeriod(deviceID uint64, period string) ([]domain.Measurement, error)
}

type measurementService struct {
	mesRepo database.MeasurementRepository
}

// NewMeasurementService тепер повертає інтерфейс MeasurementService, а не структуру.
func NewMeasurementService(mesRepo database.MeasurementRepository) MeasurementService {
	return &measurementService{
		mesRepo: mesRepo,
	}
}

// Save приймає дані від сенсорів та записує їх у БД.
// Save приймає дані від вимірювальних пристроїв і записує їх до бази даних.
func (s *measurementService) Save(m domain.Measurement) (domain.Measurement, error) {
	// Змініть m.DeviceID на m.DeviceId та m.RoomID на m.RoomId
	if m.DeviceId == 0 || m.RoomId == 0 {
		return domain.Measurement{}, errors.New("invalid device or room id")
	}

	mess, err := s.mesRepo.Save(m)
	if err != nil {
		log.Printf("measurementService.Save error: %s", err)
		return domain.Measurement{}, err
	}

	return mess, nil
}

// Find шукає конкретне вимірювання за його ID (повертає конкретний тип замість interface{}).
func (s *measurementService) Find(id uint64) (interface{}, error) { // І тут
	mess, err := s.mesRepo.Find(id)
	if err != nil {
		log.Printf("measurementService.Find error: %s", err)
		return nil, err
	}
	return mess, nil
}

// Update оновлює дані вимірювання.
func (s *measurementService) Update(m domain.Measurement) (domain.Measurement, error) {
	mes, err := s.mesRepo.Update(m)
	if err != nil {
		log.Printf("measurementService.Update error: %s", err)
		return domain.Measurement{}, err
	}

	return mes, nil
}

// Delete реалізує видалення вимірювання (зазвичай софт-деліт через deleted_date).
func (s *measurementService) Delete(id uint64) error {
	err := s.mesRepo.Delete(id)
	if err != nil {
		log.Printf("measurementService.Delete error: %s", err)
		return err
	}

	return nil
}

// GetByDeviceAndPeriod реалізує вимогу ТЗ щодо фільтрації адміністратором за день/тиждень/місяць.
func (s *measurementService) GetByDeviceAndPeriod(deviceID uint64, period string) ([]domain.Measurement, error) {
	var fromTime time.Time
	now := time.Now()

	// Обчислюємо початкову точку часу залежно від періоду
	switch period {
	case "day":
		fromTime = now.AddDate(0, 0, -1)
	case "week":
		fromTime = now.AddDate(0, 0, -7)
	case "month":
		fromTime = now.AddDate(0, -1, 0)
	default:
		return nil, errors.New("invalid period: must be day, week or month")
	}

	// Виклик репозиторію з передачею розрахованого часу відсікання
	mess, err := s.mesRepo.FindListByDeviceAndTime(deviceID, fromTime)
	if err != nil {
		log.Printf("measurementService.GetByDeviceAndPeriod error: %s", err)
		return nil, err
	}

	return mess, nil
}
