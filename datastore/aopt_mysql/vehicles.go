package aopt_mysql

import (
	"aopt-db-sync/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	vehiclesFields = "id, tec_doc_id, brand_id, name, year_from, year_to"
	vehiclesTable  = "vehicles"
)

// NewVehiclesRepository create data access layer VehiclesRepository
func NewVehiclesRepository(db *sqlx.DB) *VehiclesRepository {
	return &VehiclesRepository{db: db}
}

// VehiclesRepository data access layer for vehicles
type VehiclesRepository struct {
	db *sqlx.DB
}

type vehiclesDTO struct {
	ID       string    `db:"id"`
	TecDocID string    `db:"tec_doc_id"`
	BrandID  string    `db:"brand_id"`
	Name     string    `db:"name"`
	YearFrom time.Time `db:"year_from"`
	YearTo   time.Time `db:"year_to"`
}

func (dto *vehiclesDTO) Entity() *models.Vehicle {
	out := &models.Vehicle{
		ID:       dto.ID,
		TecDocID: dto.TecDocID,
		BrandID:  dto.BrandID,
		Name:     dto.Name,
		YearFrom: dto.YearFrom,
		YearTo:   dto.YearTo,
	}

	return out
}

func (repo *VehiclesRepository) GetAll(ctx context.Context) ([]*models.Vehicle, error) {
	var dtos []*vehiclesDTO

	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s;`, vehiclesFields, vehiclesTable))
	if err != nil {
		return nil, err
	}

	out := make([]*models.Vehicle, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
