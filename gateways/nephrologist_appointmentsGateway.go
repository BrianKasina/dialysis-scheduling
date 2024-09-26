package gateways

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/models"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "github.com/gorilla/mux"
)

type NephrologistAppointmentGateway struct {
    db *sql.DB
}

// Initialize Nephrologist Gateway
func NewNephrologistAppointmentGateway(db *sql.DB) *NephrologistAppointmentGateway {
    return &NephrologistAppointmentGateway{db: db}
}

// Retrieve nephrologist appointments with joined patient and staff data
func (ng *NephrologistAppointmentGateway) GetAppointments() ([]models.NephrologistAppointment, error) {
    rows, err := ng.db.Query(`
        SELECT na.appointment_id, na.date, na.time, na.status, p.name AS patient_name, s.name AS staff_name
        FROM nephrologist_appointment na
        JOIN patient p ON na.patient_id = p.patient_id
        JOIN hospital_staff s ON na.staff_id = s.staff_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var appointments []models.NephrologistAppointment
    for rows.Next() {
        var appointment models.NephrologistAppointment
        if err := rows.Scan(&appointment.ID, &appointment.Date, &appointment.Time, &appointment.Status, &appointment.PatientName, &appointment.StaffName); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return appointments, nil
}

// SearchAppointments searches for nephrologist appointments based on a query
func (ng *NephrologistAppointmentGateway) SearchAppointments(query string) ([]models.NephrologistAppointment, error) {
    searchQuery := "%" + query + "%"
    rows, err := ng.db.Query(`
        SELECT na.appointment_id, na.date, na.time, na.status, p.name AS patient_name, s.name AS staff_name
        FROM nephrologist_appointment na
        JOIN patient p ON na.patient_id = p.patient_id
        JOIN hospital_staff s ON na.staff_id = s.staff_id
        WHERE CONCAT(na.date, ' ', na.time, ' ', p.name, ' ', s.name) LIKE ?
    `, searchQuery)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var appointments []models.NephrologistAppointment
    for rows.Next() {
        var appointment models.NephrologistAppointment
        if err := rows.Scan(&appointment.ID, &appointment.Date, &appointment.Time, &appointment.Status, &appointment.PatientName, &appointment.StaffName); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

// Create new nephrologist appointment
func (ng *NephrologistAppointmentGateway) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment models.NephrologistAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = ng.db.Exec(`
        INSERT INTO nephrologist_appointment (date, time, status, patient_id, staff_id)
        VALUES (?, ?, ?, ?, ?)`,
        appointment.Date, appointment.Time, appointment.Status, appointment.PatientID, appointment.StaffID,
    )
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create nephrologist appointment")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(appointment)
}

// Update or cancel nephrologist appointment
func (ng *NephrologistAppointmentGateway) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var appointment models.NephrologistAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = ng.db.Exec(`
        UPDATE nephrologist_appointment 
        SET date = ?, time = ?, status = ?, staff_id = ?
        WHERE appointment_id = ?`,
        appointment.Date, appointment.Time, appointment.Status, appointment.StaffID, id,
    )
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update nephrologist appointment")
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(appointment)
}

// Delete nephrologist appointment
func (ng *NephrologistAppointmentGateway) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    _, err := ng.db.Exec("DELETE FROM nephrologist_appointment WHERE appointment_id = ?", id)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete nephrologist appointment")
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Nephrologist appointment deleted successfully"})
}