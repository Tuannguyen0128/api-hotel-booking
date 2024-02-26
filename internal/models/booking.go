package models

import (
	"time"
)

type Booking struct {
	ID            string    `json:"id"`
	GuestID       string    `json:"guest_id"`
	AccountID     string    `json:"account_id"`
	CheckinTime   time.Time `json:"checkin_time"`
	CheckoutTime  time.Time `json:"checkout_time"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}
