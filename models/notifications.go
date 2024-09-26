package models

type Notification struct {
    ID        int    `json:"id" db:"notification_id"`
    Message   string `json:"message" db:"message"`
    SentDate  string `json:"sent_date" db:"sent_date"`
    SentTime  string `json:"sent_time" db:"sent_time"`
    AdminID   int    `json:"admin_id" db:"admin_id"`
    PatientID int    `json:"patient_id" db:"patient_id"`
}