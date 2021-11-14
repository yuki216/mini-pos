package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"mini_pos/api/common"
	middlwr "mini_pos/api/middleware"
	"mini_pos/config"
	"mini_pos/constants"
	"net/http"
	"strconv"

	"mini_pos/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// PaymentHandler  represent the httphandler for Payment
type PaymentHandler struct {
	AUsecase domain.PaymentUsecase
}

// NewPaymentHandler will initialize the Payments/ resources endpoint
func NewPaymentHandler(e *echo.Echo, us domain.PaymentUsecase, cfg config.Config) {
	handler := &PaymentHandler{
		AUsecase: us,
	}
	PaymentV1:= e.Group("")

	PaymentV1.Use(middlwr.JWTMiddleware(cfg))
	PaymentV1.POST("/payment", handler.Payment)
	PaymentV1.GET("/payment/customer", handler.GetByCustomers)
	PaymentV1.GET("/payment/customer/:id", handler.GetByCustomerID)
	PaymentV1.GET("/payment/product", handler.GetByProducts)
	PaymentV1.GET("/payment/product/:id", handler.GetByProductID)
}

func isRequestValid(m *domain.RequestPayment) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the Payment by given request body
func (a *PaymentHandler) Payment(c echo.Context) (err error) {
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
	var Payment domain.RequestPayment
	err = c.Bind(&Payment)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&Payment); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Payment(ctx, Payment)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(Payment))
}

// GetByID will get Purchase by given id
func (a *PaymentHandler) GetByCustomerID(c echo.Context) error {
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
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}

	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		idP=0
	}
	fmt.Println(idP)
	id := int64(idP)
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByCustomerID(ctx, id)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

// GetByID will get Purchase by given id
func (a *PaymentHandler) GetByCustomers(c echo.Context) error {
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
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}

	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		idP=0
	}
	fmt.Println(idP)
	id := int64(idP)
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByCustomerID(ctx, id)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

// GetByID will get Purchase by given id
func (a *PaymentHandler) GetByProductID(c echo.Context) error {
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
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}

	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		idP=0
	}
	fmt.Println(idP)
	id := int64(idP)
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByProductID(ctx, id)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

// GetByID will get Purchase by given id
func (a *PaymentHandler) GetByProducts(c echo.Context) error {
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
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}

	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		idP=0
	}
	fmt.Println(idP)
	id := int64(idP)
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByProductID(ctx, id)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
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
