package models

type Patient struct {
    ID               int    `json:"id" bson:"patient_id"`
    Name             string `json:"name" bson:"name"`
    Address          string `json:"address" bson:"address"`
    PhoneNumber      string `json:"phone_number" bson:"phone_number"`
    DateOfBirth      string `json:"date_of_birth" bson:"date_of_birth"`
    Gender           string `json:"gender" bson:"gender"`
    EmergencyContact string `json:"emergency_contact" bson:"emergency_contact"`
    PaymentDetailsID int    `json:"payment_details_id,omitempty" bson:"payment_details_id"`
    PaymentName      string `json:"payment_name,omitempty" bson:"payment_name"`
    Status           string `json:"status,omitempty" bson:"status"`
    HistoryFile      string `json:"history_file,omitempty" bson:"history_file"`
}