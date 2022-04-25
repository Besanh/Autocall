package elastic

import (
	// "io"
	"net/http"
)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       map[string]interface{}
}
