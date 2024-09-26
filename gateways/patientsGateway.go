package gateways

import (
    "database/sql"
)

type Patient struct {
    ID               int    `json:"id"`
    Name             string `json:"name"`
    Address          string `json:"address"`
    PhoneNumber      string `json:"phone_number"`
    DateOfBirth      string `json:"date_of_birth"`
    Gender           string `json:"gender"`
    EmergencyContact string `json:"emergency_contact"`
    PaymentDetailsID int    `json:"payment_details_id"`
}

type PatientGateway struct {
    db *sql.DB
}

func NewPatientGateway(db *sql.DB) *PatientGateway {
    return &PatientGateway{db: db}
}

func (pg *PatientGateway) GetPatients() ([]Patient, error) {
    rows, err := pg.db.Query("SELECT * FROM patient")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []Patient
    for rows.Next() {
        var patient Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentDetailsID); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) CreatePatient(patient *Patient) error {
    _, err := pg.db.Exec("INSERT INTO patient (name, address, phone_number, date_of_birth, gender, emergency_contact, payment_details_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
        patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.Gender, patient.EmergencyContact, patient.PaymentDetailsID)
    return err
}

func (pg *PatientGateway) UpdatePatient(patient *Patient) error {
    _, err := pg.db.Exec("UPDATE patient SET name = ?, address = ?, phone_number = ?, date_of_birth = ?, gender = ?, emergency_contact = ?, payment_details_id = ? WHERE patient_id = ?",
        patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.Gender, patient.EmergencyContact, patient.PaymentDetailsID, patient.ID)
    return err
}

func (pg *PatientGateway) DeletePatient(patientID string) error {
    _, err := pg.db.Exec("DELETE FROM patient WHERE patient_id = ?", patientID)
    return err
}