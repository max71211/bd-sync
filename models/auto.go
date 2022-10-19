package models

import (
	"time"
)

type CarMark struct {
	ID         int64     `json:"id_car_mark"`
	Name       string    `json:"name"`
	NameRU     *string   `json:"name_rus"`
	DateCreate time.Time `json:"date_create"`
	DateUpdate time.Time `json:"date_update"`
	IDCarType  int       `json:"id_car_type"`
}

type CarModel struct {
	ID         string    `json:"id_car_model"`
	Name       string    `json:"name"`
	NameRU     *string   `json:"name_rus"`
	DateCreate time.Time `json:"date_create"`
	DateUpdate time.Time `json:"date_update"`
	IDCarType  int       `json:"id_car_type"`
}

type CarModification struct {
	ID                  int           `json:"id_car_modification"`
	IDCarSerie          int           `json:"id_car_serie"`
	IDCarModel          int           `json:"id_car_model"`
	Name                string        `json:"name"`
	StartProductionYear *int          `json:"start_production_year"`
	EndProductionYear   *int          `json:"end_production_year"`
	DateCreate          time.Time     `json:"date_create"`
	DateUpdate          time.Time     `json:"date_update"`
	IDCarType           int           `json:"id_car_type"`
	SerieName           string        `json:"serie_name"`
	Generation          CarGeneration `json:"generation"`
}

type CarGeneration struct {
	ID        int     `json:"id_car_generationn"`
	Name      string  `json:"name"`
	YearBegin *string `json:"year_begin"`
	YearEnd   *string `json:"year_end"`
	IDCarType int     `json:"id_car_type"`
}
