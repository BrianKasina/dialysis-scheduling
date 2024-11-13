package models

type HospitalStaff struct {
    ID             int    `json:"id" bson:"staff_id"`
    Name           string `json:"name" bson:"name"`
    Gender         string `json:"gender" bson:"gender"`
    Specialization string `json:"specialization" bson:"specialization"`
    PhoneNumber    string `json:"phone_number" bson:"phone_number"`
    Status         string `json:"status" bson:"status"`
}