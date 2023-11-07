package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hostelSotre := db.NewMongoHostelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hostel := types.Hostel{
		Name:     "Serena",
		Location: "Dushanbe",
	}
	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 129.9,
		},
	}
	insertedHostel, err := hostelSotre.InsertHostel(ctx, &hostel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room :=range rooms{
	room.HostelID = insertedHostel.ID
	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedRoom)
  }
}
