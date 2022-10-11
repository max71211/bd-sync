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
	carMarkFields = "id_car_mark, name, name_rus, date_create, date_update, id_car_type"
	carMarkTable  = "car_mark"
)

// NewCarMarRepository create data access layer CarMarRepository
func NewCarMarRepository(db *sqlx.DB) *CarMarRepository {
	return &CarMarRepository{db: db}
}

// CarMarRepository data access layer for users
type CarMarRepository struct {
	db *sqlx.DB
}

type carMarkDTO struct {
	ID         int            `db:"id_car_mark"`
	Name       string         `db:"name"`
	NameRU     sql.NullString `db:"name_rus"`
	DateCreate int64          `db:"date_create"`
	DateUpdate int64          `db:"date_update"`
	IDCarType  int            `db:"id_car_type"`
}

func (dto *carMarkDTO) Entity() *models.CarMark {
	out := &models.CarMark{
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

func (repo *CarMarRepository) GetAll(ctx context.Context) ([]*models.CarMark, error) {
	var dtos []*carMarkDTO

	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s FROM %s;`, carMarkFields, carMarkTable))
	if err != nil {
		return nil, err
	}

	out := make([]*models.CarMark, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, nil
}
