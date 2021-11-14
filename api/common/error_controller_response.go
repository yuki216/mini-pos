package common

import "net/http"

type errorControllerResponseCode string


//ControllerResponse default payload response
type ControllerResponse struct {
	Code    int `json:"code"`
	Message string                      `json:"message"`
	Data    interface{}                 `json:"data"`
}

//NewBadRequestResponse bad request format response
func NewBadRequestResponse() (int, ControllerResponse) {
	return http.StatusBadRequest, ControllerResponse{
		http.StatusBadRequest,
		"Bad request",
		map[string]interface{}{},
	}
}

//NewForbiddenResponse default for Forbidden error response
func NewForbiddenResponse() (int, ControllerResponse) {
	return http.StatusForbidden, ControllerResponse{
		http.StatusForbidden,
		"Forbidden",
		map[string]interface{}{},
	}
}

//NewForbiddenResponse default for Forbidden error response
func NewAuthorizeResponse() (int, ControllerResponse) {
	return http.StatusUnauthorized, ControllerResponse{
		http.StatusUnauthorized,
		"UnAuthorized",
		map[string]interface{}{},
	}
}

func NewErrorResponse(controll ControllerResponse)(int, ControllerResponse) {
	return http.StatusBadRequest, ControllerResponse{
		controll.Code,
		controll.Message,
		map[string]interface{}{},
	}
}
