package gateways

import (
    "database/sql"
	"github.com/BrianKasina/dialysis-scheduling/models"
)

// type Patient struct {
//     ID               int    `json:"id"`
//     Name             string `json:"name"`
//     Address          string `json:"address"`
//     PhoneNumber      string `json:"phone_number"`
//     DateOfBirth      string `json:"date_of_birth"`
//     Gender           string `json:"gender"`
//     EmergencyContact string `json:"emergency_contact"`
//     PaymentDetailsID int    `json:"payment_details_id"`
// }

type PatientGateway struct {
    db *sql.DB
}

func NewPatientGateway(db *sql.DB) *PatientGateway {
    return &PatientGateway{db: db}
}

func (pg *PatientGateway) GetPatients() ([]models.Patient, error) {
    rows, err := pg.db.Query("SELECT * FROM patient")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentDetailsID); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) SearchPatients(query string) ([]models.Patient, error) {
    searchQuery := "%" + query + "%"
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, pd.payment_name
        FROM patient p
        LEFT JOIN payment_details pd ON p.payment_details_id = pd.payment_details_id
        WHERE CONCAT(p.name, ' ', p.address, ' ', p.phone_number) LIKE ?
    `, searchQuery)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentName); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}
func (pg *PatientGateway) GetPatientsWithPayment() ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, pd.payment_name
        FROM patient p
        JOIN payment_details pd ON p.payment_details_id = pd.payment_details_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentName); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetPatientsWithDialysisAppointments() ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, da.date AS appointment_date, da.time AS appointment_time, da.status AS appointment_status
        FROM patient p
        JOIN dialysis_appointment da ON p.patient_id = da.patient_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentDetailsID); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetPatientsWithNephrologistAppointments() ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, na.date AS appointment_date, na.time AS appointment_time, na.status AS appointment_status
        FROM patient p
        JOIN nephrologist_appointment na ON p.patient_id = na.patient_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentDetailsID); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetPatientsWithNotifications() ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, n.message AS notification_message, n.sent_date AS notification_date, n.sent_time AS notification_time
        FROM patient p
        JOIN notifications n ON p.patient_id = n.patient_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.PaymentDetailsID); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) CreatePatient(patient *models.Patient) error {
    _, err := pg.db.Exec("INSERT INTO patient (name, address, phone_number, date_of_birth, gender, emergency_contact, payment_details_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
        patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.Gender, patient.EmergencyContact, patient.PaymentDetailsID)
    return err
}

func (pg *PatientGateway) UpdatePatient(patient *models.Patient) error {
    _, err := pg.db.Exec("UPDATE patient SET name = ?, address = ?, phone_number = ?, date_of_birth = ?, gender = ?, emergency_contact = ?, payment_details_id = ? WHERE patient_id = ?",
        patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.Gender, patient.EmergencyContact, patient.PaymentDetailsID, patient.ID)
    return err
}

func (pg *PatientGateway) DeletePatient(patientID string) error {
    _, err := pg.db.Exec("DELETE FROM patient WHERE patient_id = ?", patientID)
    return err
}