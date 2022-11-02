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
	GetBrandByID(ctx context.Context, brandID int64) (*models.Brand, error)
	GetBrandByName(ctx context.Context, name string) (*models.Brand, error)
	UpsertBrand(ctx context.Context, brand *models.Brand) (*models.Brand, error)

	GetVehicles(ctx context.Context) ([]*models.Vehicle, error)
	GetVehiclesByBrandAndName(ctx context.Context, brandID int64, name string) (*models.Vehicle, error)
	UpsertVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error)

	GetModifications(ctx context.Context) ([]*models.Modification, error)
}

type autoDataManager interface {
	GetMarks(ctx context.Context) ([]*models.CarMark, error)
	GetModels(ctx context.Context) ([]*models.CarModel, error)
	GetModifications(ctx context.Context) ([]*models.CarModification, error)
	GetModelsByMarkID(ctx context.Context, markID int64) ([]*models.CarModel, error)
	GetModificationsByModelID(ctx context.Context, carModelID int64) ([]*models.CarModification, error)
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

	for i, cm := range carMarks {
		log.Println("##### BRAND:", cm.Name, "#######")
		log.Println("################################")
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

		// region update_brand

		//brandUpdated, err := useCase.aoptManager.UpsertBrand(ctx, brand)
		//if err != nil {
		//	log.Fatal("upsert brand error", err)
		//}
		//log.Println("Updated brand", brandUpdated)

		// endregion update_brand

		carModels, err := useCase.autoManager.GetModelsByMarkID(ctx, cm.ID)
		if err != nil {
			log.Fatal("GET models err", err)
		}

		for _, cmd := range carModels {
			name := useCase.transliterate(cmd.Name)
			vehicle, err := useCase.aoptManager.GetVehiclesByBrandAndName(ctx, brand.ID, name)
			if errors.Is(err, models.ErrNoObject) {
				log.Println("get vehicle by name:", cmd.Name)
				err = nil
			}
			if err != nil {
				log.Fatal("get vehicle by name error:", err)
			}

			if vehicle != nil {
				vehicle.AutoID = brand.AutoID
			} else {
				vehicle = &models.Vehicle{
					AutoID:   brand.AutoID,
					TecDocID: brand.TecDocID,
					BrandID:  brand.ID,
					Name:     cmd.Name,
				}
			}

			log.Println("Model:", cmd.Name, "Vehicle:", vehicle.Name)

			// region update_vehicle

			//updatedVehicle, err := useCase.aoptManager.UpsertVehicle(ctx, vehicle)
			//if err != nil {
			//	log.Fatal("upsert vehicle error", err)
			//}
			//log.Println("updated vehicle", updatedVehicle)

			// endregion update_vehicle

			mdf, err := useCase.autoManager.GetModificationsByModelID(ctx, cmd.ID)
			if err != nil {
				log.Fatal("GET modifications err", err)
			}

			log.Println("MODEL:", cmd.Name, "| MODEL_MODIFICATIONS:", len(mdf))
		}
		log.Println()
		if i == 5 {
			break
		}
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
