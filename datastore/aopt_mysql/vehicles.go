package aopt_mysql

import (
	"aopt-db-sync/datastore/attributes"
	"aopt-db-sync/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	vehiclesFields = "id, tec_doc_id, auto_id, brand_id, name, year_from, year_to"
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

func newVehiclesDTO(in *models.Vehicle) *vehiclesDTO {
	out := &vehiclesDTO{
		ID:   in.ID,
		Name: in.Name,
	}

	if in.TecDocID != nil {
		out.TecDocID = sql.NullInt64{
			Int64: *in.TecDocID,
			Valid: true,
		}
	}

	if in.AutoID != nil {
		out.AutoID = sql.NullInt64{
			Int64: *in.AutoID,
			Valid: true,
		}
	}

	return out
}

type vehiclesDTO struct {
	ID       int64          `db:"id"`
	TecDocID sql.NullInt64  `db:"tec_doc_id"`
	AutoID   sql.NullInt64  `db:"auto_id"`
	BrandID  sql.NullInt64  `db:"brand_id"`
	Name     string         `db:"name"`
	YearFrom sql.NullString `db:"year_from"`
	YearTo   sql.NullString `db:"year_to"`
}

func (dto *vehiclesDTO) Entity() *models.Vehicle {
	out := &models.Vehicle{
		ID:   dto.ID,
		Name: dto.Name,
	}

	if dto.TecDocID.Valid {
		out.TecDocID = &dto.TecDocID.Int64
	}
	if dto.AutoID.Valid {
		out.AutoID = &dto.AutoID.Int64
	}
	if dto.YearFrom.Valid {
		t, _ := time.Parse("2006-01-02", dto.YearFrom.String)
		out.YearFrom = &t
	}
	if dto.YearTo.Valid {
		t, _ := time.Parse("2006-01-02", dto.YearTo.String)
		out.YearTo = &t
	}

	return out
}

func (repo *VehiclesRepository) Get(ctx context.Context, in *models.VehicleFilter) ([]*models.Vehicle, error) {
	query, args, err := repo.filteredQuery(in)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var dtos []vehiclesDTO
	err = repo.db.SelectContext(ctx, &dtos, fmt.Sprintf("%s;", query), args...)
	if err != nil {
		return nil, err
	}

	out := make([]*models.Vehicle, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, err
}

func (repo *VehiclesRepository) filteredQuery(in *models.VehicleFilter) (string, []interface{}, error) {
	conditions := make([]attributes.ConditionAttribute, 0, 0)
	if in != nil && in.ID != nil {
		conditions = append(conditions,
			attributes.NewStrictCondition(vehiclesTable, "id", *in.ID))
	}
	if in != nil && in.BrandID != nil {
		conditions = append(conditions,
			attributes.NewStrictCondition(vehiclesTable, "brand_id", *in.BrandID))
	}
	if in != nil && in.Name != nil {
		nameCond := attributes.NewStrictCondition("", "name", strings.ToLower(*in.Name))
		nameCond.ValueKey = "name"
		conditions = append(conditions, nameCond)
	}
	whereSets, whereMapper, err := attributes.PrepareConditionsWithoutCheckTable(" AND ", conditions...)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	searchConditions := make([]attributes.ConditionAttribute, 0, 0)
	if in != nil && in.Search != nil {
		searchConditions = append(searchConditions,
			attributes.NewCondition(vehiclesTable, "name", "LIKE", fmt.Sprintf("%s%%", *in.Search)))
	}
	searchSets, searchMapper, err := attributes.PrepareConditions(attributes.CheckTable(vehiclesTable), " AND ", searchConditions...)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	where := whereSets
	mapper := attributes.CombineMapper(whereMapper, searchMapper)

	if len(searchSets) > 0 {
		if where == "" {
			where = searchSets
		} else {
			where = fmt.Sprintf("%s AND (%s)", where, searchSets)
		}
	}

	if where != "" {
		where = "WHERE " + where
	}

	query, args, err := sqlx.Named(fmt.Sprintf(`SELECT %s FROM %s %s`, vehiclesFields, vehiclesTable, where), mapper)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return repo.db.Rebind(query), args, nil
}

func (repo *VehiclesRepository) Upsert(ctx context.Context, in *models.Vehicle) (*models.Vehicle, error) {
	dto := newVehiclesDTO(in)

	query, args, err := repo.db.BindNamed(fmt.Sprintf(`
INSERT INTO %s (id, tec_doc_id, brand_id, name, year_from, year_to)
VALUES (:id, :tec_doc_id, :brand_id, :name, :year_from, :year_to)
ON DUPLICATE KEY UPDATE auto_id = :auto_id, 
tec_doc_id = :tec_doc_id, 
name       = :name, 
year_from  = :year_from,
year_to    = :year_to;`, vehiclesTable), dto)
	if err != nil {
		return nil, err
	}

	query = repo.db.Rebind(query)

	result, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	dto.ID, _ = result.LastInsertId()

	return dto.Entity(), nil
}
