package usecase

import (
	"context"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type CustomerUsecase struct {
	CustomerRepo    domain.CustomerRepository
	contextTimeout time.Duration
}

// NewCustomerUsecase will create new an CustomerUsecase object representation of domain.CustomerUsecase interface
func NewCustomerUsecase(a domain.CustomerRepository, timeout time.Duration) domain.CustomerUsecase {
	return &CustomerUsecase{
		CustomerRepo:    a,
		contextTimeout: timeout,
	}
}


func (a *CustomerUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Customer, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.CustomerRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res)<=0 {
		nextCursor = ""
		err = constants.ErrDataNotFound
	}
	return
}

func (a *CustomerUsecase) GetByID(c context.Context, id int64) (res domain.Customer, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.CustomerRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (a *CustomerUsecase) Update(c context.Context, ar *domain.Customer) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.CustomerRepo.Update(ctx, ar)
}

func (a *CustomerUsecase) GetByName(c context.Context, name string) (res domain.Customer, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.CustomerRepo.GetByName(ctx, name)
	if err != nil {
		return
	}

	return
}

func (a *CustomerUsecase) Store(c context.Context, m *domain.Customer) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedCustomer, _ := a.GetByName(ctx, m.Name)
	if existedCustomer != (domain.Customer{}) {
		return domain.ErrConflict
	}
	err = a.CustomerRepo.Store(ctx, m)
	return
}

func (a *CustomerUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedCustomer, err := a.CustomerRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedCustomer == (domain.Customer{}) {
		return domain.ErrNotFound
	}
	return a.CustomerRepo.Delete(ctx, id)
}
