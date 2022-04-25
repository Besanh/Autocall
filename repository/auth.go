package repository

import (
	"autocall/common/model"
	IRedis "autocall/internal/redis"
	"encoding/json"
)

const redisAccessTokenUser = "access_token_user"
const redisAccessTokenKey = "access_token_key"

type AuthRepository struct {
}

func NewAuthRepository() AuthRepository {
	return AuthRepository{}
}

var AuthRepo AuthRepository

func (repo *AuthRepository) GetAccessTokenFromCache(clientId string) (interface{}, error) {
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

func (repo *AuthRepository) InsertAccessTokenCache(token model.AccessToken) error {
	clientId := token.ClienID
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

func (repo *AuthRepository) DeleteAccessTokenCache(token model.AccessToken) error {
	clientId := token.ClienID
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

func (repo *AuthRepository) GetAuthenFromCache(token string) (interface{}, error) {
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
