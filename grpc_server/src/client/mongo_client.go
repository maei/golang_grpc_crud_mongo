package client

import (
	"context"
	"github.com/maei/shared_utils_go/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	MongoClient     mongoClientInterface     = &mongoClient{}
	MongoDatabase   mongoDatabaseInterface   = &mongoDatabase{}
	MongoCollection mongoCollectionInterface = &mongoCollection{}
)

type mongoClientInterface interface {
	setClient(client *mongo.Client) *mongo.Client
	Disconnect(context.Context) error
}

type mongoDatabaseInterface interface {
	setDB(client *mongo.Client) *mongo.Database
}

type mongoCollectionInterface interface {
	setCollection(db *mongo.Database)
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

type mongoClient struct {
	client *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	col *mongo.Collection
}

func InitClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Error("Error while connecting to Database", err)
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Error("Error while connect with context", err)
		panic(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	c := MongoClient.setClient(client)
	return c

}

func InitCollection(db *mongo.Database) {
	MongoCollection.setCollection(db)
}

func (collection *mongoCollection) setCollection(db *mongo.Database) {
	collection.col = db.Collection("test")
}

func InitDatabase(client *mongo.Client) *mongo.Database {
	d := MongoDatabase.setDB(client)
	return d
}

func (d *mongoDatabase) setDB(client *mongo.Client) *mongo.Database {
	d.db = client.Database("mydb")
	return d.db

}

func (c *mongoClient) setClient(client *mongo.Client) *mongo.Client {
	c.client = client
	return c.client
}

func (c *mongoClient) Disconnect(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	if err != nil {
		logger.Error("Error disconnecting MongoDB Client", err)
		panic(err)
	}
	return nil
}

func (collection *mongoCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	res, err := collection.col.InsertOne(ctx, document)
	if err != nil {
		logger.Error("MongoDB: Cannot InsertOne operation", err)
		return nil, err
	}
	return res, err
}
