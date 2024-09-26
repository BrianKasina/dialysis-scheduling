package models

type PatientHistory struct {
    ID          int    `json:"id" db:"history_id"`
    PatientID   int    `json:"patient_id" db:"patient_id"`
    HistoryFile string `json:"history_file" db:"history_file"`
}