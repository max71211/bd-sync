package aopt_mysql

import (
	"aopt-db-sync/datastore/attributes"
	"aopt-db-sync/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

const (
	modificationsFields = "id, vehicle_id, tec_doc_id, auto_id, modification_name, construction_type, cylinder_capacity_liter, fuel_type, impulsion_type, power_hp, year_from, year_to, vehicle_type_id"
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
	ID                    int64         `db:"id"`
	VehicleID             int64         `db:"vehicle_id"`
	TecDocID              sql.NullInt64 `db:"tec_doc_id"`
	AutoID                sql.NullInt64 `db:"auto_id"`
	BrandID               int64         `db:"brand_id"`
	VehicleTypeID         string        `db:"vehicle_type_id"`
	ModificationName      string        `db:"modification_name"`
	ConstructionType      string        `db:"construction_type"`
	CylinderCapacityLiter string        `db:"cylinder_capacity_liter"`
	FuelType              string        `db:"fuel_type"`
	ImpulsionType         string        `db:"impulsion_type"`
	PowerHp               string        `db:"power_hp"`
	YearFrom              time.Time     `db:"year_from"`
	YearTo                time.Time     `db:"year_to"`
}

func (dto *modificationsDTO) Entity() *models.Modification {
	out := &models.Modification{
		ID:                    dto.ID,
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

	if dto.TecDocID.Valid {
		out.TecDocID = &dto.TecDocID.Int64
	}
	if dto.AutoID.Valid {
		out.AutoID = &dto.AutoID.Int64
	}

	return out
}

func (repo *ModificationsRepository) Get(ctx context.Context, in *models.ModificationFilter) ([]*models.Modification, error) {
	query, args, err := repo.filteredQuery(in)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var dtos []modificationsDTO
	err = repo.db.SelectContext(ctx, &dtos, fmt.Sprintf("%s;", query), args...)
	if err != nil {
		return nil, err
	}

	out := make([]*models.Modification, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, err
}

func (repo *ModificationsRepository) filteredQuery(in *models.ModificationFilter) (string, []interface{}, error) {
	conditions := make([]attributes.ConditionAttribute, 0, 0)
	if in != nil && in.ID != nil {
		conditions = append(conditions,
			attributes.NewStrictCondition(brandTable, "id", *in.ID))
	}
	if in != nil && in.VehicleID != nil {
		conditions = append(conditions,
			attributes.NewStrictCondition(brandTable, "vehicle_id", *in.VehicleID))
	}
	whereSets, whereMapper, err := attributes.PrepareConditionsWithoutCheckTable(" AND ", conditions...)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	searchConditions := make([]attributes.ConditionAttribute, 0, 1)
	if in != nil && in.Search != nil {
		searchConditions = append(searchConditions,
			attributes.NewCondition(brandTable, "name", "LIKE", fmt.Sprintf("%%%s%%", *in.Search)))
	}
	searchSets, searchMapper, err := attributes.PrepareConditions(attributes.CheckTable(brandTable), " AND ", searchConditions...)
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

	query, args, err := sqlx.Named(fmt.Sprintf(`SELECT %s FROM %s %s`, modificationsFields, modificationsTable, where), mapper)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return repo.db.Rebind(query), args, nil
}

func (repo *ModificationsRepository) Upsert(ctx context.Context, in *models.Vehicle) (*models.Vehicle, error) {
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
