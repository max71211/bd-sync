package auto_mysql

import (
	"aopt-db-sync/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	carModificationFields = "id_car_modification, id_car_serie, id_car_model, name, start_production_year, end_production_year, date_create, date_update, id_car_type"
	carModificationTable  = "car_modification"
)

// NewCarModificationRepository create data access layer CarModificationRepository
func NewCarModificationRepository(db *sqlx.DB) *CarModificationRepository {
	return &CarModificationRepository{db: db}
}

// CarModificationRepository data access layer for car modification
type CarModificationRepository struct {
	db *sqlx.DB
}

type carModificationDTO struct {
	ID                  int    `db:"id_car_modification"`
	IDCarSerie          int    `db:"id_car_serie"`
	IDCarModel          int    `db:"id_car_model"`
	Name                string `db:"name"`
	StartProductionYear *int   `db:"start_production_year"`
	EndProductionYear   *int   `db:"end_production_year"`
	DateCreate          int64  `db:"date_create"`
	DateUpdate          int64  `db:"date_update"`
	IDCarType           int    `db:"id_car_type"`
}

func (dto *carModificationDTO) Entity() *models.CarModification {
	out := &models.CarModification{
		ID:                  dto.ID,
		IDCarSerie:          dto.IDCarSerie,
		IDCarModel:          dto.IDCarModel,
		Name:                dto.Name,
		StartProductionYear: dto.StartProductionYear,
		EndProductionYear:   dto.EndProductionYear,
		DateCreate:          time.Unix(dto.DateCreate, 0),
		DateUpdate:          time.Unix(dto.DateUpdate, 0),
		IDCarType:           dto.IDCarType,
	}

	return out
}

func (repo *CarModificationRepository) GetAll(ctx context.Context) ([]*models.CarModification, error) {
	var dtos []*carModificationDTO
	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s;`, carModelFields, carModelTable))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarModification, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}

func (repo *CarModificationRepository) GetByCarModelID(ctx context.Context, carModelID int) ([]*models.CarModification, error) {
	var dtos []*carModificationDTO
	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s WHERE id_car_model = %d;`, carModificationFields, carModificationTable, carModelID))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarModification, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
