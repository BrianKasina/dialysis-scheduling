package gateways

import (
    "database/sql"
	"github.com/BrianKasina/dialysis-scheduling/models"
)


type PatientGateway struct {
    db *sql.DB
}

func NewPatientGateway(db *sql.DB) *PatientGateway {
    return &PatientGateway{db: db}
}

func (pg *PatientGateway) GetPatients( limit, offset int) ([]models.Patient, error) {
    rows, err := pg.db.Query("SELECT * FROM patient LIMIT ? OFFSET ?", limit, offset)
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

func (pg *PatientGateway) SearchPatients(query string, limit, offset int) ([]models.Patient, error) {
    searchQuery := "%" + query + "%"
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, pd.payment_name
        FROM patient p
        LEFT JOIN payment_details pd ON p.payment_details_id = pd.payment_details_id
        WHERE CONCAT(p.name, ' ', p.address, ' ', p.phone_number) LIKE ? LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
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
func (pg *PatientGateway) GetPatientsWithPayment(limit, offset int) ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, pd.payment_name
        FROM patient p
        JOIN payment_details pd ON p.payment_details_id = pd.payment_details_id
        LIMIT ? OFFSET ?
    `, limit, offset)
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

func (pg *PatientGateway) GetPatientsWithDialysisAppointments(limit, offset int) ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, da.date , da.time , da.status 
        FROM patient p
        JOIN dialysis_appointment da ON p.patient_id = da.patient_id LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.Date, &patient.Time, &patient.Status); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetPatientsWithNephrologistAppointments(limit, offset int) ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, na.date, na.time, na.status
        FROM patient p
        JOIN nephrologist_appointment na ON p.patient_id = na.patient_id 
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.Date, &patient.Time, &patient.Status); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetPatientsWithNotifications(limit, offset int) ([]models.Patient, error) {
    rows, err := pg.db.Query(`
        SELECT p.patient_id, p.name, p.address, p.phone_number, p.date_of_birth, p.gender, p.emergency_contact, n.message, n.sent_date AS Date_sent, n.sent_time AS Time_sent
        FROM patient p
        JOIN notifications n ON p.patient_id = n.patient_id
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var patients []models.Patient
    for rows.Next() {
        var patient models.Patient
        if err := rows.Scan(&patient.ID, &patient.Name, &patient.Address, &patient.PhoneNumber, &patient.DateOfBirth, &patient.Gender, &patient.EmergencyContact, &patient.Date_sent, &patient.Time_sent, &patient.Message); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetTotalPatientCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = pg.db.QueryRow(`
            SELECT COUNT(*)
            FROM patient p
            LEFT JOIN payment_details pd ON p.payment_details_id = pd.payment_details_id
            WHERE CONCAT(p.name, ' ', p.address, ' ', p.phone_number) LIKE ?
        `, searchQuery)
    } else {
        row = pg.db.QueryRow("SELECT COUNT(*) FROM patient")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
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