package models

type Notification struct {
    ID          int    `json:"id" bson:"notification_id"`
    Message     string `json:"message" bson:"message"`
    SentDate    string `json:"sent_date" bson:"sent_date"`
    SentTime    string `json:"sent_time" bson:"sent_time"`
    AdminID     int    `json:"admin_id,omitempty" bson:"admin_id"`
    AdminName   string `json:"admin_name,omitempty" bson:"admin_name"`
    PatientID   int    `json:"patient_id,omitempty" bson:"patient_id"`
    PatientName string `json:"patient_name,omitempty" bson:"patient_name"`
}