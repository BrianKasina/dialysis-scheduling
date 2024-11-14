package gateways

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/BrianKasina/dialysis-scheduling/models"
)

type HospitalStaffGateway struct {
    collection *mongo.Collection
}

func NewHospitalStaffGateway(db *mongo.Database) *HospitalStaffGateway {
    return &HospitalStaffGateway{
        collection: db.Collection("hospital_staff"),
    }
}

func (hsg *HospitalStaffGateway) GetHospitalStaff(limit, offset int) ([]models.HospitalStaff, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := hsg.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var staff []models.HospitalStaff
    for cursor.Next(ctx) {
        var member models.HospitalStaff
        if err := cursor.Decode(&member); err != nil {
            return nil, err
        }
        staff = append(staff, member)
    }
    return staff, nil
}

func (hsg *HospitalStaffGateway) SearchHospitalStaff(query string, limit, offset int) ([]models.HospitalStaff, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": query, "$options": "i"}},
            {"specialization": bson.M{"$regex": query, "$options": "i"}},
            {"phonenumber": bson.M{"$regex": query, "$options": "i"}},
        },
    }
    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := hsg.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }

    var staff []models.HospitalStaff
    for cursor.Next(ctx) {
        var member models.HospitalStaff
        if err := cursor.Decode(&member); err != nil {
            return nil, err
        }
        staff = append(staff, member)
    }
    return staff, nil
}

func (hsg *HospitalStaffGateway) GetTotalStaffCount(query string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": query, "$options": "i"}},
            {"specialization": bson.M{"$regex": query, "$options": "i"}},
            {"phonenumber": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    count, err := hsg.collection.CountDocuments(ctx, filter)
    return int(count), err
}

func (hsg *HospitalStaffGateway) CreateHospitalStaff(member *models.HospitalStaff) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := hsg.collection.InsertOne(ctx, member)
    return err
}

func (hsg *HospitalStaffGateway) UpdateHospitalStaff(staff *models.HospitalStaff) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"staff_id": staff.ID}
    update := bson.M{
        "$set": bson.M{
            "name":             staff.Name,
            "phone_number":     staff.PhoneNumber,
            "gender":           staff.Gender,
            "status":           staff.Status,
            "specialization":   staff.Specialization,
        },
    }

    result, err := hsg.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no staff found with ID %d", staff.ID)
    }

    return nil
}

func (hsg *HospitalStaffGateway) DeleteHospitalStaff(staffID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := hsg.collection.DeleteOne(ctx, bson.M{"staff_id": staffID})
    return err
}