package models

type SystemAdmin struct {
    ID          int    `json:"id" db:"admin_id"`
    Name        string `json:"name" db:"name"`
    Email       string `json:"email" db:"email"`
    PhoneNumber string `json:"phone_number" db:"phonenumber"`
}