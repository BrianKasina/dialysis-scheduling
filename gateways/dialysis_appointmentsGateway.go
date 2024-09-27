package gateways

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/models"
	"github.com/BrianKasina/dialysis-scheduling/utils"
    "github.com/gorilla/mux"
)

// DialysisGateway handles database operations for dialysis appointments
type DialysisGateway struct {
    db *sql.DB
}

// NewDialysisGateway creates a new instance of DialysisGateway
func NewDialysisGateway(db *sql.DB) *DialysisGateway {
    return &DialysisGateway{db: db}
}

// GetAppointments retrieves dialysis appointments with patient and staff details
func (dg *DialysisGateway) GetAppointments(limit, offset int) ([]models.DialysisAppointment, error) {
    rows, err := dg.db.Query(`
        SELECT da.appointment_id, da.date, da.time, da.status, p.name AS patient_name, s.name AS staff_name
        FROM dialysis_appointment da
        JOIN patient p ON da.patient_id = p.patient_id
        JOIN hospital_staff s ON da.staff_id = s.staff_id
        LIMIT ? OFFSET ?
    ` , limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var appointments []models.DialysisAppointment
    for rows.Next() {
        var appointment models.DialysisAppointment
        if err := rows.Scan(&appointment.ID, &appointment.Date, &appointment.Time, &appointment.Status, &appointment.PatientName, &appointment.StaffName); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

// SearchAppointments searches for dialysis appointments based on a query
func (dg *DialysisGateway) SearchAppointments(query string, limit, offset int) ([]models.DialysisAppointment, error) {
    searchQuery := "%" + query + "%"
    rows, err := dg.db.Query(`
        SELECT da.appointment_id, da.date, da.time, da.status, p.name AS patient_name, s.name AS staff_name
        FROM dialysis_appointment da
        JOIN patient p ON da.patient_id = p.patient_id
        JOIN hospital_staff s ON da.staff_id = s.staff_id
        WHERE CONCAT(da.date, ' ', da.time, ' ', p.name, ' ', s.name) LIKE ? LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var appointments []models.DialysisAppointment
    for rows.Next() {
        var appointment models.DialysisAppointment
        if err := rows.Scan(&appointment.ID, &appointment.Date, &appointment.Time, &appointment.Status, &appointment.PatientName, &appointment.StaffName); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

func (dg *DialysisGateway) GetTotalDialysisAppointmentCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = dg.db.QueryRow(`
            SELECT COUNT(*)
            FROM dialysis_appointment da
            JOIN patient p ON da.patient_id = p.patient_id
            JOIN hospital_staff s ON da.staff_id = s.staff_id
            WHERE CONCAT(da.date, ' ', da.time, ' ', p.name, ' ', s.name) LIKE ?
        `, searchQuery)
    } else {
        row = dg.db.QueryRow("SELECT COUNT(*) FROM dialysis_appointment")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

// CreateAppointment creates a new dialysis appointment
func (dg *DialysisGateway) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment models.DialysisAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = dg.db.Exec(
        `INSERT INTO dialysis_appointment (
            date, time, status, patient_id, staff_id
        ) VALUES (
            ?, ?, ?, 
            (SELECT patient_id FROM patients WHERE name LIKE ? LIMIT 1), 
            (SELECT staff_id FROM hospital_staff WHERE name LIKE ? LIMIT 1)
        )`,
        appointment.Date, appointment.Time, appointment.Status,appointment.PatientName, appointment.StaffName,
    )
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create dialysis appointment")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(appointment)
}

// UpdateAppointment updates an existing dialysis appointment
func (dg *DialysisGateway) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment models.DialysisAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = dg.db.Exec(
        `UPDATE dialysis_appointment SET 
            date = ?, time = ?, status = ?, 
            staff_id = (SELECT staff_id FROM hospital_staff WHERE name LIKE ? LIMIT 1) 
            WHERE appointment_id = ?`,
        appointment.Date, appointment.Time, appointment.Status, appointment.StaffName, appointment.ID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update dialysis appointment")
        return
    }

    json.NewEncoder(w).Encode(appointment)
}

// DeleteAppointment deletes a dialysis appointment by its ID
func (dg *DialysisGateway) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    appointmentID := vars["id"]

    _, err := dg.db.Exec("DELETE FROM dialysis_appointment WHERE appointment_id = ?", appointmentID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete dialysis appointment")
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Dialysis appointment deleted successfully"})
}
