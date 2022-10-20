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
	Get(ctx context.Context, in *models.VehicleFilter) ([]*models.Vehicle, error)
	Upsert(ctx context.Context, in *models.Vehicle) (*models.Vehicle, error)
}
type modificationRepo interface {
	Get(ctx context.Context, in *models.ModificationFilter) ([]*models.Modification, error)
	Upsert(ctx context.Context, in *models.Vehicle) (*models.Vehicle, error)
}

type AoptUseCase struct {
	brandRepo        brandRepo
	vehicleRepo      vehicleRepo
	modificationRepo modificationRepo
}

func (useCase AoptUseCase) GetBrands(ctx context.Context) ([]*models.Brand, error) {
	return useCase.brandRepo.Get(ctx, nil)
}

func (useCase AoptUseCase) GetBrandByID(ctx context.Context, brandID int64) (*models.Brand, error) {
	brands, err := useCase.brandRepo.Get(ctx, &models.BrandFilter{ID: &brandID})
	if err != nil {
		return nil, err
	}
	if len(brands) == 0 {
		return nil, nil
	}

	return brands[0], nil
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
	return useCase.vehicleRepo.Get(ctx, nil)
}

func (useCase AoptUseCase) GetVehicleByID(ctx context.Context, vehicleID int64) (*models.Vehicle, error) {
	vehicles, err := useCase.vehicleRepo.Get(ctx, &models.VehicleFilter{ID: &vehicleID})
	if err != nil {
		return nil, err
	}
	if len(vehicles) == 0 {
		return nil, nil
	}

	return vehicles[0], nil
}

func (useCase AoptUseCase) GetVehiclesByBrandAndName(ctx context.Context, brandID int64, name string) (*models.Vehicle, error) {
	result, err := useCase.vehicleRepo.Get(ctx, &models.VehicleFilter{BrandID: &brandID, Name: &name})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, models.ErrNoObject
	}

	return result[0], nil
}

func (useCase AoptUseCase) GetVehiclesByName(ctx context.Context, name string) ([]*models.Vehicle, error) {
	return useCase.vehicleRepo.Get(ctx, &models.VehicleFilter{Name: &name})
}

func (useCase AoptUseCase) UpsertVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
	return useCase.vehicleRepo.Upsert(ctx, vehicle)
}

func (useCase AoptUseCase) GetModifications(ctx context.Context) ([]*models.Modification, error) {
	return useCase.modificationRepo.Get(ctx, nil)
}
