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
	carModelFields = "id_car_model, id_car_mark, name, name_rus, date_create, date_update, id_car_type"
	carModelTable  = "car_model"
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
	ID         string         `db:"id_car_model"`
	IDCarMar   string         `db:"id_car_mark"`
	Name       string         `db:"name"`
	NameRU     sql.NullString `db:"name_rus"`
	DateCreate int64          `db:"date_create"`
	DateUpdate int64          `db:"date_update"`
	IDCarType  int            `db:"id_car_type"`
}

func (dto *carModelDTO) Entity() *models.CarModel {
	out := &models.CarModel{
		ID:         dto.ID,
		Name:       dto.Name,
		DateCreate: time.Unix(dto.DateCreate, 0),
		DateUpdate: time.Unix(dto.DateUpdate, 0),
		IDCarType:  dto.IDCarType,
	}

	if dto.NameRU.Valid {
		out.NameRU = &dto.NameRU.String
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

func (repo *CarModelRepository) GetByCarMark(ctx context.Context, carMarkID int) ([]*models.CarModel, error) {
	var dtos []*carModelDTO
	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s WHERE id_car_mark = %d;`, carModelFields, carModelTable, carMarkID))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarModel, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
