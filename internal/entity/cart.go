package entity

import "time"

type UserCart struct {
	Id        string
	UserId    string
	ProductId string
	Quantity  int
	CreatedAt time.Time
	CreatedBy *string
	UpdatedAt *time.Time
	UpdatedBy *string

	Product *Product
}
