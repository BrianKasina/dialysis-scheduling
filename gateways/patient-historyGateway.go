package gateways

import (
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "context"
    "time"
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




