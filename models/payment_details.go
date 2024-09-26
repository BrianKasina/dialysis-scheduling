package models

type PaymentDetails struct {
    ID          int    `json:"id" db:"payment_details_id"`
    PaymentName string `json:"payment_name" db:"payment_name"`
}