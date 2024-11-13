package gateways

import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type NotificationGateway struct {
    collection *mongo.Collection
}

func NewNotificationGateway(db *mongo.Database) *NotificationGateway {
    return &NotificationGateway{
        collection: db.Collection("notifications"),
    }
}

func (ng *NotificationGateway) GetNotifications(limit, offset int) ([]models.Notification, error) {
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

    var notifications []models.Notification
    for cursor.Next(ctx) {
        var notification models.Notification
        if err := cursor.Decode(&notification); err != nil {
            return nil, err
        }
        notifications = append(notifications, notification)
    }
    return notifications, nil
}

func (ng *NotificationGateway) SearchNotifications(query string, limit, offset int) ([]models.Notification, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"message": bson.M{"$regex": query, "$options": "i"}},
            {"admin_name": bson.M{"$regex": query, "$options": "i"}},
            {"patient_name": bson.M{"$regex": query, "$options": "i"}},


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

    var notifications []models.Notification
    for cursor.Next(ctx) {
        var notification models.Notification
        if err := cursor.Decode(&notification); err != nil {
            return nil, err
        }
        notifications = append(notifications, notification)
    }
    return notifications, nil
}

func (ng *NotificationGateway) GetTotalNotificationCount(query string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"message": bson.M{"$regex": query, "$options": "i"}},
            {"admin_name": bson.M{"$regex": query, "$options": "i"}},
            {"patient_name": bson.M{"$regex": query, "$options": "i"}},
        },
    }

    count, err := ng.collection.CountDocuments(ctx, filter)
    return int(count), err
}

func (ng *NotificationGateway) CreateNotification(notification *models.Notification) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := ng.collection.InsertOne(ctx, notification)
    return err

}

func (ng *NotificationGateway) UpdateNotification(notification *models.Notification) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := ng.collection.ReplaceOne(ctx, bson.M{"notification_id": notification.ID}, notification)
    return err
}

func (ng *NotificationGateway) DeleteNotification(notificationID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := ng.collection.DeleteOne(ctx, bson.M{"notification_id": notificationID})
    return err
}