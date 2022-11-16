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
	carModelFields     = "car_model.id_car_model, id_car_mark, name, name_rus, date_create, date_update, id_car_type"
	carModelTable      = "car_model"
	carModelYearFields = `min(year.year) year_from, max(year.year) year_to`
	carModelYearTable  = `year`
)

// NewCarModelRepository create data access layer CarMarRepository
func NewCarModelRepository(db *sqlx.DB) *CarModelRepository {
	return &CarModelRepository{db: db}
}

// CarModelRepository data access layer for car models
type CarModelRepository struct {
	db *sqlx.DB
}

type carModelDTO struct {
	ID         int64          `db:"id_car_model"`
	IDCarMar   int64          `db:"id_car_mark"`
	Name       string         `db:"name"`
	NameRU     sql.NullString `db:"name_rus"`
	DateCreate int64          `db:"date_create"`
	DateUpdate int64          `db:"date_update"`
	IDCarType  int            `db:"id_car_type"`
	YearFrom   int            `db:"year_from"`
	YearTo     int            `db:"year_to"`
}

func (dto *carModelDTO) Entity() *models.CarModel {
	out := &models.CarModel{
		ID:         dto.ID,
		Name:       dto.Name,
		DateCreate: time.Unix(dto.DateCreate, 0),
		DateUpdate: time.Unix(dto.DateUpdate, 0),
		IDCarType:  dto.IDCarType,
		YearFrom:   &dto.YearFrom,
	}

	if dto.NameRU.Valid {
		out.NameRU = &dto.NameRU.String
	}

	if dto.YearTo > dto.YearFrom {
		out.YearTo = &dto.YearTo
	}

	return out
}

func (repo *CarModelRepository) GetAll(ctx context.Context) ([]*models.CarModel, error) {
	var dtos []*carModelDTO
	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s;`, carModelFields, carModelTable))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarModel, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}

func (repo *CarModelRepository) GetByCarMark(ctx context.Context, carMarkID int64) ([]*models.CarModel, error) {
	var dtos []*carModelDTO
	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s 
FROM %s 
JOIN %s ON car_model.id_car_model = year.id_car_model
WHERE id_car_mark = %d 
GROUP BY %s;`, fmt.Sprintf("%s, %s", carModelFields, carModelYearFields),
			carModelTable, carModelYearTable, carMarkID, carModelFields))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarModel, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
