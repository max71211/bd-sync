package usecases

import (
	"aopt-db-sync/models"
	"context"
	"errors"
	"fmt"
	"github.com/alexsergivan/transliterator"
	"go.uber.org/zap"
	"log"
	"time"
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
	GetBrandByName(ctx context.Context, name string) (*models.Brand, error)
	UpsertBrand(ctx context.Context, brand *models.Brand) (*models.Brand, error)

	GetVehiclesByBrandAndName(ctx context.Context, brandID int64, name string) (*models.Vehicle, error)
	UpsertVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error)

	GetModificationByVehicleAndName(ctx context.Context, vehicleID int64, name string) (*models.Modification, error)
	UpsertModification(ctx context.Context, modification *models.Modification) (*models.Modification, error)
}

type autoDataManager interface {
	GetMarks(ctx context.Context) ([]*models.CarMark, error)
	GetModelsByMarkID(ctx context.Context, markID int64) ([]*models.CarModel, error)
	GetModificationsByModelID(ctx context.Context, carModelID int64) ([]*models.CarModification, error)
	GetCarCharacteristic(ctx context.Context, carModificationID int64) (*models.CarCharacteristic, error)
}

type SyncUseCase struct {
	aoptManager aoptDataManager
	autoManager autoDataManager
}

func (useCase SyncUseCase) SyncData(ctx context.Context, updateBrands, updateVehicles, updateModifications bool) {
	log.Println(updateBrands, "|", updateVehicles, "|", updateModifications)
	carMarks, err := useCase.autoManager.GetMarks(ctx)
	if err != nil {
		log.Fatal("Get car marks error:", zap.Error(err))
	}

	for _, cm := range carMarks {
		name := useCase.transliterate(cm.Name)
		brand, err := useCase.aoptManager.GetBrandByName(ctx, name)
		if errors.Is(err, models.ErrNoObject) {
			log.Println("get brand by name:", cm.Name)
			err = nil
		}
		if err != nil {
			log.Fatal("get brand error:", err)
		}

		// region update_brand

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

		if updateBrands {
			brand, err = useCase.aoptManager.UpsertBrand(ctx, brand)
			if err != nil {
				log.Fatal("upsert brand error", err)
			}
		}

		if brand.ID == 0 {
			continue
		}

		// endregion update_brand

		carModels, err := useCase.autoManager.GetModelsByMarkID(ctx, cm.ID)
		if err != nil {
			log.Fatal("GET models err", err)
		}

		for _, cmd := range carModels {
			name = useCase.transliterate(cmd.Name)
			vehicle, err := useCase.aoptManager.GetVehiclesByBrandAndName(ctx, brand.ID, name)
			if errors.Is(err, models.ErrNoObject) {
				log.Println("get vehicle by name:", cmd.Name)
				err = nil
			}
			if err != nil {
				log.Fatal("get vehicle by name error:", err)
			}

			// region update_vehicle

			if vehicle != nil {
				vehicle.AutoID = brand.AutoID
			} else {
				var tecDocID int64 = 0
				vehicle = &models.Vehicle{
					ID:       0,
					AutoID:   &cmd.ID,
					TecDocID: &tecDocID,
					BrandID:  brand.ID,
					Name:     cmd.Name,
				}

				if cmd.YearFrom != nil {
					start, err := time.Parse("2006-01-02", fmt.Sprintf("%d-01-01", *cmd.YearFrom))
					if err == nil {
						vehicle.YearFrom = &start
					}
				}

				if cmd.YearTo != nil {
					end, err := time.Parse("2006-01-02", fmt.Sprintf("%d-12-31", *cmd.YearTo))
					if err == nil {
						vehicle.YearTo = &end
					}
				}
			}

			if updateVehicles {
				vehicle, err = useCase.aoptManager.UpsertVehicle(ctx, vehicle)
				if err != nil {
					log.Fatal("upsert vehicle error:", err)
				}
			}

			if vehicle.ID == 0 {
				continue
			}

			// endregion update_vehicle

			autoModifications, err := useCase.autoManager.GetModificationsByModelID(ctx, cmd.ID)
			if errors.Is(err, models.ErrNoObject) {
				err = nil
			}
			if err != nil {
				log.Fatal("GET modifications err", err)
			}

			for _, mdf := range autoModifications {
				name := useCase.transliterate(mdf.Name)
				modification, err := useCase.aoptManager.GetModificationByVehicleAndName(ctx, vehicle.ID, name)
				if errors.Is(err, models.ErrNoObject) {
					log.Println("get modification by name:", mdf.Name)
					err = nil
				}

				// region update_modification

				if modification != nil {
					modification.AutoID = &mdf.ID
				} else {
					var tecDocID int64 = 0
					modification = &models.Modification{
						ID:               0,
						AutoID:           &mdf.ID,
						TecDocID:         &tecDocID,
						VehicleID:        vehicle.ID,
						BrandID:          brand.ID,
						VehicleTypeID:    0,
						ModificationName: mdf.Name,
						ConstructionType: mdf.SerieName,
					}

					if mdf.StartProductionYear != nil {
						start, err := time.Parse("2006-01-02", fmt.Sprintf("%s-01-01", *mdf.StartProductionYear))
						if err == nil {
							modification.YearFrom = &start
						}
					}
					if modification.YearFrom == nil && mdf.Generation.YearBegin != nil {
						start, err := time.Parse("2006-01-02", fmt.Sprintf("%s-01-01", *mdf.Generation.YearBegin))
						if err == nil {
							modification.YearFrom = &start
						}
					}

					if mdf.EndProductionYear != nil {
						end, err := time.Parse("2006-01-02", fmt.Sprintf("%s-12-31", *mdf.EndProductionYear))
						if err == nil {
							modification.YearTo = &end
						}
					}
					if modification.YearTo == nil && mdf.Generation.YearEnd != nil {
						end, err := time.Parse("2006-01-02", fmt.Sprintf("%s-12-31", *mdf.Generation.YearEnd))
						if err == nil {
							modification.YearTo = &end
						}
					}

					characteristic, err := useCase.autoManager.GetCarCharacteristic(ctx, mdf.ID)
					if err != nil {
						log.Println("get car characteristic error:", err)
					}

					if characteristic != nil {
						modification.FuelType = characteristic.FuelType
						modification.ImpulsionType = characteristic.ImpulsionType
						modification.PowerHp = characteristic.HorsePower
						modification.CylinderCapacityLiter = characteristic.CylinderCapacity
					}
				}

				if updateModifications {
					modification, err = useCase.aoptManager.UpsertModification(ctx, modification)
					if err != nil {
						log.Println("upsert modification error:", err)
					}
				}

				// endregion update_modification
			}
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
