package request

import (
	"link-manager/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		response.Json(*w, err.Error(), 400)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		response.Json(*w, err.Error(), 400)
		return nil, err
	}
	return &body, nil
}
