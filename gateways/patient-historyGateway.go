package gateways

import (
	"database/sql"
)

type PatientHistoryGateway struct {
    db *sql.DB
}

func NewPatientHistoryGateway(db *sql.DB) *PatientHistoryGateway {
    return &PatientHistoryGateway{db: db}
}

// Create a patient history record in the database
func (phg *PatientHistoryGateway) CreatePatientHistory(patientName string, historyFile string) error {
// Use a subquery to get the patient ID from the patient's name
    _, err := phg.db.Exec(`
        INSERT INTO patient_history (patient_id, history_file)
        VALUES (
            (SELECT patient_id FROM patient WHERE name = ?),
            ?
        )
    `, patientName, historyFile)
    return err
	}



