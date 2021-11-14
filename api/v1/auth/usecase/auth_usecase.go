package usecase

import (
	"context"
	"fmt"
	"mini_pos/constants"
	"mini_pos/domain"
	"time"
)

type authUsecase struct {
	authRepo domain.AuthRepository
	contextTimeout time.Duration
}

func NewAuthUseCase ( auth domain.AuthRepository, timeout time.Duration) domain.AuthUseCase {
	return &authUsecase{
	 authRepo: auth,
	 contextTimeout: timeout,
	}
}

func (a *authUsecase) GetEmailUser(c context.Context, username string)(res domain.User, err error){
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.authRepo.GetEmailUser(ctx,username)
	fmt.Println(res,err)
	if err != nil {
		fmt.Println(err)
		err = constants.ErrEmailNotFound
		return
	}
	return res, nil
}

func (a *authUsecase) Register(c context.Context, register domain.RegisterUser)(res domain.User, err error){
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	id, err := a.authRepo.Register(ctx, register)
	if err != nil {
		return
	}
	res.ID=       id
	res.Name=     register.Name
	res.UserName= register.UserName
	res.Email=    register.Email
	return res, nil
}
