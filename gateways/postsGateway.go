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

type PostGateway struct {
    collection *mongo.Collection
}

func NewPostGateway(db *mongo.Database) *PostGateway {
    return &PostGateway{
        collection: db.Collection("posts"),
    }
}

func (pg *PostGateway) GetPosts(limit, offset int) ([]models.Post, error) {
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

    var posts []models.Post
    for cursor.Next(ctx) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (pg *PostGateway) SearchPosts(query string, limit, offset int) ([]models.Post, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{
        "$or": []bson.M{
            {"title": bson.M{"$regex": query, "$options": "i"}},
            {"content": bson.M{"$regex": query, "$options": "i"}},
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

    var posts []models.Post
    for cursor.Next(ctx) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (pg *PostGateway) GetTotalPostCount(query string) (int, error) {
    return 0, nil
}

func (pg *PostGateway) CreatePost(post *models.Post) error {
    _, err := pg.collection.InsertOne(context.Background(),
        bson.M{
            "title": post.Title,
            "content": post.Content,
            "post_date": post.PostDate,
            "post_time": post.PostTime,
            "admin_name": post.AdminName,
        })
    return err
}

func (pg *PostGateway) UpdatePost(post *models.Post) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"post_id": post.ID}
    update := bson.M{
        "$set": bson.M{
            "title":      post.Title,
            "content":    post.Content,
            "post_date":  post.PostDate,
            "post_time":  post.PostTime,
            "admin_name": post.AdminName,
        },
    }

    result, err := pg.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("no post found with ID %d", post.ID)
    }

    return nil
}

func (pg *PostGateway) DeletePost(postID string) error {
    _, err := pg.collection.DeleteOne(context.Background(), bson.M{"post_id": postID})
    return err
}