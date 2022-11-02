package models

import "time"

type BrandFilter struct {
	Search    *string `json:"search"`
	ID        *int64  `json:"id"`
	AutoID    *int64  `json:"auto_id"`
	TecDocID  *int64  `json:"tec_doc_id"`
	Name      *string `json:"name"`
	ISPopular *bool   `json:"is_popular"`
}

type Brand struct {
	ID        int64  `json:"id"`
	AutoID    *int64 `json:"auto_id"`
	TecDocID  *int64 `json:"tec_doc_id"`
	Name      string `json:"name"`
	ISPopular bool   `json:"is_popular"`
}

type VehicleFilter struct {
	Search   *string `json:"search"`
	ID       *int64  `json:"id"`
	TecDocID *int64  `json:"tec_doc_id"`
	AutoID   *int64  `json:"auto_id"`
	BrandID  *int64  `json:"brand_id"`
	Name     *string `json:"name"`
}

type Vehicle struct {
	ID       int64      `json:"id"`
	TecDocID *int64     `json:"tec_doc_id"`
	AutoID   *int64     `json:"auto_id"`
	BrandID  int64      `json:"brand_id"`
	Name     string     `json:"name"`
	YearFrom *time.Time `json:"year_from"`
	YearTo   *time.Time `json:"year_to"`
}

type ModificationFilter struct {
	Search                *string    `json:"search"`
	ID                    *int64     `json:"id"`
	VehicleID             *int64     `json:"vehicle_id"`
	TecDocID              *int64     `json:"tec_doc_id"`
	AutoID                *int64     `json:"auto_id"`
	BrandID               *int64     `json:"brand_id"`
	VehicleTypeID         *string    `json:"vehicle_type_id"`
	ModificationName      *string    `json:"modification_name"`
	ConstructionType      *string    `json:"construction_type"`
	CylinderCapacityLiter *string    `json:"cylinder_capacity_liter"`
	FuelType              *string    `json:"fuel_type"`
	ImpulsionType         *string    `json:"impulsion_type"`
	PowerHp               *string    `json:"power_hp"`
	YearFrom              *time.Time `json:"year_from"`
	YearTo                *time.Time `json:"year_to"`
}

type Modification struct {
	ID                    int64     `json:"id"`
	VehicleID             int64     `json:"vehicle_id"`
	TecDocID              *int64    `json:"tec_doc_id"`
	AutoID                *int64    `json:"auto_id"`
	BrandID               int64     `json:"brand_id"`
	VehicleTypeID         string    `json:"vehicle_type_id"`
	ModificationName      string    `json:"modification_name"`
	ConstructionType      string    `json:"construction_type"`
	CylinderCapacityLiter string    `json:"cylinder_capacity_liter"`
	FuelType              string    `json:"fuel_type"`
	ImpulsionType         string    `json:"impulsion_type"`
	PowerHp               string    `json:"power_hp"`
	YearFrom              time.Time `json:"year_from"`
	YearTo                time.Time `json:"year_to"`
}
