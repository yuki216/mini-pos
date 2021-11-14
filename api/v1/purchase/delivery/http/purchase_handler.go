package http

import (
	"fmt"
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

// PurchaseHandler  represent the httphandler for Purchase
type PurchaseHandler struct {
	AUsecase domain.PurchaseUsecase
}

// NewPurchaseHandler will initialize the Purchases/ resources endpoint
func NewPurchaseHandler(e *echo.Echo, us domain.PurchaseUsecase, cfg config.Config) {
	handler := &PurchaseHandler{
		AUsecase: us,
	}
	PurchaseV1:= e.Group("")

	PurchaseV1.Use(middlwr.JWTMiddleware(cfg))
	PurchaseV1.GET("/purchases", handler.FetchPurchase)
	PurchaseV1.GET("/purchases/supplier", handler.GetBySuppliers)
	PurchaseV1.GET("/purchases/supplier/:id", handler.GetBySupplierID)
	PurchaseV1.GET("/purchases/product", handler.GetByProducts)
	PurchaseV1.GET("/purchases/product/:id", handler.GetByProductID)
	PurchaseV1.POST("/purchases", handler.Store)
	PurchaseV1.PUT("/purchases", handler.Update)
	PurchaseV1.GET("/purchases/:id", handler.GetByID)
	PurchaseV1.DELETE("/purchases/:id", handler.Delete)
}

// FetchPurchase will fetch the Purchase based on given params
func (a *PurchaseHandler) FetchPurchase(c echo.Context) error {
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

// GetByID will get Purchase by given id
func (a *PurchaseHandler) GetByID(c echo.Context) error {
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

// GetByID will get Purchase by given id
func (a *PurchaseHandler) GetBySupplierID(c echo.Context) error {
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

	prod, err := a.AUsecase.GetBySupplierID(ctx, id)
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
func (a *PurchaseHandler) GetBySuppliers(c echo.Context) error {
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
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetBySupplierID(ctx, 0)
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
func (a *PurchaseHandler) GetByProductID(c echo.Context) error {
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
func (a *PurchaseHandler) GetByProducts(c echo.Context) error {
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
	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetByProductID(ctx, 0)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

func isRequestValid(m *domain.RequestPurchase) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update will the Purchase by given request body
func (a *PurchaseHandler) Update(c echo.Context) (err error) {
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
	var Purchase domain.RequestPurchase
	err = c.Bind(&Purchase)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&Purchase); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	Purchase.OutletID = int64(claims.RoleID)
	err = a.AUsecase.Update(ctx, &Purchase)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(Purchase))
}

// Store will store the Purchase by given request body
func (a *PurchaseHandler) Store(c echo.Context) (err error) {
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
	var Purchase domain.RequestPurchase
	err = c.Bind(&Purchase)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&Purchase); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &Purchase)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(Purchase))
}

// Delete will delete Purchase by given param
func (a *PurchaseHandler) Delete(c echo.Context) error {
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
