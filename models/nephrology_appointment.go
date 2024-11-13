package models

type NephrologistAppointment struct {
    ID        int    `json:"id" bson:"appointment_id"`
    Date      string `json:"date" bson:"date"`
    Time      string `json:"time" bson:"time"`
    Status    string `json:"status" bson:"status"`
    PatientID int    `json:"patient_id,omitempty" bson:"patient_id"`
    StaffID   int    `json:"staff_id,omitempty" bson:"staff_id"`
    StaffName string `json:"staff_name,omitempty" bson:"staff_name"`
    PatientName string `json:"patient_name,omitempty" bson:"patient_name"`
}