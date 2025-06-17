package entity

import "time"

type Product struct {
	Id            string
	Name          string
	Description   string
	Price         float64
	ImageFileName string
	CreatedAt     time.Time
	CreatedBy     *string
	UpdatedAt     time.Time
	UpdatedBy     *string
	DeletedAt     time.Time
	DeletedBy     *string
	IsDeleted     bool
}
