package db

import (
	"context"

	"github.com/JovidYnwa/hostel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HostelStore interface {
	InsertHostel(context.Context, *types.Hostel) (*types.Hostel, error)
	Update(context.Context, Map, Map) error
	GetHostels(context.Context, Map, *Pagination) ([]*types.Hostel, error)
	GetHostelByID(context.Context, string) (*types.Hostel, error)
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

func (s *MongoHostelStore) GetHostelByID(ctx context.Context, id string) (*types.Hostel, error) {
	var hostel types.Hostel
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hostel); err != nil {
		return nil, err
	}
	return &hostel, nil
}

func (s *MongoHostelStore) GetHostels(ctx context.Context, filter Map, pag *Pagination) ([]*types.Hostel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pag.Page - 1) * pag.Limit)
	opts.SetLimit(pag.Limit)
	resp, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hostels []*types.Hostel
	if err := resp.All(ctx, &hostels); err != nil {
		return nil, err
	}
	return hostels, nil
}

func (s *MongoHostelStore) Update(ctx context.Context, filter Map, update Map) error {
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
