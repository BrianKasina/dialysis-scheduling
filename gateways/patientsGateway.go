package gateways

import (
    "go.mongodb.org/mongo-driver/mongo"
	"github.com/BrianKasina/dialysis-scheduling/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "fmt"
    "time"
)


type PatientGateway struct {
    collection *mongo.Collection
}

func NewPatientGateway(db *mongo.Database) *PatientGateway {
    return &PatientGateway{
        collection: db.Collection("patients"),
    }
}

func (pg *PatientGateway) GetPatients( limit, offset int) ([]models.Patient, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := pg.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var patients []models.Patient
    for cursor.Next(ctx) {
        var patient models.Patient
        if err := cursor.Decode(&patient); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

//get all patients with a non patient_history field
func (pg *PatientGateway) GetPatientsWithHistory(limit, offset int) ([]models.Patient, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := pg.collection.Find(ctx, bson.M{"patient_history": bson.M{"$exists": true}}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var patients []models.Patient
    for cursor.Next(ctx) {
        var patient models.Patient
        if err := cursor.Decode(&patient); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}


func (pg *PatientGateway) SearchPatients(query string, limit, offset int) ([]models.Patient, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": query, "$options": "i"}},
            {"address": bson.M{"$regex": query, "$options": "i"}},
            {"phone_number": bson.M{"$regex": query, "$options": "i"}},
        },
    }
    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := pg.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var patients []models.Patient
    for cursor.Next(ctx) {
        var patient models.Patient
        if err := cursor.Decode(&patient); err != nil {
            return nil, err
        }
        patients = append(patients, patient)
    }
    return patients, nil
}

func (pg *PatientGateway) GetTotalPatientCount(query string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": query, "$options": "i"}},
            {"address": bson.M{"$regex": query, "$options": "i"}},
            {"phone_number": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    count, err := pg.collection.CountDocuments(ctx, filter)
    return int(count), err
}

func (pg *PatientGateway) CreatePatient(patient *models.Patient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := pg.collection.InsertOne(ctx, patient)
    return err
}

func (pg *PatientGateway) UpdatePatient(patient *models.Patient) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"patient_id": patient.ID}
    update := bson.M{
        "$set": bson.M{
            "name":             patient.Name,
            "address":          patient.Address,
            "phone_number":     patient.PhoneNumber,
            "date_of_birth":    patient.DateOfBirth,
            "gender":           patient.Gender,
            "emergency_contact": patient.EmergencyContact,
            "payment_details_id": patient.PaymentDetailsID,
            "payment_name":     patient.PaymentName,
            "status":           patient.Status,
            "history_file":     patient.HistoryFile,
        },
    }

    result, err := pg.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no patient found with ID %d", patient.ID)
    }

    return nil
}

func (pg *PatientGateway) DeletePatient(patientID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"patient_id": patientID}
    _, err := pg.collection.DeleteOne(ctx, filter)
    return err
}