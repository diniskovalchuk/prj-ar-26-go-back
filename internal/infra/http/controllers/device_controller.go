package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/diniskovalchuk/prj2/internal/app"
	"github.com/diniskovalchuk/prj2/internal/domain"
)

type DeviceController struct {
	orgService app.OrganizationService
}

func NewDeviceController(os app.OrganizationService) DeviceController {
	return DeviceController{
		orgService: os,
	}
}

func (c DeviceController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		org, err := requests.Bind(r,
			requests.OrganizationRequest{},
			domain.Organization{})
		if err != nil {
			log.Printf("DeviceController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		org.UserId = user.Id

		org, err = c.orgService.Save(org)
		if err != nil {
			log.Printf("DeviceController.Save(c.orgService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		devDto := resources.DeviceDto{}
		devDto = devDto.DomainToDto(org)
		Success(w, devDto)
	}
}

func (c DeviceController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		orgs, err := c.orgService.FindList(user.Id)
		if err != nil {
			log.Printf("DeviceController.FindList(c.orgService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDtoCollection(orgs))
	}
}

func (c DeviceController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDto(org))
	}
}

func (c DeviceController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		newDev, err := requests.Bind(r,
			requests.DeviceRequest{},
			domain.Device{})
		if err != nil {
			log.Printf("DeviceController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		dev.Name = newDev.Name
		dev.Description = newDev.Description
		dev.City = newDev.City
		dev.Address = newDev.Address
		dev.Lat = newDev.Lat
		dev.Lon = newDev.Lon

		org, err = c.orgService.Update(org)
		if err != nil {
			log.Printf("DeviceController.Update(c.orgService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDto(org))
	}
}

func (c DeviceController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		org := r.Context().Value(OrgKey).(domain.Organization)

		if user.Id != org.UserId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := c.orgService.Delete(org.Id)
		if err != nil {
			log.Printf("OrganizationController.Delete(c.orgService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
