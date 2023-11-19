package db

import (
	"context"

	"github.com/JovidYnwa/hostel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HostelStore interface {
	InsertHostel(context.Context, *types.Hostel) (*types.Hostel, error)
	Update(context.Context, bson.M, bson.M) error
}

type MongoHostelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHostelStore(client *mongo.Client) *MongoHostelStore {
	return &MongoHostelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hostels"),
	}
}

func (s *MongoHostelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx,filter, update)
	return err
}

func (s *MongoHostelStore) InsertHostel(ctx context.Context, hostel *types.Hostel) (*types.Hostel, error) {
	resp, err := s.coll.InsertOne(ctx, hostel)
	if err != nil {
		return nil, err
	}
	hostel.ID = resp.InsertedID.(primitive.ObjectID)
	return hostel, nil
}
