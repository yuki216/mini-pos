package usecase

import (
	"context"
	"errors"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type PaymentUsecase struct {
	OrderRepo    domain.OrderRepository
	PaymentRepo    domain.PaymentRepository
	contextTimeout time.Duration
}

// NewPaymentUsecase will create new an PaymentUsecase object representation of domain.PaymentUsecase interface
func NewPaymentUsecase(order domain.OrderRepository, a domain.PaymentRepository, timeout time.Duration) domain.PaymentUsecase {
	return &PaymentUsecase{
		OrderRepo:    order,
		PaymentRepo:    a,
		contextTimeout: timeout,
	}
}


func (a *PaymentUsecase) Payment(c context.Context, m domain.RequestPayment) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	//check payment ammount must be same
	id,err:= a.OrderRepo.GetStatus(ctx)
	if err != nil {
		err = errors.New("error: cannot checkout order empty")
		return
	}
	m.Order = id
	err = a.OrderRepo.Update(ctx, &domain.Order{
						ID:         m.Order,
						Customer:   domain.Customer{
							ID:        m.CustomerID,
						},
						IsCheckout: 1,
					})
	if err != nil {
		return
	}
	err = a.PaymentRepo.Payment(ctx, m)
	return
}

func (a *PaymentUsecase) GetByCustomerID(c context.Context, id int64) ([]domain.PaymentCustomer, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res,  err := a.PaymentRepo.PaymentByCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(res)<=0 {
		return nil, constants.ErrDataNotFound
	}
	return res,nil
}

func (a *PaymentUsecase) GetByProductID(c context.Context, id int64) ([]domain.PaymentProduct, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res,  err := a.PaymentRepo.PaymentByProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(res)<=0 {
		return nil, constants.ErrDataNotFound
	}
	return res,nil
}