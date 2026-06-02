package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/diniskovalchuk/prj2/internal/app"
	"github.com/diniskovalchuk/prj2/internal/domain"
	"github.com/diniskovalchuk/prj2/internal/infra/http/requests"
	"github.com/diniskovalchuk/prj2/internal/infra/http/resources"
)

type DeviceController struct {
	devService app.DeviceService
}

func NewDeviceController(ds app.OrganizationService) DeviceController {
	return DeviceController{
		devService: ds,
	}
}

func (c DeviceController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dev, err := requests.Bind(r,
			requests.DeviceRequest{},
			domain.Device{})
		if err != nil {
			log.Printf("DeviceController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		organization := r.Context().Value(OrgKey).(domain.Organization)
		dev.OrganizationId = organizationId.Id

		dev, err = c.devService.Save(dev)
		if err != nil {
			log.Printf("DeviceController.Save(c.orgService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		devDto := resources.DeviceDto{}
		devDto = devDto.DomainToDto(dev)
		Success(w, devDto)
	}
}

func (c DeviceController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)

		dev, err := c.devService.FindList(organization.Id)
		if err != nil {
			log.Printf("DeviceController.FindList(c.orgService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDtoCollection(dev))
	}
}

func (c DeviceController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(OrgKey).(domain.Organization)
		dev := r.Context().Value(DevKey).(domain.Device)

		if organization.Id != dev.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDto(dev))
	}
}

func (c DeviceController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(UserKey).(domain.User)
		dev := r.Context().Value(DevKey).(domain.Device)

		if organization.Id != dev.OrganizationId {
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

		dev, err = c.devService.Update(dev)
		if err != nil {
			log.Printf("DeviceController.Update(c.orgService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.DeviceDto{}.DomainToDto(dev))
	}
}

func (c DeviceController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := r.Context().Value(UserKey).(domain.User)
		dev := r.Context().Value(DevKey).(domain.Device)

		if organization.Id != dev.OrganizationId {
			Forbidden(w, errors.New("access denied"))
			return
		}

		err := c.devService.Delete(dev.Id)
		if err != nil {
			log.Printf("DeviceController.Delete(c.devService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
