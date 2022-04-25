package util

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHeader(c *gin.Context, name string) (string, bool) {
	name = strings.Title(name)
	header := c.Request.Header[name]
	if len(header) < 1 {
		return "", false
	}
	if header[0] == "" {
		return "", false
	}
	return header[0], true
}
