package auth

import (
	"autocall/common/log"
	"autocall/common/model"
	"autocall/repository"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Config struct {
	ExpiredTime int
	TokenType   string
}

var AuthConfig Config

type AuthClient struct {
	ClientID     string
	ClientSecret string
	UserId       string
	Scope        string
	User         model.User
}

type AccessTokenResponse struct {
	UserID       string `json:"user_id"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiredtime  int    `json:"expire_at"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

func NewAuthUtil(config Config) {
	AuthConfig.ExpiredTime = config.ExpiredTime
	AuthConfig.TokenType = config.TokenType
}

func ClientCredential(client AuthClient, isRefresh bool) (interface{}, error) {
	accessToken, err := CheckAccessTokenCache(client)
	if err != nil {
		return nil, err
	}
	accessTokenResponse, err := CreateResponseAccessToken(accessToken, isRefresh)
	if err != nil {
		return nil, err
	}
	return accessTokenResponse, nil
}

func CheckAccessTokenCache(client AuthClient) (model.AccessToken, error) {
	var accessToken model.AccessToken

	log.Info("Auth Middleware", "CheckAccessTokenCache - clientID", client.ClientID)
	accessTokenRes, err := repository.AuthRepo.GetAccessTokenFromCache(client.ClientID)
	if err != nil {
		log.Error("Auth Middleware", "CheckAccessTokenCache - GetAccessTokenFromCache", err)
		return accessToken, err
	}
	if accessTokenRes != "" {
		var ok bool
		accessToken, ok = accessTokenRes.(model.AccessToken)
		if !ok {
			log.Error("Auth Middleware", "CheckAccessTokenCache", err)
			return accessToken, err
		}
	}

	if accessToken.ClienID == "" {
		accessToken = CreateAccessToken(client)
		log.Info("Auth Middleware", "CheckAccessTokenCache - CreateAccessToken", accessToken)
		err := repository.AuthRepo.InsertAccessTokenCache(accessToken)
		if err != nil {
			log.Error("Auth Middleware", "CheckAccessTokenCache - InsertAccessTokenCache", err)
			return accessToken, err
		}
	} else {
		timein := accessToken.Createdtime.Add(time.Second * time.Duration(accessToken.Expiredtime))
		if timein.Sub(time.Now().Local()) <= 0 {
			log.Info("Auth Middleware", "CheckAccessTokenCache", "accesstoken already expired and create new accesstoken")
			err := repository.AuthRepo.DeleteAccessTokenCache(accessToken)
			if err != nil {
				return accessToken, err // tao error
			}
			accessToken = CreateAccessToken(client)
			log.Info("Auth Middleware", "CheckAccessTokenCache - CreateAccessToken", accessToken)
			err = repository.AuthRepo.InsertAccessTokenCache(accessToken)
			if err != nil {
				log.Error("Auth Middleware", "CheckAccessTokenCache - InsertAccessTokenCache", err)
				return accessToken, err
			}
		} else {
			log.Info("Auth Middleware", "CheckAccessTokenCache", "accesstoken already existed")
		}
	}
	return accessToken, nil
}

func CheckAccessToken(token string) (interface{}, error) {
	res, err := repository.AuthRepo.GetAuthenFromCache(token)
	if err != nil {
		log.Error("Auth Middleware", "CheckAccessToken - GetAuthenFromCache", err)
		return nil, err
	}
	accessToken, ok := res.(model.AccessToken)
	if !ok {
		return nil, errors.New("get access token fail")
	}
	// timein := accessToken.Createdtime.Add(time.Second * time.Duration(accessToken.Expiredtime))
	// if timein.Sub(time.Now().Local()) <= 0 {
	// 	// token already expired
	// 	log.Error("Auth Middleware", "CheckAccessToken", "Token expired")
	// 	return nil, errors.New("token is expired")
	// }
	return accessToken, nil
}

func CreateResponseAccessToken(token model.AccessToken, isRefresh bool) (AccessTokenResponse, error) {
	response := AccessTokenResponse{}
	if token.Token == "" {
		return response, errors.New("token is null")
	}
	if !isRefresh {
		response = AccessTokenResponse{
			UserID:      token.UserID,
			Token:       token.Token,
			Expiredtime: token.Expiredtime,
			TokenType:   AuthConfig.TokenType, // config
			Scope:       token.Scope,
		}
	} else {
		response = AccessTokenResponse{
			UserID:       token.UserID,
			Token:        token.Token,
			RefreshToken: token.RefreshToken,
			Expiredtime:  token.Expiredtime,
			TokenType:    AuthConfig.TokenType, // config
			Scope:        token.Scope,
		}
	}

	return response, nil
}

//GenerateToken
func GenerateToken(id string) string {
	uuidNew, _ := uuid.NewRandom()
	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
	token := strings.Replace(uuidNew.String(), "-", "", -1)
	token = token + "-" + idEnc
	return token
}

//GenerateRefreshToken
func GenerateRefreshToken(id string) string {
	uuidNew, _ := uuid.NewRandom()
	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
	token := strings.Replace(uuidNew.String(), "-", "", -1)
	token = token + "-" + idEnc
	return token
}

//CreateAccessToken -CreateAccessToken
func CreateAccessToken(client AuthClient) model.AccessToken {
	accesstoken := model.AccessToken{
		ClienID:      client.ClientID,
		UserID:       client.UserId,
		Token:        GenerateToken(client.ClientID),
		RefreshToken: GenerateRefreshToken(client.ClientID),
		Createdtime:  time.Now().Local(),
		Expiredtime:  AuthConfig.ExpiredTime,
		Scope:        client.Scope,
		TokenType:    AuthConfig.TokenType,
	}
	jwtData := make(map[string]string)
	jwtData["username"] = client.User.UserName
	jwtData["level"] = client.User.Level
	jwtData["group_id"] = client.User.GroupID
	accesstoken.JWT = GenerateJWT(client.UserId, jwtData)
	return accesstoken
}

func GenerateJWT(id string, data map[string]string) string {
	claim := jwt.MapClaims{
		"iss": "crm-api",
		"sub": "tel4vn",
		"aud": "tel4vn",
		"jti": id,
		"id":  id,
	}
	if len(data) > 0 {
		for key, value := range data {
			claim[key] = value
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtToken, _ := token.SignedString([]byte("secret"))
	return jwtToken
}
