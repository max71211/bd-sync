package models

import "time"

type BrandFilter struct {
	Search    *string `json:"search"`
	ID        *string `json:"id"`
	AutoID    *string `json:"auto_id"`
	TecDocID  *string `json:"tec_doc_id"`
	Name      *string `json:"name"`
	ISPopular *bool   `json:"is_popular"`
}

type Brand struct {
	ID        string  `json:"id"`
	AutoID    *string `json:"auto_id"`
	TecDocID  *string `json:"tec_doc_id"`
	Name      string  `json:"name"`
	ISPopular bool    `json:"is_popular"`
}

type VehicleFilter struct {
	ID       *string `json:"id"`
	TecDocID *string `json:"tec_doc_id"`
	AutoID   *string `json:"auto_id"`
	BrandID  *string `json:"brand_id"`
	Name     *string `json:"name"`
}

type Vehicle struct {
	ID       string    `json:"id"`
	TecDocID string    `json:"tec_doc_id"`
	AutoID   string    `json:"auto_id"`
	BrandID  string    `json:"brand_id"`
	Name     string    `json:"name"`
	YearFrom time.Time `json:"year_from"`
	YearTo   time.Time `json:"year_to"`
}

type ModificationFilter struct {
	ID                    *string    `json:"id"`
	VehicleID             *string    `json:"vehicle_id"`
	TecDocID              *string    `json:"tec_doc_id"`
	AutoID                *string    `json:"auto_id"`
	BrandID               *string    `json:"brand_id"`
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
	ID                    string    `json:"id"`
	VehicleID             string    `json:"vehicle_id"`
	TecDocID              string    `json:"tec_doc_id"`
	AutoID                string    `json:"auto_id"`
	BrandID               string    `json:"brand_id"`
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
