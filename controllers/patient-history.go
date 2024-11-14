package controllers

import (
    "archive/zip"
    "encoding/json"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "go.mongodb.org/mongo-driver/mongo"
)

type PatientHistoryController struct {
    PatientHistoryGateway *gateways.PatientHistoryGateway
}

func NewPatientHistoryController(db *mongo.Database) *PatientHistoryController {
    return &PatientHistoryController{
        PatientHistoryGateway: gateways.NewPatientHistoryGateway(db),
    }
}

func (phc *PatientHistoryController) HandlePatientHistory(w http.ResponseWriter, r *http.Request) {
    operation := r.URL.Query().Get("identifier")

    switch operation {
    case "list":
        phc.ListPatientHistory(w, r)
    case "download":
        phc.DownloadPatientHistoryZip(w, r)
    default:
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Invalid operation")
    }
}

// Upload or Update patient history files in a patient-specific folder
func (phc *PatientHistoryController) UploadPatientHistory(w http.ResponseWriter, r *http.Request) {
    patientName := r.FormValue("patient_name")

    // Parse the multipart form to retrieve files
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Error parsing form data")
        return
    }

    // Ensure patient-specific folder exists
    patientFolder := filepath.Join("patients-history-folder", patientName)
    os.MkdirAll(patientFolder, os.ModePerm)

    // Prepare a slice to hold the file names
    var fileNames []string

    // Loop through files and save them in the patient's folder
    files := r.MultipartForm.File["files"]
    for _, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            utils.ErrorHandler(w, http.StatusBadRequest, err, "Error reading file")
            return
        }
        defer file.Close()

        filePath := filepath.Join(patientFolder, fileHeader.Filename)
        f, err := os.Create(filePath)
        if err != nil {
            utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to save file")
            return
        }
        defer f.Close()

        _, err = io.Copy(f, file)
        if err != nil {
            utils.ErrorHandler(w, http.StatusInternalServerError, err, "Error saving file content")
            return
        }

        // Append the file name to the list
        fileNames = append(fileNames, fileHeader.Filename)
    }

    // Update the patient's history in the database with the file names
    err = phc.PatientHistoryGateway.CreateOrUpdatePatientHistory(patientName, fileNames)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update patient history")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Files uploaded successfully"})
}


// List contents of a patient's folder
func (phc *PatientHistoryController) ListPatientHistory(w http.ResponseWriter, r *http.Request) {
    patientName := r.URL.Query().Get("patient_name")
    patientFolder := filepath.Join("patients-history-folder", patientName)

    files, err := os.ReadDir(patientFolder)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Error reading patient history folder")
        return
    }

    fileNames := []string{}
    for _, file := range files {
        if !file.IsDir() {
            fileNames = append(fileNames, file.Name())
        }
    }

    json.NewEncoder(w).Encode(fileNames)
}

// Download patient folder as a zip file
func (phc *PatientHistoryController) DownloadPatientHistoryZip(w http.ResponseWriter, r *http.Request) {
    patientName := r.URL.Query().Get("patient_name")
    patientFolder := filepath.Join("patients-history-folder", patientName)

    zipFileName := patientName + ".zip"
    w.Header().Set("Content-Disposition", "attachment; filename="+zipFileName)
    w.Header().Set("Content-Type", "application/zip")

    zipWriter := zip.NewWriter(w)
    defer zipWriter.Close()

    err := filepath.Walk(patientFolder, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        f, err := zipWriter.Create(info.Name())
        if err != nil {
            return err
        }

        _, err = io.Copy(f, file)
        return err
    })

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Error creating zip file")
        return
    }
}
