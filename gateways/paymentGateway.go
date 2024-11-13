package gateways

import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type PaymentDetailsGateway struct {
    collection *mongo.Collection
}

func NewPaymentDetailsGateway(db *mongo.Database) *PaymentDetailsGateway {
    return &PaymentDetailsGateway{
        collection: db.Collection("payment_details"),
    }
}

func (pg *PaymentDetailsGateway) GetPaymentDetails(limit, offset int) ([]models.PaymentDetails, error) {
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

    var paymentDetails []models.PaymentDetails
    for cursor.Next(ctx) {
        var paymentDetail models.PaymentDetails
        if err := cursor.Decode(&paymentDetail); err != nil {
            return nil, err
        }
        paymentDetails = append(paymentDetails, paymentDetail)
    }
    return paymentDetails, nil
}

func (pg *PaymentDetailsGateway) SearchPaymentDetails(query string, limit, offset int) ([]models.PaymentDetails, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"payment_name": bson.M{"$regex": query, "$options": "i"}},
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

    var paymentDetails []models.PaymentDetails
    for cursor.Next(ctx) {
        var paymentDetail models.PaymentDetails
        if err := cursor.Decode(&paymentDetail); err != nil {
            return nil, err
        }
        paymentDetails = append(paymentDetails, paymentDetail)
    }
    return paymentDetails, nil
}

func (pg *PaymentDetailsGateway) GetTotalPaymentDetailsCount(query string) (int, error) {
    return 0, nil
}

func (pg *PaymentDetailsGateway) CreatePaymentDetail(paymentDetail *models.PaymentDetails) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := pg.collection.InsertOne(ctx, paymentDetail)
    return err
}

func (pg *PaymentDetailsGateway) UpdatePaymentDetail(paymentDetail *models.PaymentDetails) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"payment_details_id": paymentDetail.ID}
    update := bson.M{"$set": paymentDetail}
    _, err := pg.collection.UpdateOne(ctx, filter, update)
    return err
}

func (pg *PaymentDetailsGateway) DeletePaymentDetail(paymentDetailID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"payment_details_id": paymentDetailID}
    _, err := pg.collection.DeleteOne(ctx, filter)
    return err
}