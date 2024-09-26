package models

type NephrologistAppointment struct {
    ID        int    `json:"id" db:"appointment_id"`
    Date      string `json:"date" db:"date"`
    Time      string `json:"time" db:"time"`
    Status    string `json:"status" db:"status"`
    PatientID int    `json:"patient_id" db:"patient_id"`
    StaffID   int    `json:"staff_id" db:"staff_id"`
}