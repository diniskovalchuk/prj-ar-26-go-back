package app

import (
	"log"

	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/database"
)

type roomService struct {
	roomRepo database.RoomRepository
}

type RoomService interface {
	Save(r domain.Room) (domain.Room, error)
	FindList(uId uint64) ([]domain.Room, error)
	Find(id uint64) (interface{}, error)
	Update(r domain.Room) (domain.Room, error)
	Delete(id uint64) error
}

func NewRoomService(
	rr database.RoomRepository) RoomService {
	return roomService{
		roomRepo: rr,
	}
}

func (s roomService) Save(r domain.Room) (domain.Room, error) {
	room, err := s.roomRepo.Save(r)
	if err != nil {
		log.Printf("roomService.Save(s.roomRepo.Save): %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s roomService) FindList(oId uint64) ([]domain.Room, error) {
	rooms, err := s.roomRepo.FindList(oId)
	if err != nil {
		log.Printf("roomService.FindList(s.roomRepo.FindList): %s", err)
		return nil, err
	}

	return rooms, nil
}

func (s roomService) Find(id uint64) (interface{}, error) {
	room, err := s.roomRepo.Find(id)
	if err != nil {
		log.Printf("roomService.Find(s.roomRepo.Find): %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s roomService) Update(r domain.Room) (domain.Room, error) {
	room, err := s.roomRepo.Update(r)
	if err != nil {
		log.Printf("roomService.Update(s.roomRepo.Update): %s", err)
		return domain.Room{}, err
	}

	return room, nil
}

func (s roomService) Delete(id uint64) error {
	err := s.roomRepo.Delete(id)
	if err != nil {
		log.Printf("roomService.Delete(s.roomRepo.Delete): %s", err)
		return err
	}

	return nil
}
