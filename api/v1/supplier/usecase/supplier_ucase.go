package usecase

import (
	"context"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type SupplierUsecase struct {
	SupplierRepo    domain.SupplierRepository
	contextTimeout time.Duration
}

// NewSupplierUsecase will create new an SupplierUsecase object representation of domain.SupplierUsecase interface
func NewSupplierUsecase(a domain.SupplierRepository, timeout time.Duration) domain.SupplierUsecase {
	return &SupplierUsecase{
		SupplierRepo:    a,
		contextTimeout: timeout,
	}
}


func (a *SupplierUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Supplier, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.SupplierRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res)<=0 {
		nextCursor = ""
		err = constants.ErrDataNotFound
	}
	return
}

func (a *SupplierUsecase) GetByID(c context.Context, id int64) (res domain.Supplier, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.SupplierRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (a *SupplierUsecase) Update(c context.Context, ar *domain.Supplier) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.SupplierRepo.Update(ctx, ar)
}

func (a *SupplierUsecase) GetByName(c context.Context, name string) (res domain.Supplier, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.SupplierRepo.GetByName(ctx, name)
	if err != nil {
		return
	}

	return
}

func (a *SupplierUsecase) Store(c context.Context, m *domain.Supplier) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedSupplier, _ := a.GetByName(ctx, m.Name)
	if existedSupplier != (domain.Supplier{}) {
		return domain.ErrConflict
	}
	err = a.SupplierRepo.Store(ctx, m)
	return
}

func (a *SupplierUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedSupplier, err := a.SupplierRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedSupplier == (domain.Supplier{}) {
		return domain.ErrNotFound
	}
	return a.SupplierRepo.Delete(ctx, id)
}
