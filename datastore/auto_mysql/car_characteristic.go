package auto_mysql

import (
	"aopt-db-sync/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

const (
	carCharacteristicFields      = `car_characteristic.id_car_characteristic, car_characteristic.name`
	carCharacteristicTable       = "car_characteristic"
	carCharacteristicValueFields = `car_characteristic_value.value, car_characteristic_value.unit`
	carCharacteristicValueTable  = "car_characteristic_value"
)

// NewCarCharacteristicRepository create data access layer CarModificationRepository
func NewCarCharacteristicRepository(db *sqlx.DB) *CarCharacteristicRepository {
	return &CarCharacteristicRepository{db: db}
}

// CarCharacteristicRepository data access layer for car modification
type CarCharacteristicRepository struct {
	db *sqlx.DB
}

type carCharacteristicDTO struct {
	ID    int64          `db:"id_car_characteristic"`
	Name  string         `db:"name"`
	Value sql.NullString `db:"value"`
	Unit  sql.NullString `db:"unit"`
}

func (repo *CarCharacteristicRepository) CarModificationCharacteristic(ctx context.Context, carModelID int64) (*models.CarCharacteristic, error) {
	var dtos []*carCharacteristicDTO
	err := sqlx.SelectContext(ctx, repo.db, &dtos,
		fmt.Sprintf(`SELECT %s  
FROM %s
JOIN %s ON car_characteristic_value.id_car_characteristic = car_characteristic.id_car_characteristic
WHERE car_characteristic_value.id_car_modification = %d 
AND car_characteristic.id_car_characteristic IN (12,13,14,27);`,
			fmt.Sprintf("%s, %s", carCharacteristicFields, carCharacteristicValueFields),
			carCharacteristicValueTable, carCharacteristicTable, carModelID),
	)
	if err != nil {
		return nil, err
	}

	out := &models.CarCharacteristic{}
	for _, d := range dtos {
		switch d.ID {
		case 12:
			if d.Value.Valid {
				out.FuelType = d.Value.String
			}
		case 13:
			if d.Value.Valid {
				cc, err := strconv.Atoi(d.Value.String)
				if err == nil {
					out.CylinderCapacity = cc
				}
			}
		case 14:
			if d.Value.Valid {
				hp, err := strconv.Atoi(d.Value.String)
				if err == nil {
					out.HorsePower = hp
				}
			}
		case 27:
			if d.Value.Valid {
				out.ImpulsionType = d.Value.String
			}
		}
	}

	return out, nil
}
