package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BrianKasina/dialysis-scheduling/models"
	"github.com/BrianKasina/dialysis-scheduling/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DialysisGateway handles database operations for dialysis appointments
type DialysisGateway struct {
    collection *mongo.Collection
}

// NewDialysisGateway creates a new instance of DialysisGateway
func NewDialysisGateway(db *mongo.Database) *DialysisGateway {
    return &DialysisGateway{
        collection: db.Collection("dialysis_appointments"),
    }
}

// GetAppointments retrieves dialysis appointments with patient and staff details
func (dg *DialysisGateway) GetAppointments(limit, offset int) ([]models.DialysisAppointment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := dg.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var appointments []models.DialysisAppointment
    for cursor.Next(ctx) {
        var appointment models.DialysisAppointment
        if err := cursor.Decode(&appointment); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

// SearchAppointments searches for dialysis appointments based on a query
func (dg *DialysisGateway) SearchAppointments(query string, limit, offset int) ([]models.DialysisAppointment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"staff_name": bson.M{"$regex": query, "$options": "i"}},
            {"patient_name": bson.M{"$regex": query, "$options": "i"}},
            {"date": bson.M{"$regex": query, "$options": "i"}},
            {"time": bson.M{"$regex": query, "$options": "i"}},
            {"status": bson.M{"$regex": query, "$options": "i"}},
        },
    }
    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := dg.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var appointments []models.DialysisAppointment
    for cursor.Next(ctx) {
        var appointment models.DialysisAppointment
        if err := cursor.Decode(&appointment); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

func (dg *DialysisGateway) GetTotalDialysisAppointmentCount(query string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"staff_name": bson.M{"$regex": query, "$options": "i"}},
            {"patient_name": bson.M{"$regex": query, "$options": "i"}},
            {"date": bson.M{"$regex": query, "$options": "i"}},
            {"time": bson.M{"$regex": query, "$options": "i"}},
            {"status": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    count, err := dg.collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    return int(count), nil
}

// CreateAppointment creates a new dialysis appointment
func (dg *DialysisGateway) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var appointment models.DialysisAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = dg.collection.InsertOne(ctx, appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create dialysis appointment")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(appointment)

}

func (dg *DialysisGateway) UpdateAppointment(appointment *models.DialysisAppointment) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"appointment_id": appointment.ID}
    update := bson.M{
        "$set": bson.M{
            "date":         appointment.Date,
            "time":         appointment.Time,
            "status":       appointment.Status,
            "staff_name":   appointment.StaffName,
            "patient_name": appointment.PatientName,
        },
    }

    result, err := dg.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no appointment found with ID %d", appointment.ID)
    }

    return nil
}

// DeleteAppointment deletes a dialysis appointment by its ID
func (dg *DialysisGateway) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    vars := mux.Vars(r)
    appointmentID := vars["id"]

    _, err := dg.collection.DeleteOne(ctx, bson.M{"appointment_id": appointmentID})
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete dialysis appointment")
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Dialysis appointment deleted successfully"})
}
