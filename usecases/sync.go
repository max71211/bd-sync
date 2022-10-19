package usecases

import (
	"aopt-db-sync/models"
	"context"
	"errors"
	"github.com/alexsergivan/transliterator"
	"go.uber.org/zap"
	"log"
)

var nameMapper = map[string]string{
	"ВАЗ (Lada)": "LADA",
	"Volkswagen": "VW",
}

func NewSyncUseCase(aoptDataManager aoptDataManager, autoDataManager autoDataManager) *SyncUseCase {
	return &SyncUseCase{
		aoptManager: aoptDataManager,
		autoManager: autoDataManager,
	}
}

type aoptDataManager interface {
	GetBrands(ctx context.Context) ([]*models.Brand, error)
	GetBrandByID(ctx context.Context, brandID string) ([]*models.Brand, error)
	GetBrandByName(ctx context.Context, name string) (*models.Brand, error)
	UpsertBrand(ctx context.Context, brand *models.Brand) (*models.Brand, error)

	GetVehicles(ctx context.Context) ([]*models.Vehicle, error)
	GetModifications(ctx context.Context) ([]*models.Modification, error)
}

type autoDataManager interface {
	GetMarks(ctx context.Context) ([]*models.CarMark, error)
	GetModels(ctx context.Context) ([]*models.CarModel, error)
	GetModifications(ctx context.Context) ([]*models.CarModification, error)
	GetModelsByMarkID(ctx context.Context, markID int) ([]*models.CarModel, error)
	GetModificationsByModelID(ctx context.Context, carModelID int) ([]*models.CarModification, error)
}

type SyncUseCase struct {
	aoptManager aoptDataManager
	autoManager autoDataManager
}

func (useCase SyncUseCase) SyncData(ctx context.Context) {
	carMarks, err := useCase.autoManager.GetMarks(ctx)
	if err != nil {
		log.Fatal("Get car marks error:", zap.Error(err))
	}

	for _, cm := range carMarks {
		cm.Name = useCase.transliterate(cm.Name)
		brand, err := useCase.aoptManager.GetBrandByName(ctx, cm.Name)
		if errors.Is(err, models.ErrNoObject) {
			err = nil
		}
		if err != nil {
			log.Fatal("get brand error:", err)
		}

		if brand != nil {
			brand.AutoID = &cm.ID
		} else {
			var tecDocID int64 = 0
			brand = &models.Brand{
				ID:       0,
				Name:     cm.Name,
				AutoID:   &cm.ID,
				TecDocID: &tecDocID,
			}
		}

		brandUpdated, err := useCase.aoptManager.UpsertBrand(ctx, brand)
		if err != nil {
			log.Fatal("upsert brand error", err)
		}

		log.Println("Updated brand", brandUpdated)

		//carModels, err := useCase.autoManager.GetModelsByMarkID(ctx, cm.ID)
		//if err != nil {
		//	log.Println("GET models err", err)
		//	continue
		//}
		//
		//log.Println("MARK:", cm.Name, "| CAR_MODELS:", len(carModels))
		//for _, cmd := range carModels {
		//	cmdf, err := useCase.autoManager.GetModificationsByModelID(ctx, cmd.ID)
		//	if err != nil {
		//		log.Println("GET modifications err", err)
		//		continue
		//	}
		//	log.Println("MODEL:", cmd.Name, "| MODEL_MODIFICATIONS:", len(cmdf))
		//}
		log.Println("################################")
	}
}

func (useCase SyncUseCase) transliterate(s string) string {
	trans := transliterator.NewTransliterator(nil)
	s = trans.Transliterate(s, "en")
	if name, ok := nameMapper[s]; ok {
		s = name
	}

	return s
}
