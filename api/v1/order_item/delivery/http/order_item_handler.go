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

// OrderItemHandler  represent the httphandler for OrderItem
type OrderItemHandler struct {
	AUsecase domain.OrderItemUsecase
}

// NewOrderItemHandler will initialize the OrderItems/ resources endpoint
func NewOrderItemHandler(e *echo.Echo, us domain.OrderItemUsecase, cfg config.Config) {
	handler := &OrderItemHandler{
		AUsecase: us,
	}
	OrderItemV1:= e.Group("")

	OrderItemV1.Use(middlwr.JWTMiddleware(cfg))
	OrderItemV1.GET("/cart", handler.FetchOrderItem)
	OrderItemV1.GET("/cart/order/:id", handler.FetchByOrderID)
	OrderItemV1.POST("/cart", handler.Store)
	OrderItemV1.PUT("/cart", handler.Update)
	OrderItemV1.GET("/cart/:id", handler.GetByID)
	OrderItemV1.DELETE("/cart/:id", handler.Delete)
}

// FetchOrderItem will fetch the OrderItem based on given params
func (a *OrderItemHandler) FetchByOrderID(c echo.Context) error {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
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
	ctx := c.Request().Context()

	listAr, err := a.AUsecase.FetchByOrderID(ctx, int64(idP))
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(listAr))
}

func (a *OrderItemHandler) FetchOrderItem(c echo.Context) error {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
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

// GetByID will get OrderItem by given id
func (a *OrderItemHandler) GetByID(c echo.Context) error {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
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
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

func isRequestValid(m *domain.RequestOrderItem) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update will the OrderItem by given request body
func (a *OrderItemHandler) Update(c echo.Context) (err error) {
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
	var OrderItem domain.RequestOrderItem
	err = c.Bind(&OrderItem)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&OrderItem); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Update(ctx, &OrderItem)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(OrderItem))
}


// Store will store the OrderItem by given request body
func (a *OrderItemHandler) Store(c echo.Context) (err error) {
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
	var OrderItem domain.RequestOrderItem
	err = c.Bind(&OrderItem)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&OrderItem); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &OrderItem)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(OrderItem))
}

// Delete will delete OrderItem by given param
func (a *OrderItemHandler) Delete(c echo.Context) error {
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
			Code:    getStatusCode(err),
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
