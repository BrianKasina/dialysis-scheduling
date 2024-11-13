package models

type SystemAdmin struct {
    ID          int    `json:"id" bson:"admin_id"`
    Name        string `json:"name" bson:"name"`
    Email       string `json:"email" bson:"email"`
    PhoneNumber string `json:"phone_number" bson:"phone_number"`
}