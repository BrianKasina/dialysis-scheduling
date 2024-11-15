package utils

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
)

type Database struct {
    Host     string
    Name     string
    User     string
    Port     string
    Password string
    Client   *mongo.Client
}

func NewDatabase(port,host, name, user, password string) (*Database, error) {
    db := &Database{Host: host, Name: name, User: user, Password: password, Port: port}
    err := db.connect()
    if err != nil {
        return nil, err
    }
    return db, nil
}

func (db *Database) connect() error {
    // Use the MongoDB Stable API version
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)

    // Format the URI with the necessary credentials and options
    uri := fmt.Sprintf("mongodb+srv://%s:%s@%s.pbt7o.mongodb.net/?retryWrites=true&w=majority&appName=dialysis-database",
    db.User, db.Password, db.Name)

    clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

    // Set up a context with a timeout for the connection attempt
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    // Ping the database to ensure the connection is successful
    if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
        return fmt.Errorf("ping failed: %v", err)
    }

    db.Client = client
    fmt.Println("Successfully connected to MongoDB!")
    return nil
}



// GetConnection returns the MongoDB database
func (db *Database) GetConnection() (*mongo.Database, error) {
    if db.Client == nil {
        return nil, fmt.Errorf("no MongoDB client initialized")
    }
    return db.Client.Database(db.Name), nil
}

// Close closes the MongoDB connection
func (db *Database) Close() error {
    if db.Client == nil {
        return fmt.Errorf("no MongoDB client initialized")
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return db.Client.Disconnect(ctx)
}