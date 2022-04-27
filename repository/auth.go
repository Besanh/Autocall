package repository

import (
	"autocall/common/model"
	IRedis "autocall/internal/redis"
	"encoding/json"
)

const redisAccessTokenUser = "access_token_user"
const redisAccessTokenKey = "access_token_key"

type (
	IAuth interface {
		GetAccessTokenFromCache(clientId string) (interface{}, error)
		InsertAccessTokenCache(token model.AccessToken) error
		DeleteAccessTokenCache(token model.AccessToken) error
		GetAuthenFromCache(token string) (interface{}, error)
	}
	Auth struct{}
)

var AuthRepo Auth

func NewAuthRepository() Auth {
	return Auth{}
}

func (repo *Auth) GetAccessTokenFromCache(clientId string) (interface{}, error) {
	res, err := IRedis.Redis.HMGet(redisAccessTokenUser, clientId)
	if err != nil {
		return nil, err
	}
	accessTokenResponse := model.AccessToken{}
	if len(res) == 0 {
		return nil, nil
	} else {
		accessToken, ok := res[0].(string)
		if ok {
			err := json.Unmarshal([]byte(accessToken), &accessTokenResponse)
			if err != nil {
				return nil, err
			}
		}
		return accessTokenResponse, nil
	}
}

func (repo *Auth) InsertAccessTokenCache(token model.AccessToken) error {
	clientId := token.ClientID
	accessToken := token.Token
	jsonEncodeToken, err := json.Marshal(token)
	if err != nil {
		return err
	}
	jsonEncodeString := string(jsonEncodeToken)
	clientStoreInfo := map[string]interface{}{clientId: jsonEncodeString}
	accessTokenStoreInfo := map[string]interface{}{accessToken: jsonEncodeString}
	err = IRedis.Redis.HMSet(redisAccessTokenUser, clientStoreInfo)
	if err != nil {
		return err
	}
	err = IRedis.Redis.HMSet(redisAccessTokenKey, accessTokenStoreInfo)
	if err != nil {
		return err
	}
	return err
}

func (repo *Auth) DeleteAccessTokenCache(token model.AccessToken) error {
	clientId := token.ClientID
	accessToken := token.Token
	err := IRedis.Redis.HMDel(redisAccessTokenUser, clientId)
	if err != nil {
		return err
	}
	err = IRedis.Redis.HMDel(redisAccessTokenKey, accessToken)
	if err != nil {
		return err
	}
	return err
}

func (repo *Auth) GetAuthenFromCache(token string) (interface{}, error) {
	res, err := IRedis.Redis.HMGet(redisAccessTokenKey, token)
	if err != nil {
		return nil, err
	}
	accessTokenResponse := model.AccessToken{}
	if len(res) == 0 {
		return nil, nil
	} else {
		accessToken, ok := res[0].(string)
		if ok {
			err := json.Unmarshal([]byte(accessToken), &accessTokenResponse)
			if err != nil {
				return nil, err
			}
		}
		return accessTokenResponse, nil
	}
}
