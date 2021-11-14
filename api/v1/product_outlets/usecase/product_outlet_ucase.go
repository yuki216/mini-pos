package usecase

import (
	"context"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type ProductOutletsUsecase struct {
	ProductOutletsRepo    domain.ProductOutletsRepository
	productRepo    domain.ProductRepository
	outletRepo     domain.OutletRepository
	supRepo domain.SupplierRepository
	contextTimeout time.Duration
}

// NewProductOutletsUsecase will create new an ProductOutletsUsecase object representation of domain.ProductOutletsUsecase interface
func NewProductOutletsUsecase(a domain.ProductOutletsRepository, prod domain.ProductRepository,out domain.OutletRepository, sup domain.SupplierRepository, timeout time.Duration) domain.ProductOutletsUsecase {
	return &ProductOutletsUsecase{
		ProductOutletsRepo:    a,
		productRepo:     prod,
		outletRepo:     out,
		supRepo: sup,
		contextTimeout: timeout,
	}
}


func (a *ProductOutletsUsecase) fillProductOutletsDetails(c context.Context, data []domain.ProductOutlets) ([]domain.ProductOutlets, error) {

	for index, ProductOutlets := range data {
		if ProductOutlets.Product.ID != 0 {
			prod, err := a.productRepo.GetByID(c, ProductOutlets.Product.ID)
			if err == nil {
				data[index].Product = prod
			}
		}

		if ProductOutlets.Outlet.ID != 0 {
			out, err := a.outletRepo.GetByID(c, ProductOutlets.Outlet.ID)
			if err == nil {
				data[index].Outlet = out
			}

		}

		if ProductOutlets.Supplier.ID != 0 {
			sup, err := a.supRepo.GetByID(c, ProductOutlets.Supplier.ID)
			if err == nil {
				data[index].Supplier = sup
			}

		}
	}


	return data, nil
}

func (a *ProductOutletsUsecase) Fetch(c context.Context, cursor string, num int64, outlet int64) (res []domain.ProductOutlets, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.ProductOutletsRepo.Fetch(ctx, cursor, num,  outlet)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillProductOutletsDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	if len(res)<=0 {
		nextCursor = ""
		err = constants.ErrDataNotFound
	}
	return
}

func (a *ProductOutletsUsecase) GetByID(c context.Context, id int64, outlet int64) (res domain.ProductOutlets, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.ProductOutletsRepo.GetByID(ctx, id, outlet)
	if err != nil {
		return
	}

	resProduct, err := a.productRepo.GetByID(ctx, res.Product.ID)
	if err != nil {
		return domain.ProductOutlets{}, err
	}
	res.Product = resProduct

	resOutlet, err := a.outletRepo.GetByID(ctx, res.Outlet.ID)
	if err != nil {
		return domain.ProductOutlets{}, err
	}
	res.Outlet = resOutlet

	resSupplier, err := a.supRepo.GetByID(ctx, res.Supplier.ID)
	if err != nil {
		return domain.ProductOutlets{}, err
	}
	res.Supplier = resSupplier
	return
}

func (a *ProductOutletsUsecase) Update(c context.Context, ar *domain.RequestProductOutlets) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.ProductOutletsRepo.Update(ctx, ar)
}

func (a *ProductOutletsUsecase) UpdateStock(c context.Context,  stock domain.ProductOutletsStock) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.ProductOutletsRepo.UpdateStock(ctx, stock)
}

func (a *ProductOutletsUsecase) GetBySku(c context.Context, sku string) (res domain.ProductOutlets, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.ProductOutletsRepo.GetBySku(ctx, sku)
	if err != nil {
		return
	}

	resProduct, err := a.productRepo.GetByID(ctx, res.Product.ID)
	if err != nil {
		return domain.ProductOutlets{}, err
	}
	res.Product = resProduct

	resOutlet, err := a.outletRepo.GetByID(ctx, res.Outlet.ID)
	if err != nil {
		return domain.ProductOutlets{}, err
	}
	res.Outlet = resOutlet
	return
}

func (a *ProductOutletsUsecase) Store(c context.Context, m *domain.RequestProductOutlets) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedProductOutlets, _ := a.GetBySku(ctx, m.Sku)
	if existedProductOutlets != (domain.ProductOutlets{}) {
		return domain.ErrConflict
	}
	err = a.ProductOutletsRepo.Store(ctx, m)
	return
}

func (a *ProductOutletsUsecase) Delete(c context.Context, id int64, outlet int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedProductOutlets, err := a.ProductOutletsRepo.GetByID(ctx, id, outlet)
	if err != nil {
		return
	}
	if existedProductOutlets == (domain.ProductOutlets{}) {
		return domain.ErrNotFound
	}
	return a.ProductOutletsRepo.Delete(ctx, id)
}
