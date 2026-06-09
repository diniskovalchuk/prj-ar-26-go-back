package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"prj-ar-26-go-back/internal/app"
	"prj-ar-26-go-back/internal/domain"
	"prj-ar-26-go-back/internal/infra/http/requests"
	"prj-ar-26-go-back/internal/infra/http/resources"
)

type MeasurementController struct {
	measService app.MeasurementService
}

func NewMeasurementController(ms app.MeasurementService) MeasurementController {
	return MeasurementController{
		measService: ms,
	}
}

// Save приймає дані від вимірювальних пристроїв і записує їх до БД.
func (c MeasurementController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqType requests.MeasurementRequest
		var domainModel domain.Measurement

		// Bind зчитує DeviceId, RoomId та Value з JSON-тіла запиту від пристрою
		req, err := requests.Bind[requests.MeasurementRequest, domain.Measurement](r, reqType, domainModel)
		if err != nil {
			log.Printf("MeasurementController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		// Зберігаємо вимірювання (DeviceId береться з розпарсеного req, а не з користувача)
		req, err = c.measService.Save(req)
		if err != nil {
			log.Printf("MeasurementController.Save(c.measService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		measurementDto := resources.MeasurementDto{}
		Success(w, measurementDto.DomainToDto(req))
	}
}

// GetAdminReport реалізує вимогу перегляду даних адміном за день/тиждень/місяць.
// Очікує query-параметри: ?device_id=X&period=day|week|month
func (c MeasurementController) GetAdminReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Перевірка ролі користувача (ТЗ: "Адміністратор системи має мати змогу...")
		user := r.Context().Value(UserKey).(domain.User)
		if user.Role != "admin" { // Перевірте, як саме у вашому domain.User реалізовані ролі
			Forbidden(w, errors.New("only administrators can access this report"))
			return
		}

		// Зчитування та валідація query-параметра device_id
		deviceIDStr := r.URL.Query().Get("device_id")
		if deviceIDStr == "" {
			BadRequest(w, errors.New("missing device_id parameter"))
			return
		}
		deviceID, err := strconv.ParseUint(deviceIDStr, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid device_id parameter"))
			return
		}

		// Зчитування періоду (day, week, month)
		period := r.URL.Query().Get("period")
		if period == "" {
			period = "day" // Значення за замовчуванням
		}

		// Отримання фільтрованого списку з сервісу
		meas, err := c.measService.GetByDeviceAndPeriod(deviceID, period)
		if err != nil {
			log.Printf("MeasurementController.GetAdminReport error: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.MeasurementDto{}.DomainToDtoCollection(meas))
	}
}

func (c MeasurementController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримуємо умовний roomId з query параметрів, щоб повернути вимірювання для кімнати
		roomIDStr := r.URL.Query().Get("room_id")
		roomID, _ := strconv.ParseUint(roomIDStr, 10, 64)

		meas, err := c.measService.GetByDeviceAndPeriod(roomID, "day") // Або використайте свій кастомний метод
		if err != nil {
			log.Printf("MeasurementController.FindList error: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.MeasurementDto{}.DomainToDtoCollection(meas))
	}
}

func (c MeasurementController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meas := r.Context().Value(MesKey).(domain.Measurement)
		Success(w, resources.MeasurementDto{}.DomainToDto(meas))
	}
}

func (c MeasurementController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meas := r.Context().Value(MesKey).(domain.Measurement)

		var reqType requests.MeasurementRequest
		var domainModel domain.Measurement
		req, err := requests.Bind[requests.MeasurementRequest, domain.Measurement](r, reqType, domainModel)
		if err != nil {
			log.Printf("MeasurementController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		meas.DeviceId = req.DeviceId
		meas.RoomId = req.RoomId
		meas.Value = req.Value

		meas, err = c.measService.Update(meas)
		if err != nil {
			log.Printf("MeasurementController.Update(c.measService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.MeasurementDto{}.DomainToDto(meas))
	}
}

func (c MeasurementController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meas := r.Context().Value(MesKey).(domain.Measurement)

		err := c.measService.Delete(meas.Id)
		if err != nil {
			log.Printf("MeasurementController.Delete(c.measService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
