package gateways

import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/models"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "github.com/gorilla/mux"
)

type NephrologistAppointmentGateway struct {
    collection *mongo.Collection
    collection2 *mongo.Collection
}

// Initialize Nephrologist Gateway
func NewNephrologistAppointmentGateway(db *mongo.Database) *NephrologistAppointmentGateway {
    return &NephrologistAppointmentGateway{
        collection: db.Collection("nephrologist_appointments"),
        collection2: db.Collection("patients"),
    }
}

// Retrieve nephrologist appointments with joined patient and staff data
func (ng *NephrologistAppointmentGateway) GetAppointments(limit, offset int) ([]models.NephrologistAppointment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := ng.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var appointments []models.NephrologistAppointment
    for cursor.Next(ctx) {
        var appointment models.NephrologistAppointment
        if err := cursor.Decode(&appointment); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

// SearchAppointments searches for nephrologist appointments based on a query
func (ng *NephrologistAppointmentGateway) SearchAppointments(query string, limit, offset int ) ([]models.NephrologistAppointment, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"date": bson.M{"$regex": query, "$options": "i"}},
            {"time": bson.M{"$regex": query, "$options": "i"}},
            {"status": bson.M{"$regex": query, "$options": "i"}},
            {"patient_name": bson.M{"$regex": query, "$options": "i"}},
            {"staff_name": bson.M{"$regex": query, "$options": "i"}},
        },
    }
    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := ng.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var appointments []models.NephrologistAppointment
    for cursor.Next(ctx) {
        var appointment models.NephrologistAppointment
        if err := cursor.Decode(&appointment); err != nil {
            return nil, err
        }
        appointments = append(appointments, appointment)
    }
    return appointments, nil
}

func (ng *NephrologistAppointmentGateway) GetTotalNephrologistAppointmentCount(query string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"date": bson.M{"$regex": query, "$options": "i"}},
            {"time": bson.M{"$regex": query, "$options": "i"}},
            {"status": bson.M{"$regex": query, "$options": "i"}},
            {"patient_name": bson.M{"$regex": query, "$options": "i"}},
            {"staff_name": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    count, err := ng.collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    return int(count), nil
}

// Create new nephrologist appointment
func (ng *NephrologistAppointmentGateway) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var appointment models.NephrologistAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = ng.collection.InsertOne(ctx, appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create nephrologist appointment")
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(appointment)

}

// Update or cancel nephrologist appointment
func (ng *NephrologistAppointmentGateway) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var appointment models.NephrologistAppointment
    err := json.NewDecoder(r.Body).Decode(&appointment)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    _, err = ng.collection.UpdateOne(ctx, bson.M{"appointment_id": appointment.ID}, bson.M{"$set": appointment})
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update nephrologist appointment")
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(appointment)
}

// Delete nephrologist appointment
func (ng *NephrologistAppointmentGateway) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    vars := mux.Vars(r)
    appointmentID := vars["id"]
    if appointmentID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing appointment ID")
        return
    }

    _, err := ng.collection.DeleteOne(ctx, bson.M{"appointment_id": appointmentID})
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete nephrologist appointment")
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Nephrologist appointment deleted successfully"})
}