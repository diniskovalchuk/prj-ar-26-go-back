package controllers

import (
	"errors"
	"log"
	"net/http"

	"prj-ar-26-go-back/internal/app"
	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/http/requests"
	"prj-ar-26-go-back/internal/infra/http/resources"
)

type RoomController struct {
	roomService app.RoomService
}

func NewRoomController(rs app.RoomService) RoomController {
	return RoomController{
		roomService: rs,
	}
}

func (c RoomController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room, err := requests.Bind(r,
			requests.RoomRequest{},
			domain.Room{})
		if err != nil {
			log.Printf("RoomController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		room.OrganizationId = user.Id

		room, err = c.roomService.Save(room)
		if err != nil {
			log.Printf("RoomController.Save(c.roomService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		roomDto := resources.RoomDto{}
		roomDto = roomDto.DomainToDto(room)
		Success(w, roomDto)
	}
}

func (c RoomController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(UserKey).(domain.User)

		rooms, err := c.roomService.FindList(organization.Id)
		if err != nil {
			log.Printf("RoomController.FindList(c.roomService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		roomList := rooms

		Success(w, resources.RoomDto{}.DomainToDtoCollection(roomList))
	}
}

func (c RoomController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		room := r.Context().Value(RoomKey).(domain.Room)

		if user.Id != room.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.RoomDto{}.DomainToDto(room))
	}
}

func (c RoomController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		room := r.Context().Value(RoomKey).(domain.Room)

		if user.Id != room.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		newRoom, err := requests.Bind(r,
			requests.RoomRequest{},
			domain.Room{})
		if err != nil {
			log.Printf("RoomController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		room.Name = newRoom.Name
		room.Description = newRoom.Description

		room, err = c.roomService.Update(room)
		if err != nil {
			log.Printf("RoomController.Update(c.roomService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.RoomDto{}.DomainToDto(room))
	}
}

func (c RoomController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		room := r.Context().Value(RoomKey).(domain.Room)

		if user.Id != room.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := c.roomService.Delete(room.Id)
		if err != nil {
			log.Printf("RoomController.Delete(c.roomService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
