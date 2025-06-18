package entity

import "time"

type UserCart struct {
	Id        string     `db:"id"`
	UserId    string     `db:"user_id"`
	ProductId string     `db:"product_id"`
	Quantity  int        `db:"quantity"`
	CreatedAt time.Time  `db:"created_at"`
	CreatedBy *string    `db:"created_by"`
	UpdatedAt *time.Time `db:"updated_at"` // <-- ini harus pointer
	UpdatedBy *string    `db:"updated_by"`
}
