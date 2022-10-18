package usecases

import (
	"aopt-db-sync/models"
	"context"
)

func NewAoptUseCase(brandRepo brandRepo, vehicleRepo vehicleRepo, modificationRepo modificationRepo) *AoptUseCase {
	return &AoptUseCase{
		brandRepo:        brandRepo,
		vehicleRepo:      vehicleRepo,
		modificationRepo: modificationRepo,
	}
}

type brandRepo interface {
	Get(ctx context.Context, in *models.BrandFilter) ([]*models.Brand, error)
	Upsert(ctx context.Context, in *models.Brand) (*models.Brand, error)
}
type vehicleRepo interface {
	GetAll(ctx context.Context) ([]*models.Vehicle, error)
}
type modificationRepo interface {
	GetAll(ctx context.Context) ([]*models.Modification, error)
}

type AoptUseCase struct {
	brandRepo        brandRepo
	vehicleRepo      vehicleRepo
	modificationRepo modificationRepo
}

func (useCase AoptUseCase) GetBrands(ctx context.Context) ([]*models.Brand, error) {
	return useCase.brandRepo.Get(ctx, nil)
}

func (useCase AoptUseCase) GetBrandByID(ctx context.Context, brandID string) ([]*models.Brand, error) {
	return useCase.brandRepo.Get(ctx, &models.BrandFilter{ID: &brandID})
}

func (useCase AoptUseCase) GetBrandByName(ctx context.Context, name string) (*models.Brand, error) {
	result, err := useCase.brandRepo.Get(ctx, &models.BrandFilter{Name: &name})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, models.ErrNoObject
	}

	return result[0], nil
}

func (useCase AoptUseCase) UpsertBrand(ctx context.Context, brand *models.Brand) (*models.Brand, error) {
	return useCase.brandRepo.Upsert(ctx, brand)
}

func (useCase AoptUseCase) GetVehicles(ctx context.Context) ([]*models.Vehicle, error) {
	return useCase.vehicleRepo.GetAll(ctx)
}

func (useCase AoptUseCase) GetModifications(ctx context.Context) ([]*models.Modification, error) {
	return useCase.modificationRepo.GetAll(ctx)
}
