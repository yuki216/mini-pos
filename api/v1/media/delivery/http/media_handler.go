package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mini_pos/api/common"
	middlwr "mini_pos/api/middleware"
	"mini_pos/config"
	"mini_pos/constants"
	"mini_pos/domain"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// CustomerHandler  represent the httphandler for Customer
type MediaHandler struct {
	AUsecase domain.MediaUseCase
}

// NewCustomerHandler will initialize the Customers/ resources endpoint
func NewCustomerHandler(e *echo.Echo, us domain.MediaUseCase, cfg config.Config) {
	handler := &MediaHandler{
		AUsecase: us,
	}
	CustomerV1:= e.Group("")

	CustomerV1.Use(middlwr.JWTMiddleware(cfg))
	CustomerV1.POST("/media/upload", handler.UploadMedia)
}

// FetchCustomer will fetch the Customer based on given params
func (a *MediaHandler) UploadMedia(c echo.Context) error {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}
	if claims.RoleID != 1 {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrNotAuthorized.Error(),
			Data:    nil,
		}))
	}
	form, err := c.FormFile("file")
	if err != nil {
		fmt.Println(form,err)
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}
	filename, err := a.AUsecase.UploadMedia(form)

	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(filename))
}

