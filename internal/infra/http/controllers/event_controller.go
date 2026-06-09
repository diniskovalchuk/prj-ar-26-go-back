package controllers

import (
	"errors"
	"log"
	"net/http"
	"prj-ar-26-go-back/internal/app"
	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/http/requests"
	"strconv"
)

type EventController struct {
	service app.EventService
}

func NewEventController(s app.EventService) EventController {
	return EventController{service: s}
}

// Save приймає запити від виконавчих пристроїв ("on"/"off")
func (c EventController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqType requests.EventRequest
		var domainModel domain.Event

		req, err := requests.Bind(r, reqType, domainModel)
		if err != nil {
			log.Printf("EventController.Save validation error: %s", err)
			BadRequest(w, err)
			return
		}

		res, err := c.service.Save(req)
		if err != nil {
			log.Printf("EventController.Save service error: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, res) // Або через DTO, якщо потрібно сховати поля
	}
}

// GetEnergyReport — аналітика для Адміністратора системи
// GET /api/v1/events/energy-report?period=week&room_id=3
func (c EventController) GetEnergyReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		if user.Role != "admin" {
			Forbidden(w, errors.New("access denied: admin only"))
			return
		}

		period := r.URL.Query().Get("period")
		if period == "" {
			period = "day"
		}

		var roomId uint64
		if roomIdStr := r.URL.Query().Get("room_id"); roomIdStr != "" {
			roomId, _ = strconv.ParseUint(roomIdStr, 10, 64)
		}

		report, err := c.service.GetEnergyReport(roomId, period)
		if err != nil {
			BadRequest(w, err)
			return
		}

		Success(w, report)
	}
}
