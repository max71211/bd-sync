package usecases

import (
	"aopt-db-sync/models"
	"context"
)

func NewAutoUseCase(markRepo carMarkRepo, modelRepo carModelRepo, modificationRepo carModificationRepo) *AutoUseCase {
	return &AutoUseCase{
		markRepo:         markRepo,
		modelRepo:        modelRepo,
		modificationRepo: modificationRepo,
	}
}

type carMarkRepo interface {
	GetAll(ctx context.Context) ([]*models.CarMark, error)
}
type carModelRepo interface {
	GetAll(ctx context.Context) ([]*models.CarModel, error)
	GetByCarMark(ctx context.Context, carMarkID int) ([]*models.CarModel, error)
}
type carModificationRepo interface {
	GetAll(ctx context.Context) ([]*models.CarModification, error)
	GetByCarModelID(ctx context.Context, carModelID int) ([]*models.CarModification, error)
}

type AutoUseCase struct {
	markRepo         carMarkRepo
	modelRepo        carModelRepo
	modificationRepo carModificationRepo
}

func (useCase AutoUseCase) GetMarks(ctx context.Context) ([]*models.CarMark, error) {
	return useCase.markRepo.GetAll(ctx)
}

func (useCase AutoUseCase) GetModels(ctx context.Context) ([]*models.CarModel, error) {
	return useCase.modelRepo.GetAll(ctx)
}

func (useCase AutoUseCase) GetModelsByMarkID(ctx context.Context, markID int) ([]*models.CarModel, error) {
	return useCase.modelRepo.GetByCarMark(ctx, markID)
}

func (useCase AutoUseCase) GetModifications(ctx context.Context) ([]*models.CarModification, error) {
	return useCase.modificationRepo.GetAll(ctx)
}

func (useCase AutoUseCase) GetModificationsByModelID(ctx context.Context, carModelID int) ([]*models.CarModification, error) {
	return useCase.modificationRepo.GetByCarModelID(ctx, carModelID)
}
