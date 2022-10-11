package auto_mysql

import (
	"aopt-db-sync/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	carModificationFields = `id_car_modification, id_car_serie, id_car_model, car_modification.name, start_production_year, end_production_year, date_create, date_update, id_car_type`
	carModificationTable  = "car_modification"
	carSerieFields        = `car_serie.name as serie_name`
	carSerieTable         = "car_serie"
	carGenerationFields   = `car_generation.name as generation_name, car_generation.year_begin as generation_year_begin, car_generation.year_end as generation_year_end`
	carGenerationTable    = "car_generation"
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
	ID                  int            `db:"id_car_modification"`
	IDCarSerie          int            `db:"id_car_serie"`
	IDCarModel          int            `db:"id_car_model"`
	Name                string         `db:"name"`
	StartProductionYear *int           `db:"start_production_year"`
	EndProductionYear   *int           `db:"end_production_year"`
	SerieName           string         `db:"serie_name"`
	GenerationName      string         `db:"generation_name"`
	GenerationYearBegin string         `db:"generation_year_begin"`
	GenerationYearEnd   sql.NullString `db:"generation_year_end"`
	DateCreate          int64          `db:"date_create"`
	DateUpdate          int64          `db:"date_update"`
	IDCarType           int            `db:"id_car_type"`
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
		fmt.Sprintf(`SELECT %s FROM %s 
JOIN %s ON %[1]s.id_car_serie = %[2]s.id_car_serie
JOIN %s ON %[2]s.id_car_generation = %[3]s.id_car_generation
WHERE id_car_model = %d;`, fmt.Sprintf("%s, %s, %s", carModificationFields, carSerieFields, carGenerationFields), carModificationTable, carSerieTable, carGenerationTable, carModelID))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarModification, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
