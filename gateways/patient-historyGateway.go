package gateways

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PatientHistoryGateway struct {
    collection *mongo.Collection
    collection2 *mongo.Collection
}

func NewPatientHistoryGateway(db *mongo.Database) *PatientHistoryGateway {
    return &PatientHistoryGateway{
        collection: db.Collection("patient_history"),
        collection2: db.Collection("patients"),
    }
}

//create a patient history file, adding a new patient history file record to the patient whose name mathces the patientName
func (phg *PatientHistoryGateway) CreatePatientHistory(patientName string, patientHistoryFile string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := phg.collection.InsertOne(ctx, bson.M{"patient_name": patientName, "patient_history_file": patientHistoryFile})
    if err != nil {
        return err
    }

    //update the patients collection to include the patient history file
    _, err = phg.collection2.UpdateOne(ctx, bson.M{"name": patientName}, bson.M{"$set": bson.M{"history_file": patientHistoryFile}})
    if err != nil {
        return err
    }
    return nil
}

//delete a patient history file, removing the patient history file record from the patient whose name mathces the patientName
func (phg *PatientHistoryGateway) DeletePatientHistory(patientName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := phg.collection.DeleteOne(ctx, bson.M{"patient name": patientName})
    if err != nil {
        return err
    }

    //update the patients collection to remove the patient history file
    _, err = phg.collection2.UpdateOne(ctx, bson.M{"name": patientName}, bson.M{"$unset": bson.M{"history_file": ""}})
    if err != nil {
        return err
    }
    return nil
}

// Add or update patient history files for a specific patient
func (phg *PatientHistoryGateway) CreateOrUpdatePatientHistory(patientName string, files []string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Update or insert history in `patient_history` collection
    _, err := phg.collection.UpdateOne(ctx, bson.M{"patient_name": patientName},
        bson.M{"$set": bson.M{"patient_history_files": files}}, options.Update().SetUpsert(true))
    if err != nil {
        return err
    }

    // Update history reference in `patients` collection
    _, err = phg.collection2.UpdateOne(ctx, bson.M{"name": patientName},
        bson.M{"$set": bson.M{"history_files": files}})
    return err
}




