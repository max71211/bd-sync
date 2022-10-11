package models

import (
	"time"
)

type CarMark struct {
	ID         int       `json:"id_car_mark"`
	Name       string    `json:"name"`
	NameRU     *string   `json:"name_rus"`
	DateCreate time.Time `json:"date_create"`
	DateUpdate time.Time `json:"date_update"`
	IDCarType  int       `json:"id_car_type"`
}
