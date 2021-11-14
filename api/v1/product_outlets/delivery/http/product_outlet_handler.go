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

// ProductOutletsHandler  represent the httphandler for ProductOutlets
type ProductOutletsHandler struct {
	AUsecase domain.ProductOutletsUsecase
}

// NewProductOutletsHandler will initialize the ProductOutletss/ resources endpoint
func NewProductOutletsHandler(e *echo.Echo, us domain.ProductOutletsUsecase, cfg config.Config) {
	handler := &ProductOutletsHandler{
		AUsecase: us,
	}
	ProductOutletsV1:= e.Group("")

	ProductOutletsV1.Use(middlwr.JWTMiddleware(cfg))
	ProductOutletsV1.GET("/product_outlets", handler.FetchProductOutlets)
	ProductOutletsV1.POST("/product_outlets", handler.Store)
	ProductOutletsV1.PUT("/product_outlets", handler.Update)
	ProductOutletsV1.PUT("/product_outlets/stock", handler.UpdateStock)
	ProductOutletsV1.GET("/product_outlets/:id", handler.GetByID)
	ProductOutletsV1.DELETE("/product_outlets/:id", handler.Delete)
}

// FetchProductOutlets will fetch the ProductOutlets based on given params
func (a *ProductOutletsHandler) FetchProductOutlets(c echo.Context) error {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}
	outletS := c.QueryParam("outlet")
	outlet, err := strconv.Atoi(outletS)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: constants.ErrParamsIsNotInvalid.Error(),
			Data:    nil,
		}))
	}
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num), int64(outlet))
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

// GetByID will get ProductOutlets by given id
func (a *ProductOutletsHandler) GetByID(c echo.Context) error {
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

	outletS := c.QueryParam("outlet")
	outlet, err := strconv.Atoi(outletS)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	prod, err := a.AUsecase.GetByID(ctx, id, int64(outlet))
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

func (a *ProductOutletsHandler) GetBySKU(c echo.Context) error {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}

	sku := c.Param("sku")

	ctx := c.Request().Context()

	prod, err := a.AUsecase.GetBySku(ctx, sku)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(prod))
}

func isRequestValid(m *domain.RequestProductOutlets) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update will the ProductOutlets by given request body
func (a *ProductOutletsHandler) Update(c echo.Context) (err error) {
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
	var ProductOutlets domain.RequestProductOutlets
	err = c.Bind(&ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&ProductOutlets); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	ProductOutlets.OutletID = int64(claims.RoleID)
	err = a.AUsecase.Update(ctx, &ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(ProductOutlets))
}

// Update will the ProductOutlets by given request body
func (a *ProductOutletsHandler) UpdateStock(c echo.Context) (err error) {
	claims := middlwr.GetTokenFromContext(c)
	if claims == nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: constants.ErrTokenAlreadyExpired.Error(),
			Data:    nil,
		}))
	}
	var ProductOutlets domain.ProductOutletsStock
	err = c.Bind(&ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	validate := validator.New()
	err = validate.Struct(ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	err = a.AUsecase.UpdateStock(ctx, ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(ProductOutlets))
}

// Store will store the ProductOutlets by given request body
func (a *ProductOutletsHandler) Store(c echo.Context) (err error) {
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
	var ProductOutlets domain.RequestProductOutlets
	err = c.Bind(&ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	var ok bool
	if ok, err = isRequestValid(&ProductOutlets); !ok {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()
	err = a.AUsecase.Store(ctx, &ProductOutlets)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	return c.JSON(common.NewSuccessResponse(ProductOutlets))
}

// Delete will delete ProductOutlets by given param
func (a *ProductOutletsHandler) Delete(c echo.Context) error {
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

	err = a.AUsecase.Delete(ctx, id, int64(claims.RoleID))
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
