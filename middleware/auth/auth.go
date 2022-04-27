package auth

import (
	"autocall/common/log"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	AuthUrl string
)

type AuthUser struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Level   string   `json:"level"`
	Scopes  []string `json:"scopes"`
	GroupID string   `json:"group_id"`
}

func (a *AuthUser) SetId(id string) {
	a.Id = id
}

func (a *AuthUser) GetId() string {
	return a.Id
}

func (a *AuthUser) SetName(name string) {
	a.Name = name
}

func (a *AuthUser) GetName() string {
	return a.Name
}

func (a *AuthUser) SetLevel(level string) {
	a.Level = level
}

func (a *AuthUser) GetLevel() string {
	return a.Level
}

func (a *AuthUser) SetScopes(scopes []string) {
	a.Scopes = scopes
}

func (a *AuthUser) GetScopes() []string {
	return a.Scopes
}

func (a *AuthUser) SetGroupID(groupID string) {
	a.GroupID = groupID
}

func (a *AuthUser) GetGroupID() string {
	return a.GroupID
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if len(token) < 1 {
			log.Error("AuthMiddleware", "validateBasicAuth", "invalid credentials")
			c.JSON(
				http.StatusUnauthorized,
				map[string]interface{}{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
		AuthUser, err := PostToAuthAPI(token)
		if err != nil {
			log.Error("AuthMiddleware", "validateBasicAuth", err.Error())
			c.JSON(
				http.StatusUnauthorized,
				map[string]interface{}{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
		c.Set("user", AuthUser)
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel, ok := GetUserLevel(c)
		if !ok || (userLevel != "admin" && userLevel != "superadmin") {
			c.JSON(
				http.StatusForbidden,
				map[string]interface{}{
					"error": http.StatusText(http.StatusForbidden),
				},
			)
			c.Abort()
			return
		}
	}
}

func PostToAuthAPI(token string) (*AuthUser, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, err := http.NewRequest("POST", AuthUrl, nil)
	if err != nil {
		log.Error("Auth", "PostToAuthAPI", err)
		return nil, err
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	client.Transport = tr
	res, err := client.Do(req)
	if err != nil {
		log.Error("Auth", "PostToAuthAPI", err)
		return nil, err
	}
	defer res.Body.Close()
	log.Info("Auth", "PostToAuthAPI", res.StatusCode)
	AuthUser := new(AuthUser)

	err = json.NewDecoder(res.Body).Decode(AuthUser)
	if err != nil {
		log.Error("Auth", "PostToAuthAPI", err)
		return nil, err
	}
	return AuthUser, nil
}

func GetUser(c *gin.Context) (*AuthUser, bool) {
	tmp, isExist := c.Get("user")
	if isExist {
		user, ok := tmp.(*AuthUser)
		return user, ok
	} else {
		return nil, false
	}
}

func GetUserId(c *gin.Context) (string, bool) {
	user, ok := GetUser(c)
	if !ok {
		return "", false
	} else {
		return user.Id, true
	}
}

func GetUserLevel(c *gin.Context) (string, bool) {
	user, ok := GetUser(c)
	if !ok {
		return "", false
	} else {
		return user.Level, true
	}
}

func GetUserName(c *gin.Context) (string, bool) {
	user, ok := GetUser(c)
	if !ok {
		return "", false
	} else {
		return user.Name, true
	}
}

func GetUserGroupID(c *gin.Context) (string, bool) {
	user, ok := GetUser(c)
	if !ok {
		return "", false
	}
	return user.GroupID, true
}
