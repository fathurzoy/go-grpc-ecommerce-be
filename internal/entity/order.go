package entity

import "time"

const (
	OrderStatusCodeUnpaid   = "unpaid"
	OrderStatusCodePaid     = "paid"
	OrderStatusCodeShipped  = "shipped"
	OrderStatusCodeCanceled = "canceled"
	OrderStatusCodeDone     = "done"
	OrderStatusCodeExpired  = "expired"
)

type Order struct {
	Id                   string
	Number               string
	UserId               string
	OrderStatusCode      string
	UserFullName         string
	Address              string
	PhoneNumber          string
	Notes                *string
	Total                float64
	ExpiredAt            *time.Time
	CreatedAt            time.Time
	CreatedBy            string
	UpdatedAt            *time.Time
	UpdatedBy            *string
	DeletedAt            *time.Time
	DeletedBy            *string
	IsDeleted            bool
	XenditInvoiceId      *string
	XenditInvoiceUrl     *string
	XenditPaidAt         *time.Time
	XenditPaymentMethod  *string
	XenditPaymentChannel *string

	Items []*OrderItem
}

type OrderItem struct {
	Id                   string
	ProductId            string
	ProductName          string
	ProductImageFileName string
	ProductPrice         float64
	Quantity             int64
	OrderId              string
	CreatedAt            time.Time
	CreatedBy            string
	UpdatedAt            *time.Time
	UpdatedBy            *string
	DeletedAt            *time.Time
	DeletedBy            *string
	IsDeleted            bool
}
