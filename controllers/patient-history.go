package controllers

import (
    "encoding/json"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "database/sql"
)

type PatientHistoryController struct {
    PatientHistoryGateway *gateways.PatientHistoryGateway
}

func NewPatientHistoryController(db *sql.DB) *PatientHistoryController {
    return &PatientHistoryController{
        PatientHistoryGateway: gateways.NewPatientHistoryGateway(db),
    }
}

// Upload patient history file
func (phc *PatientHistoryController) UploadPatientHistory(w http.ResponseWriter, r *http.Request) {
    // Get patient_id from the form
    patientName := r.FormValue("patient_name")

    // Parse multipart form, the file should be uploaded as 'file'
	err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Error parsing form data")
        return
    }
    file, handler, err := r.FormFile("file")
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Error reading file")
        return
    }
    defer file.Close()

    // Save file to the "patients-history-folder"
    filePath := filepath.Join("patients-history-folder", handler.Filename)
    f, err := os.Create(filePath)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to save file")
        return
    }
    defer f.Close()

    // Copy file content to the new file
    _, err = io.Copy(f, file)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Error saving file content")
        return
    }

    // Update the patient's history in the database
    err = phc.PatientHistoryGateway.CreatePatientHistory(patientName, handler.Filename)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update patient history")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}

// Open a specific history file
func (phc *PatientHistoryController) OpenPatientHistory(w http.ResponseWriter, r *http.Request) {
    historyFile := r.URL.Query().Get("filename")

    // Generate full path to the file
    filePath := filepath.Join("patients-history-folder", historyFile)

    // Check if file exists
    _, err := os.Stat(filePath)
    if os.IsNotExist(err) {
        utils.ErrorHandler(w, http.StatusNotFound, nil, "File not found")
        return
    }

	// Set headers to prompt the client to open the file in their default PDF reader
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "inline; filename="+historyFile)

    // Serve the file
    http.ServeFile(w, r, filePath)
}
