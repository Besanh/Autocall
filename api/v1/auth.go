package v1

import (
	"autocall/common/log"
	mdw "autocall/internal/middleware"
	"autocall/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	Auth struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	AuthHandler struct {
		UserService service.IUserService
	}
)

func NewAuthHandler(r *gin.Engine, user service.IUserService) {
	handler := &AuthHandler{
		UserService: user,
	}
	Group := r.Group("v1/auth")
	{
		Group.GET("check-auth", mdw.AuthMiddleware(), handler.CheckAuth)
		Group.POST("token", handler.GenerateToken)
	}
}

/**
* Check auth user
 */
func (data *AuthHandler) CheckAuth(c *gin.Context) {
	user, isExisted := c.Get("user")
	log.Info("AuthHandler", "CheckAuth", user)
	c.JSON(http.StatusOK, gin.H{
		"isExisted": isExisted,
		"user":      user,
	})
}

func (data *AuthHandler) GenerateToken(c *gin.Context) {
	log.Info("AuthHandler", "GenerateToken", "start")
	var body map[string]interface{}
	err := c.BindJSON(&body)
	if err != nil {
		log.Error("AuthHandler", "GenerateToken - error", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	apiKey, ok := body["api_key"].(string)
	if !ok || apiKey == "" {
		log.Error("AuthHandler", "GenerateToken - error", body)
		c.JSON(http.StatusBadRequest, "api_key must be not null")
		return
	}
	isRefresh, ok := body["refresh"].(bool)
	if !ok {
		log.Info("AuthHandler", "GenerateToken - error", isRefresh)
		isRefresh = false
	}
	code, result := data.UserService.GenerateTokenByApiKey(apiKey, isRefresh)
	c.JSON(code, result)
	return
}
