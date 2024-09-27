package models

type Notification struct {
    ID        int    `json:"id" db:"notification_id"`
    Message   string `json:"message" db:"message"`
    SentDate  string `json:"sent_date" db:"sent_date"`
    SentTime  string `json:"sent_time" db:"sent_time"`
    AdminID   int    `json:"admin_id,omitempty" db:"admin_id"`
    AdminName string `json:"admin_name,omitempty" db:"admin_name"`
    PatientID int    `json:"patient_id,omitempty" db:"patient_id"`
    PatientName string `json:"patient_name,omitempty" db:"patient_name"`
}