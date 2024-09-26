package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/utils"
)

func GetAppointments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    rows, err := db.Query("SELECT * FROM dialysis_appointment")
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch appointments")
        return
    }
    defer rows.Close()

    var appointments []map[string]interface{}
    for rows.Next() {
        var appointment map[string]interface{}
        err := rows.Scan(&appointment["appointment_id"], &appointment["date"], &appointment["time"], &appointment["status"], &appointment["patient_id"], &appointment["staff_id"])
        if err != nil {
            utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to scan appointment")
            return
        }
        appointments = append(appointments, appointment)
    }

    json.NewEncoder(w).Encode(appointments)
}

func CreateAppointment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var appointment map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = db.Exec("INSERT INTO dialysis_appointment (date, time, status, patient_id, staff_id) VALUES (?, ?, ?, ?, ?)",
        appointment["date"], appointment["time"], appointment["status"], appointment["patient_id"], appointment["staff_id"])
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create appointment")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Appointment created successfully"})
}

func UpdateAppointment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var appointment map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = db.Exec("UPDATE dialysis_appointment SET date = ?, time = ?, status = ?, patient_id = ?, staff_id = ? WHERE appointment_id = ?",
        appointment["date"], appointment["time"], appointment["status"], appointment["patient_id"], appointment["staff_id"], appointment["appointment_id"])
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update appointment")
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Appointment updated successfully"})
}

func DeleteAppointment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    vars := mux.Vars(r)
    appointmentID := vars["id"]

    _, err := db.Exec("DELETE FROM dialysis_appointment WHERE appointment_id = ?", appointmentID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete appointment")
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Appointment deleted successfully"})
}