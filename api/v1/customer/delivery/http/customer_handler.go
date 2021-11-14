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

// CustomerHandler  represent the httphandler for Customer
type CustomerHandler struct {
	AUsecase domain.CustomerUsecase
}

// NewCustomerHandler will initialize the Customers/ resources endpoint
func NewCustomerHandler(e *echo.Echo, us domain.CustomerUsecase, cfg config.Config) {
	handler := &CustomerHandler{
		AUsecase: us,
	}
	CustomerV1:= e.Group("")

	CustomerV1.Use(middlwr.JWTMiddleware(cfg))
	CustomerV1.GET("/customers", handler.FetchCustomer)
	CustomerV1.POST("/customers", handler.Store)
	CustomerV1.PUT("/customers", handler.Update)
	CustomerV1.GET("/customers/:id", handler.GetByID)
	CustomerV1.DELETE("/customers/:id", handler.Delete)
}

// FetchCustomer will fetch the Customer based on given params
func (a *CustomerHandler) FetchCustomer(c echo.Context) error {
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

// GetByID will get Customer by given id
func (a *CustomerHandler) GetByID(c echo.Context) error {
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

	cus, err := a.AUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(cus))
}

func (a *CustomerHandler) GetByName(c echo.Context) error {
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

	cus, err := a.AUsecase.GetByName(ctx, sku)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(cus))
}

func isRequestValid(m *domain.Customer) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update will the Customer by given request body
func (a *CustomerHandler) Update(c echo.Context) (err error) {
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
	var Customer domain.Customer
	err = c.Bind(&Customer)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&Customer); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Update(ctx, &Customer)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(Customer))
}


// Store will store the Customer by given request body
func (a *CustomerHandler) Store(c echo.Context) (err error) {
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
	var Customer domain.Customer
	err = c.Bind(&Customer)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&Customer); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &Customer)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(Customer))
}

// Delete will delete Customer by given param
func (a *CustomerHandler) Delete(c echo.Context) error {
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
