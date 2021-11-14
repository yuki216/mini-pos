package usecase

import (
	"context"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type PurchaseUsecase struct {
	ProductRepo    domain.ProductRepository
	PurchaseRepo    domain.PurchaseRepository
	catRepo     domain.CategoryRepository
	outletRepo     domain.OutletRepository
	supRepo domain.SupplierRepository
	contextTimeout time.Duration
}

// NewPurchaseUsecase will create new an PurchaseUsecase object representation of domain.PurchaseUsecase interface
func NewPurchaseUsecase(prod domain.ProductRepository, a domain.PurchaseRepository, cat domain.CategoryRepository,out domain.OutletRepository, sup domain.SupplierRepository, timeout time.Duration) domain.PurchaseUsecase {
	return &PurchaseUsecase{
		ProductRepo:    prod,
		PurchaseRepo:    a,
		catRepo:     cat,		
		outletRepo:     out,
		supRepo: sup,
		contextTimeout: timeout,
	}
}


func (a *PurchaseUsecase) fillPurchaseDetails(c context.Context, data []domain.Purchase) ([]domain.Purchase, error) {

	for index, Purchase := range data {
		if Purchase.Category.ID != 0 {
			prod, err := a.ProductRepo.GetByID(c, Purchase.Product.ID)
			if err == nil {
				data[index].Product = prod
			}
		}
		if Purchase.Category.ID != 0 {
			cat, err := a.catRepo.GetByID(c, Purchase.Outlet.ID)
			if err == nil {
				data[index].Category = cat
			}
		}

		if Purchase.Outlet.ID != 0 {
			out, err := a.outletRepo.GetByID(c, Purchase.Outlet.ID)
			if err == nil {
				data[index].Outlet = out
			}

		}

		if Purchase.Supplier.ID != 0 {
			sup, err := a.supRepo.GetByID(c, Purchase.Supplier.ID)
			if err == nil {
				data[index].Supplier = sup
			}

		}
	}


	return data, nil
}


func (a *PurchaseUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Purchase, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.PurchaseRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillPurchaseDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	if len(res)<=0 {
		nextCursor = ""
		err = constants.ErrDataNotFound
	}
	return
}

func (a *PurchaseUsecase) GetByID(c context.Context, id int64) (res domain.Purchase, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.PurchaseRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resProduct, err := a.ProductRepo.GetByID(ctx,res.Product.ID)
	if err != nil {
		return domain.Purchase{}, err
	}
	res.Product = resProduct

	resCategory, err := a.catRepo.GetByID(ctx, res.Category.ID)
	if err != nil {
		return domain.Purchase{}, err
	}
	res.Category = resCategory

	resOutlet, err := a.outletRepo.GetByID(ctx, res.Outlet.ID)
	if err != nil {
		return domain.Purchase{}, err
	}
	res.Outlet = resOutlet

	resSupplier, err := a.supRepo.GetByID(ctx, res.Supplier.ID)
	if err != nil {
		return domain.Purchase{}, err
	}
	res.Supplier = resSupplier
	return
}
func (a *PurchaseUsecase) GetBySupplierID(c context.Context, id int64) ([]domain.PurchaseBySupplier, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res,  err := a.PurchaseRepo.GetBySupplierID(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(res)<=0 {
		return nil, constants.ErrDataNotFound
	}
	return res,nil
}

func (a *PurchaseUsecase) GetByProductID(c context.Context, id int64) ([]domain.PurchaseByProduct, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res,  err := a.PurchaseRepo.GetByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(res)<=0 {
		return nil, constants.ErrDataNotFound
	}
	return res,nil
}

func (a *PurchaseUsecase) Update(c context.Context, ar *domain.RequestPurchase) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.PurchaseRepo.Update(ctx, ar)
}

func (a *PurchaseUsecase) Store(c context.Context, m *domain.RequestPurchase) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedPurchase, _ := a.GetByID(ctx, m.ID)
	if existedPurchase != (domain.Purchase{}) {
		return domain.ErrConflict
	}
	err = a.PurchaseRepo.Store(ctx, m)
	return
}

func (a *PurchaseUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedPurchase, err := a.PurchaseRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedPurchase == (domain.Purchase{}) {
		return domain.ErrNotFound
	}
	return a.PurchaseRepo.Delete(ctx, id)
}
