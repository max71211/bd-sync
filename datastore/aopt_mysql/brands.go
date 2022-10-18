package aopt_mysql

import (
	"aopt-db-sync/datastore/attributes"
	"aopt-db-sync/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"strings"
)

const (
	brandFields = "id, auto_id, tec_doc_id, name, is_popular"
	brandTable  = "brands"
)

// NewBrandsRepository create data access layer VehiclesRepository
func NewBrandsRepository(db *sqlx.DB) *BrandsRepository {
	return &BrandsRepository{db: db}
}

// BrandsRepository data access layer for vehicles
type BrandsRepository struct {
	db *sqlx.DB
}

func newBrandDTO(in *models.Brand) *brandsDTO {
	out := &brandsDTO{
		ID:        in.ID,
		Name:      in.Name,
		ISPopular: in.ISPopular,
	}

	if in.TecDocID != nil {
		out.TecDocID = sql.NullString{
			String: *in.TecDocID,
			Valid:  true,
		}
	}

	if in.AutoID != nil {
		out.AutoID = sql.NullString{
			String: *in.AutoID,
			Valid:  true,
		}
	}

	return out
}

type brandsDTO struct {
	ID        string         `db:"id"`
	AutoID    sql.NullString `db:"auto_id"`
	TecDocID  sql.NullString `db:"tec_doc_id"`
	Name      string         `db:"name"`
	ISPopular bool           `db:"is_popular"`
}

func (dto *brandsDTO) Entity() *models.Brand {
	out := &models.Brand{
		ID:        dto.ID,
		Name:      dto.Name,
		ISPopular: dto.ISPopular,
	}

	if dto.TecDocID.Valid {
		out.TecDocID = &dto.TecDocID.String
	}
	if dto.AutoID.Valid {
		out.AutoID = &dto.AutoID.String
	}

	return out
}

func (repo *BrandsRepository) Get(ctx context.Context, in *models.BrandFilter) ([]*models.Brand, error) {
	query, args, err := repo.filteredQuery(in)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var dtos []brandsDTO
	err = repo.db.SelectContext(ctx, &dtos, fmt.Sprintf("%s;", query), args...)
	if err != nil {
		return nil, err
	}

	out := make([]*models.Brand, 0, 0)
	for _, d := range dtos {
		out = append(out, d.Entity())
	}

	return out, err
}

func (repo *BrandsRepository) filteredQuery(in *models.BrandFilter) (string, []interface{}, error) {
	conditions := make([]attributes.ConditionAttribute, 0, 6)
	if in != nil && in.ID != nil {
		conditions = append(conditions,
			attributes.NewStrictCondition(brandTable, "id", *in.ID))
	}
	if in != nil && in.Name != nil {
		nameCond := attributes.NewStrictCondition("", "LOWER(name)", strings.ToLower(*in.Name))
		nameCond.ValueKey = "name"
		conditions = append(conditions, nameCond)
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

	query, args, err := sqlx.Named(fmt.Sprintf(`SELECT %s FROM %s %s`, brandFields, brandTable, where), mapper)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return repo.db.Rebind(query), args, nil
}

func (repo *BrandsRepository) Upsert(ctx context.Context, in *models.Brand) (*models.Brand, error) {
	dto := newBrandDTO(in)

	query, args, err := repo.db.BindNamed(fmt.Sprintf(`
INSERT INTO %s (auto_id, tec_doc_id, name, is_popular)
VALUES (:auto_id, :tec_doc_id, :name, :is_popular)
ON DUPLICATE KEY UPDATE auto_id = :auto_id, 
tec_doc_id = :tec_doc_id, 
name = :name, 
is_popular = :is_popular;`, brandTable), dto)
	if err != nil {
		return nil, err
	}

	log.Println("QUERY", query, "|", args)

	//result, err := repo.db.ExecContext(ctx, query, args)
	//if err != nil {
	//	return nil, err
	//}
	//
	//bID, _ := result.LastInsertId()
	//dto.ID = fmt.Sprintf("%d", bID)

	return dto.Entity(), nil
}
