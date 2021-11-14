package usecase

import (
	"context"
	"mini_pos/constants"
	"time"

	"mini_pos/domain"
)

type productUsecase struct {
	productRepo    domain.ProductRepository
	catRepo     domain.CategoryRepository
	contextTimeout time.Duration
}

// NewproductUsecase will create new an productUsecase object representation of domain.productUsecase interface
func NewproductUsecase(a domain.ProductRepository, cat domain.CategoryRepository, timeout time.Duration) domain.ProductUsecase {
	return &productUsecase{
		productRepo:    a,
		catRepo:     cat,
		contextTimeout: timeout,
	}
}


func (a *productUsecase) fillProductDetails(c context.Context, data []domain.Product) ([]domain.Product, error) {

	for index, product := range data {
		if product.Category.ID != 0 {
			cat, err := a.catRepo.GetByID(c, product.Category.ID)
			if err == nil {
				data[index].Category = domain.Category{
					ID:          cat.ID,
					Name:        cat.Name,
					Description: cat.Description,
					UpdatedAt:   cat.UpdatedAt,
					CreatedAt:   cat.CreatedAt,
				}
			}
		}
	}


	return data, nil
}

func (a *productUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Product, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.productRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillProductDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	if len(res)<=0 {
		nextCursor = ""
		err = constants.ErrDataNotFound
	}
	return
}

func (a *productUsecase) GetByID(c context.Context, id int64) (res domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.productRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resCategory, err := a.catRepo.GetByID(ctx, res.Category.ID)
	if err != nil {
		return domain.Product{}, err
	}
	res.Category = resCategory

	return
}

func (a *productUsecase) Update(c context.Context, ar *domain.RequestProduct) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.productRepo.Update(ctx, ar)
}

func (a *productUsecase) GetByName(c context.Context, name string) (res domain.Product, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.productRepo.GetByName(ctx, name)
	if err != nil {
		return
	}

	resCategory, err := a.catRepo.GetByID(ctx, res.Category.ID)
	if err != nil {
		return domain.Product{}, err
	}
	res.Category = resCategory

	return
}

func (a *productUsecase) Store(c context.Context, m *domain.RequestProduct) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedproduct, _ := a.GetByName(ctx, m.Name)
	if existedproduct != (domain.Product{}) {
		return domain.ErrConflict
	}
	err = a.productRepo.Store(ctx, m)
	return
}

func (a *productUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedproduct, err := a.productRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedproduct == (domain.Product{}) {
		return domain.ErrNotFound
	}
	return a.productRepo.Delete(ctx, id)
}
