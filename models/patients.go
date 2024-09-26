package models

type Patient struct {
    ID               int    `json:"id" db:"patient_id"`
    Name             string `json:"name" db:"name"`
    Address          string `json:"address" db:"address"`
    PhoneNumber      string `json:"phone_number" db:"phone_number"`
    DateOfBirth      string `json:"date_of_birth" db:"date_of_birth"`
    Gender           string `json:"gender" db:"gender"`
    EmergencyContact string `json:"emergency_contact" db:"emergency_contact"`
    PaymentDetailsID int    `json:"payment_details_id" db:"payment_details_id"`
}