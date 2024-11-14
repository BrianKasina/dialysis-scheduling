package gateways

import (
	"context"
	"fmt"
	"time"

	"github.com/BrianKasina/dialysis-scheduling/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminGateway struct {
    collection *mongo.Collection
}

func NewAdminGateway(db *mongo.Database) *AdminGateway {
    return &AdminGateway{
        collection: db.Collection("system_admin"),
    }
}

func (ag *AdminGateway) GetAdmins(limit, offset int) ([]models.SystemAdmin, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := ag.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var admins []models.SystemAdmin
    for cursor.Next(ctx) {
        var admin models.SystemAdmin
        if err := cursor.Decode(&admin); err != nil {
            return nil, err
        }
        admins = append(admins, admin)
    }
    return admins, nil
}

func (ag *AdminGateway) SearchAdmins(query string, limit, offset int) ([]models.SystemAdmin, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": query, "$options": "i"}},
            {"email": bson.M{"$regex": query, "$options": "i"}},
            {"phonenumber": bson.M{"$regex": query, "$options": "i"}},
        },
    }
    opts := options.Find()
    opts.SetLimit(int64(limit))
    opts.SetSkip(int64(offset))

    cursor, err := ag.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var admins []models.SystemAdmin
    for cursor.Next(ctx) {
        var admin models.SystemAdmin
        if err := cursor.Decode(&admin); err != nil {
            return nil, err
        }
        admins = append(admins, admin)
    }
    return admins, nil  
}

func (ag *AdminGateway) GetTotalAdminCount(query string) (int, error) {
    //calculate the total number of documents in the admin collection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()

    filter := bson.M{}
    if query != "" {
        filter = bson.M{
            "$or": []bson.M{
                {"name": bson.M{"$regex": query, "$options": "i"}},
                {"email": bson.M{"$regex": query, "$options": "i"}},
                {"phonenumber": bson.M{"$regex": query, "$options": "i"}},
            },
        }
    }

    count, err := ag.collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    return int(count), nil
}

func (ag *AdminGateway) CreateAdmin(admin *models.SystemAdmin) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := ag.collection.InsertOne(ctx, admin)
    return err
}

func (ag *AdminGateway) UpdateAdmin(admin *models.SystemAdmin) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"admin_id": admin.ID}
    update := bson.M{
        "$set": bson.M{
            "name":        admin.Name,
            "email":       admin.Email,
            "phonenumber": admin.PhoneNumber,
        },
    }

    result, err := ag.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no admin found with ID %d", admin.ID)
    }

    return nil
}

func (ag *AdminGateway) DeleteAdmin(adminID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"admin_id": adminID}
    _, err := ag.collection.DeleteOne(ctx, filter)
    return err
}