package aopt_mysql

import (
	"aopt-db-sync/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	modificationsFields = "id, vehicle_id, tec_doc_id, modification_name, construction_type, cylinder_capacity_liter, fuel_type, impulsion_type, power_hp, year_from, year_to, vehicle_type_id"
	modificationsTable  = "modifications"
)

// NewModificationsRepository create data access layer ModificationsRepository
func NewModificationsRepository(db *sqlx.DB) *ModificationsRepository {
	return &ModificationsRepository{db: db}
}

// ModificationsRepository data access layer for modifications
type ModificationsRepository struct {
	db *sqlx.DB
}

type modificationsDTO struct {
	ID                    string    `db:"id"`
	VehicleID             string    `db:"vehicle_id"`
	TecDocID              string    `db:"tec_doc_id"`
	BrandID               string    `db:"brand_id"`
	VehicleTypeID         string    `db:"vehicle_type_id"`
	ModificationName      string    `db:"modification_name"`
	ConstructionType      string    `db:"construction_type"`
	CylinderCapacityLiter string    `db:"cylinder_capacity_liter"`
	FuelType              string    `db:"fuel_type"`
	ImpulsionType         string    `db:"impulsion_type"`
	PowerHp               string    `db:"power_hp"`
	YearFrom              time.Time `db:"year_from"`
	YearTo                time.Time `db:"year_to"`
}

func (dto *modificationsDTO) Entity() *models.Modification {
	out := &models.Modification{
		ID:                    dto.ID,
		VehicleID:             dto.VehicleID,
		TecDocID:              dto.TecDocID,
		BrandID:               dto.BrandID,
		VehicleTypeID:         dto.VehicleTypeID,
		ModificationName:      dto.ModificationName,
		ConstructionType:      dto.ConstructionType,
		CylinderCapacityLiter: dto.CylinderCapacityLiter,
		FuelType:              dto.FuelType,
		ImpulsionType:         dto.ImpulsionType,
		PowerHp:               dto.PowerHp,
		YearFrom:              dto.YearFrom,
		YearTo:                dto.YearTo,
	}

	return out
}

func (repo *ModificationsRepository) GetAll(ctx context.Context) ([]*models.Modification, error) {
	var dtos []*modificationsDTO

	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s;`, modificationsFields, modificationsTable))
	if err != nil {
		return nil, err
	}

	out := make([]*models.Modification, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
