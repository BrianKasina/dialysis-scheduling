package models


type Patient struct {
	ID               int    `json:"id" db:"patient_id"`
	Name             string `json:"name" db:"name"`
	Address          string `json:"address" db:"address"`
	PhoneNumber      string `json:"phone_number" db:"phone_number"`
	DateOfBirth      string `json:"date_of_birth" db:"date_of_birth"`
	Gender           string `json:"gender" db:"gender"`
	EmergencyContact string `json:"emergency_contact" db:"emergency_contact"`
	PaymentDetailsID int    `json:"payment_details_id,omitempty" db:"payment_details_id"`
	PaymentName      string `json:"payment_name,omitempty" db:"payment_name"` // Added field for join
	Message          string `json:"message,omitempty" db:"message"`           // Added field for join
	Date_sent        string `json:"sent_date,omitempty" db:"sent_date"`       // Added field for join
	Time_sent        string `json:"sent_time,omitempty" db:"sent_time"`       // Added field for join
	Date             string `json:"date,omitempty" db:"date"`                 // Added field for join
	Time             string `json:"time,omitempty" db:"time"`                 // Added field for join
    Status           string `json:"status,omitempty" db:"status"`             // Added field for join
}