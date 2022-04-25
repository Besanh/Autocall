package service

import (
	"autocall/common/log"
	"autocall/repository"
	"net/http"
)

type (
	IContactService interface {
		GetContactByID(id string) (int, interface{})
	}
	ContactService struct {
		repo repository.Contact
	}
)

func NewContactService() IContactService {
	return &ContactService{}
}

func (data *ContactService) GetContactByID(id string) (int, interface{}) {
	log.Info("ContactService", "GetContactByID", id)
	resp, err := data.repo.SelectContactByID(id, "6", "6421507", "2")
	if err != nil {
		log.Error("ContactService", "SelectContactByID error", err)
		// return response.ServiceUnavailableMsg(err.Error())
		return http.StatusNotFound, err.Error()
	}
	if resp == nil {
		return http.StatusNotFound, nil
	}
	log.Info("ContactService", "GetContactByID - SelectContactByID", resp)
	return http.StatusOK, resp
}
