package http

import (
	"mini_pos/api/common"
	middlwr "mini_pos/api/middleware"
	"mini_pos/config"
	"mini_pos/constants"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"mini_pos/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// SupplierHandler  represent the httphandler for Supplier
type SupplierHandler struct {
	AUsecase domain.SupplierUsecase
}

// NewSupplierHandler will initialize the Suppliers/ resources endpoint
func NewSupplierHandler(e *echo.Echo, us domain.SupplierUsecase, cfg config.Config) {
	handler := &SupplierHandler{
		AUsecase: us,
	}
	SupplierV1:= e.Group("")

	SupplierV1.Use(middlwr.JWTMiddleware(cfg))
	SupplierV1.GET("/suppliers", handler.FetchSupplier)
	SupplierV1.POST("/suppliers", handler.Store)
	SupplierV1.PUT("/suppliers", handler.Update)
	SupplierV1.GET("/suppliers/:id", handler.GetByID)
	SupplierV1.DELETE("/suppliers/:id", handler.Delete)
}

// FetchSupplier will fetch the Supplier based on given params
func (a *SupplierHandler) FetchSupplier(c echo.Context) error {
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
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(common.NewSuccessResponse(listAr))
}

// GetByID will get Supplier by given id
func (a *SupplierHandler) GetByID(c echo.Context) error {
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
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusNotFound,
			Message: constants.ErrDataNotFound.Error(),
			Data:    nil,
		}))
	}

	id := int64(idP)
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

func (a *SupplierHandler) GetByName(c echo.Context) error {
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
	sku := c.Param("name")

	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByName(ctx, sku)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

func isRequestValid(m *domain.Supplier) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update will the Supplier by given request body
func (a *SupplierHandler) Update(c echo.Context) (err error) {
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
	var Supplier domain.Supplier
	err = c.Bind(&Supplier)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&Supplier); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Update(ctx, &Supplier)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, Supplier)
}


// Store will store the Supplier by given request body
func (a *SupplierHandler) Store(c echo.Context) (err error) {
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
	var Supplier domain.Supplier
	err = c.Bind(&Supplier)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&Supplier); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &Supplier)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, Supplier)
}

// Delete will delete Supplier by given param
func (a *SupplierHandler) Delete(c echo.Context) error {
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
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusNotFound,
			Message: constants.ErrDataNotFound.Error(),
			Data:    nil,
		}))
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = a.AUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:   getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponseWithoutData())
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
