package app

import (
	"log"

	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/database"
)

type deviceService struct {
	devRepo database.DeviceRepository
}

type DeviceService interface {
	Save(o domain.Device) (domain.Device, error)
	FindList(uId uint64) ([]domain.Device, error)
	Find(id uint64) (interface{}, error)
	Update(o domain.Device) (domain.Device, error)
	Delete(id uint64) error
}

func NewDeviceService(
	devRepo database.DeviceRepository,
) deviceService {
	return deviceService{
		devRepo: devRepo,
	}
}

func (s deviceService) Save(o domain.Device) (domain.Device, error) {
	dev, err := s.devRepo.Save(o)
	if err != nil {
		log.Printf("deviceService.Save(s.devRepo.Save): %s", err)
		return domain.Device{}, err
	}

	return dev, nil
}

func (s deviceService) FindList(uId uint64) ([]domain.Device, error) {
	devs, err := s.devRepo.FindList(uId)
	if err != nil {
		log.Printf("deviceService.FindList(s.devRepo.FindList): %s", err)
		return nil, err
	}

	return devs, nil
}

func (s deviceService) Find(id uint64) (interface{}, error) {
	dev, err := s.devRepo.Find(id)
	if err != nil {
		log.Printf("deviceService.Find(s.devRepo.Find): %s", err)
		return nil, err
	}

	return dev, nil
}

func (s deviceService) Update(o domain.Device) (domain.Device, error) {
	dev, err := s.devRepo.Update(o)
	if err != nil {
		log.Printf("deviceService.Update(s.devRepo.Update): %s", err)
		return domain.Device{}, err
	}

	return dev, nil
}

func (s deviceService) Delete(id uint64) error {
	err := s.devRepo.Delete(id)
	if err != nil {
		log.Printf("deviceService.Delete(s.devRepo.Delete): %s", err)
		return err
	}

	return nil
}
