package models

type HospitalStaff struct {
    ID             int    `json:"id" db:"staff_id"`
    Name           string `json:"name" db:"name"`
    Gender         string `json:"gender" db:"gender"`
    Specialization string `json:"specialization" db:"specialization"`
    PhoneNumber    string `json:"phone_number" db:"phonenumber"`
    Status         string `json:"status" db:"status"`
}