package models

type DialysisAppointment struct {
    ID        int    `json:"id" db:"appointment_id"`
    Date      string `json:"date" db:"date"`
    Time      string `json:"time" db:"time"`
    Status    string `json:"status" db:"status"`
    PatientID int    `json:"patient_id,omitempty" db:"patient_id"`
    StaffID   int    `json:"staff_id,omitempty" db:"staff_id"`
    StaffName string `json:"staff_name,omitempty" db:"staff_name"`
    PatientName string `json:"patient_name,omitempty" db:"patient_name"`
}