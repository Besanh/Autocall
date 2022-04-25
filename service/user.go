package service

import (
	"autocall/common/auth"
	encryptUtil "autocall/common/encrypt"
	"autocall/common/log"
	"autocall/common/model"
	"autocall/repository"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type (
	IUserService interface {
		CreateUser(entry model.User) (int, interface{})
		GenerateTokenByApiKey(apiKey string, isRefresh bool) (int, interface{})
		GetUserInGroup(groupID string, limit, offset int) (int, interface{})
	}
	UserService struct {
		repo repository.User
	}
)

func NewUserService() IUserService {
	return &UserService{}
}

func (data *UserService) CreateUser(entry model.User) (int, interface{}) {
	log.Info("UserService", "Auth", entry)
	entry.APIKey = uuid.NewString()
	entry.CreatedAt = time.Now()
	entry.Password, _ = encryptUtil.GenerateEncrypted(entry.Password)

	result, err := data.repo.CreateUser(entry)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusCreated, result
}

func (data *UserService) GenerateTokenByApiKey(apiKey string, isRefresh bool) (int, interface{}) {
	log.Info("UserService", "GenerateTokenByApiKey", apiKey)
	resp, err := data.repo.GetUserByApiKey(apiKey)
	if err != nil {
		log.Error("UserService", "GenerateTokenByApiKey - error", err)
		return http.StatusNotFound, err
	}
	if resp == nil {
		log.Error("UserService", "GenerateTokenByApiKey", "User is null")
		return http.StatusNotFound, errors.New("User is null")
	}
	user := resp.(model.User)
	log.Info("UserService", "GenerateTokenByApiKey - user", user)
	clientAuth := auth.AuthClient{
		ClientID:     apiKey,
		ClientSecret: apiKey,
		User:         user,
	}
	token, err := auth.ClientCredential(clientAuth, isRefresh)
	if err != nil {
		log.Error("UserService", "GenerateTokenByApiKey - ClientCredential - error", err)
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, token
}

func (data *UserService) GetUserInGroup(groupID string, limit, offset int) (int, interface{}) {
	log.Info("UserService", "GetUserInGroup", groupID)
	resp, total, err := data.repo.SelectUserInGroup(groupID, limit, offset)
	if err != nil {
		log.Error("UserService", "GetUserInGroup - SelectUserInGroup - error", err)
		return http.StatusNotFound, err
	}
	if total == 0 {
		log.Debug("UserService", "GetUserInGroup - SelectUserInGroup - total", 0)
		return http.StatusNotFound, nil
	}
	return http.StatusOK, map[string]interface{}{
		"total":   total,
		"success": resp,
	}
}
