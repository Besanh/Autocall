package v1

import (
	"autocall/common/log"
	"autocall/common/response"
	"autocall/service"

	"github.com/gin-gonic/gin"
)

type ContactCallHandler struct {
	ContactService service.IContactService
}

func NewContactCallHandler(r *gin.Engine, contact service.IContactService) {
	handler := &ContactCallHandler{
		ContactService: contact,
	}
	Group := r.Group("v1/contact")
	{
		Group.GET("/:id", handler.GetContactByID)
		// Group.GET("/", handler.GetAllContacts)
	}
}

func (data *ContactCallHandler) GetContactByID(c *gin.Context) {
	id := c.Param("id")
	log.Info("ContactHandler", "GetContactByID", id)
	if id == "" {
		c.JSON(response.BadRequestMsg("Required id"))
		return
	}
	code, result := data.ContactService.GetContactByID(id)
	c.JSON(code, result)
}

// func (data *ContactCallHandler) GetAllContacts(c *gin.Context) {
// 	log.Info("ContactHandler", "GetAllContacts", "ok")
// 	code, result := da
// }
