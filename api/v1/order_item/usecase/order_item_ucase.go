package usecase

import (
	"context"
	"fmt"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type OrderItemUsecase struct {
	OrderItemRepo    domain.OrderItemRepository
	OrderRepo     domain.OrderRepository
	ProductRepo     domain.ProductRepository
	contextTimeout time.Duration
}

// NewOrderItemUsecase will create new an OrderItemUsecase object representation of domain.OrderItemUsecase interface
func NewOrderItemUsecase(a domain.OrderItemRepository, order domain.OrderRepository, prod domain.ProductRepository, timeout time.Duration) domain.OrderItemUsecase {
	return &OrderItemUsecase{
		OrderItemRepo:    a,
		OrderRepo:     order,
		ProductRepo: prod,
		contextTimeout: timeout,
	}
}


func (a *OrderItemUsecase) fillOrderItemDetails(c context.Context, data []domain.OrderItem) ([]domain.OrderItem, error) {

	for index, OrderItem := range data {
		if OrderItem.Order.ID != 0 {
			order, err := a.OrderRepo.GetByID(c, OrderItem.Order.ID)
			if err == nil {
				data[index].Order = order
			}
		}

		if OrderItem.Product.ID != 0 {
			prod, err := a.ProductRepo.GetByID(c, OrderItem.Product.ID)
			if err == nil {
				data[index].Product = prod
			}
		}
	}


	return data, nil
}

func (a *OrderItemUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.OrderItem, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.OrderItemRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillOrderItemDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	if len(res)<=0 {
		nextCursor = ""
		err = constants.ErrDataNotFound
	}
	return
}

func (a *OrderItemUsecase) FetchByOrderID(c context.Context, id int64) (res []domain.OrderItem, err error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.OrderItemRepo.FetchByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}

	res, err = a.fillOrderItemDetails(ctx, res)
	if err != nil {
		return nil, err
	}

	return
}

func (a *OrderItemUsecase) GetByID(c context.Context, id int64) (res domain.OrderItem, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.OrderItemRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resOrder, err := a.OrderRepo.GetByID(ctx, res.Order.ID)
	if err != nil {
		return domain.OrderItem{}, err
	}
	res.Order = resOrder

	resProduct, err := a.ProductRepo.GetByID(ctx, res.Product.ID)
	if err != nil {
		return domain.OrderItem{}, err
	}
	res.Product = resProduct

	return
}

func (a *OrderItemUsecase) Update(c context.Context, ar *domain.RequestOrderItem) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.OrderItemRepo.Update(ctx, ar)
}


func (a *OrderItemUsecase) Store(c context.Context, m *domain.RequestOrderItem) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	//check order id
	id,err:= a.OrderRepo.GetStatus(ctx)

	if err != nil {
		id,err = a.OrderRepo.Store(ctx,&domain.Order{
			IsCheckout: 0,
		})
	}

	existedOrderItem, _ := a.OrderItemRepo.GetByOrderIDAndProdID(ctx, id, m.ProductID)
	if existedOrderItem != (domain.OrderItem{}) {
		m.ID = existedOrderItem.ID
		m.Qty += existedOrderItem.Qty
		a.OrderItemRepo.Update(ctx,m)
		return
	}
	fmt.Println(id,err)
	m.OrderID = id
	err = a.OrderItemRepo.Store(ctx, m)
	return
}

func (a *OrderItemUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedOrderItem, err := a.OrderItemRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedOrderItem == (domain.OrderItem{}) {
		return domain.ErrNotFound
	}
	return a.OrderItemRepo.Delete(ctx, id)
}
