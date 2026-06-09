package app

import (
	"errors"
	"log"
	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/database"
	"time"
)

type EventService interface {
	Save(e domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error)
	Update(e domain.Event) (domain.Event, error)
	Delete(id uint64) error
	GetEnergyReport(roomId uint64, period string) (domain.EnergyReport, error)
}

type eventService struct {
	repo database.EventRepository
}

func NewEventService(repo database.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) Save(e domain.Event) (domain.Event, error) {
	if e.DeviceId == 0 || e.RoomId == 0 || e.Action == "" {
		return domain.Event{}, errors.New("invalid event data")
	}
	return s.repo.Save(e)
}

func (s *eventService) Find(id uint64) (interface{}, error) {
	return s.repo.Find(id)
}

func (s *eventService) Update(e domain.Event) (domain.Event, error) {
	return s.repo.Update(e)
}

func (s *eventService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *eventService) GetEnergyReport(roomId uint64, period string) (domain.EnergyReport, error) {
	var fromTime time.Time
	now := time.Now()

	switch period {
	case "day":
		fromTime = now.AddDate(0, 0, -1)
	case "week":
		fromTime = now.AddDate(0, 0, -7)
	case "month":
		fromTime = now.AddDate(0, -1, 0)
	default:
		return domain.EnergyReport{}, errors.New("invalid period: use day, week or month")
	}

	events, err := s.repo.GetEventsForAnalytics(roomId, fromTime)
	if err != nil {
		log.Printf("eventService.GetEnergyReport error: %s", err)
		return domain.EnergyReport{}, err
	}

	// Алгоритм підрахунку енергії:
	// Для демонстрації рахуємо час між "on" та "off" для кожного девайса.
	// Припускаємо, що середній девайс споживає умовно 0.5 кВт за годину роботи.
	var totalEnergy float64
	deviceStates := make(map[uint64]time.Time)

	for _, event := range events {
		if event.Action == "on" {
			deviceStates[event.DeviceId] = event.CreatedDate
		} else if event.Action == "off" {
			if onTime, exists := deviceStates[event.DeviceId]; exists {
				duration := event.CreatedDate.Sub(onTime).Hours()
				totalEnergy += duration * 0.5 // 0.5 kW per hour
				delete(deviceStates, event.DeviceId)
			}
		}
	}

	// Якщо пристрій включився в цьому періоді, але ще не виключився, порахуємо роботу до поточного моменту
	for _, onTime := range deviceStates {
		duration := now.Sub(onTime).Hours()
		totalEnergy += duration * 0.5
	}

	return domain.EnergyReport{
		Period:      period,
		RoomId:      roomId,
		EnergyUsed:  totalEnergy,
		EventsCount: len(events),
	}, nil
}
