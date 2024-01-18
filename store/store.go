package store

import (
	"context"
	"fmt"
	"jiachen/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	client *mongo.Client
	dbName string
}

type collection[T any] struct {
	name            string
	mongoCollection *mongo.Collection
}

func (s *store) setClient(client *mongo.Client) {
	s.client = client
}

var storeInstance = &store{
	dbName: "jiachen",
}

func NewStore(uri string) *store {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	storeInstance.setClient(client)
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{bson.E{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return storeInstance
}

func (c *collection[T]) Init(s *store) {
	if c.name == "" {
		panic("name undefined")
	}
	c.mongoCollection = s.client.Database(s.dbName).Collection(c.name)
}

func (c *collection[T]) Find(filter interface{}, offset, limit int64) ([]*T, error) {
	cursor, err := c.mongoCollection.Find(context.TODO(), filter, &options.FindOptions{Limit: &limit, Skip: &offset})
	if err != nil {
		return nil, err
	}
	var results []*T
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results, nil
}

func (c *collection[T]) FindOne(filter interface{}) (*T, error) {
	var result = new(T)
	err := c.mongoCollection.FindOne(context.TODO(), filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *collection[T]) UpdateOne(filter, update interface{}) error {
	result, err := c.mongoCollection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (c *collection[T]) InsertOne(document interface{}) (primitive.ObjectID, error) {
	result, err := c.mongoCollection.InsertOne(context.TODO(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (c *collection[T]) Upsert(filter, update interface{}) (primitive.ObjectID, error) {
	result, err := c.mongoCollection.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}}, &options.UpdateOptions{Upsert: util.ToPointer(true)})
	if result.UpsertedCount > 0 {
		return result.UpsertedID.(primitive.ObjectID), err
	}
	return primitive.NilObjectID, err
}
