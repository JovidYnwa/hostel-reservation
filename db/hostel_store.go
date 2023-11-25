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
	GetHostels(context.Context, bson.M) ([]*types.Hostel, error)
}

type MongoHostelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

// Insert implements HostelStore.
func (*MongoHostelStore) Insert(context.Context, *types.Hostel) (*types.Hostel, error) {
	panic("unimplemented")
}

func NewMongoHostelStore(client *mongo.Client) *MongoHostelStore {
	return &MongoHostelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hostels"),
	}
}

func (s *MongoHostelStore) GetHostels(ctx context.Context, filter bson.M) ([]*types.Hostel, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hostels []*types.Hostel
	if err :=resp.All(ctx, &hostels); err != nil{
		return nil, err
	}
	return hostels, nil
}

func (s *MongoHostelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
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
