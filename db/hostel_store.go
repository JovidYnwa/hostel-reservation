package db

import (
	"context"

	"github.com/JovidYnwa/hostel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hostel interface {
	InsertHostel(context.Context, *types.Hostel)(*types.Hostel, error)
}

type MongoHostelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHostelStore{
	return &MongoHostelStore{
		client: client,
		coll: client.Database(dbname).Collection("hostel"),
	}
}


func (s *MongoHostelStore) InsertHostel(ctx context.Context, hostel *types.Hostel) (*types.Hostel, error) {
	resp, err := s.coll.InsertOne(ctx, hostel)
	if err != nil {
		return nil, err
	}
	hostel.ID = resp.InsertedID.(primitive.ObjectID)
	return hostel, nil
}
