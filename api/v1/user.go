package v1

import (
	"autocall/common/log"
	authMdw "autocall/middleware/auth"
	"autocall/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.IUserService
}

func NewUserHandler(r *gin.Engine, user service.IUserService) {
	handler := &UserHandler{
		UserService: user,
	}

	Group := r.Group("v1/user")
	{
		Group.GET("user-in-group/:groupID", authMdw.AuthMiddleware(), authMdw.CheckAdmin(), handler.GetUserInGroup)

	}
}

func (data *UserHandler) GetUserInGroup(c *gin.Context) {
	groupID := c.Param("groupID")
	limit := c.Param("limit")
	offset := c.Param("offset")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Error("UserHandler", "GetUserInGroup - error", "limit is null or fail")
		c.JSON(http.StatusNotFound, err)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Error("UserHandler", "GetUserInGroup - error", "offset is null or fail")
		c.JSON(http.StatusNotFound, err)
		return
	}
	code, result := data.UserService.GetUserInGroup(groupID, limitInt, offsetInt)
	if groupID == "" {
		log.Info("UserHandler", "GetUserInGroup", "groupID is null")
		return
	}
	c.JSON(code, result)
}
