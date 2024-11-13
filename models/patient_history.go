package models

type PatientHistory struct {
    ID          int    `json:"id" bson:"history_id"`
    PatientID   int    `json:"patient_id" bson:"patient_id"`
    PatientName string `json:"patient_name" bson:"patient_name"`
    HistoryFile string `json:"history_file" bson:"history_file"`
}