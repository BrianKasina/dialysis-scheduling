package models

type PaymentDetails struct {
    ID          int    `json:"id" bson:"payment_details_id"`
    PaymentName string `json:"payment_name" bson:"payment_name"`
}