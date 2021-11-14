package http

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"mini_pos/api/common"
	"mini_pos/api/middleware"
	"mini_pos/config"
	"mini_pos/constants"
	"mini_pos/domain"
	"mini_pos/utils"
	"net/http"
	"time"
)

type AuthHandler struct {
	AuUseCase domain.AuthUseCase
	Cfg config.Config
}

func NewAuthHandler(e *echo.Echo, us domain.AuthUseCase, cfg config.Config)  {
	handler := &AuthHandler{
		AuUseCase: us,
		Cfg: cfg,
	}

	e.POST("/auth/login", handler.Login)
	e.POST("/auth/register", handler.Register)
}
func (a *AuthHandler) Login(c echo.Context) (err error) {
	var auth domain.UserLogin
	err = c.Bind(&auth)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	//validasi email
	status := utils.ValidateEmail(auth.Email)
	if !status{
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: constants.ErrEmailIsNotValid.Error(),
			Data:    nil,
		}))
	}

	// validasi password
	status=utils.ValidatePassword(auth.Password)
	if !status{ fmt.Println(status)
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: constants.ErrInvalidPassword.Error(),
			Data:    nil,
		}))
	}

	ctx := c.Request().Context()

	user, err := a.AuUseCase.GetEmailUser(ctx, auth.Email)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}

	//fmt.Println(auth.Password, user.Password)
	if !utils.CheckPasswordHash(auth.Password, user.Password){
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: "error: invalid Password user!",
			Data:    nil,
		}))
	}

	if auth.Username != user.UserName{
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: "error: Username Not Match!",
			Data:    nil,
		}))
	}

	claims := &middleware.JwtCustomClaims{
		Name:     user.Name,
		ID:       user.ID,
		RoleID: user.RoleID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix() ,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	generate, err := token.SignedString([]byte(a.Cfg.JWTConfig.Secret))
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	res :=  domain.UserResponse{
			ID:    user.ID,
			Token: generate,
		}

	return c.JSON(common.NewSuccessResponse(res))
}

func (a *AuthHandler) Register(c echo.Context) (err error) {
	var reg domain.RegisterUser
	err = c.Bind(&reg)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Data:    nil,
		}))
	}

	//validasi email
	status := utils.ValidateEmail(reg.Email)
	if !status{ fmt.Println(status)
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Email not Valid!",
			Data:    nil,
		}))
	}

	// validasi password
	status=utils.ValidatePassword(reg.Password)
	if !status{ fmt.Println(status)
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: constants.ErrInvalidPassword.Error(),
			Data:    nil,
		}))
	}
	ctx := c.Request().Context()

	_, err = a.AuUseCase.Register(ctx, reg)
	if err != nil {
		return c.JSON(common.NewErrorResponse(common.ControllerResponse{
			Code:    getStatusCode(err),
			Message: err.Error(),
			Data:    nil,
		}))
	}
	return c.JSON(common.NewSuccessResponse(reg))
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