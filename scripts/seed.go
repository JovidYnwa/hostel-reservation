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

func main(){
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hostelSotre := db.NewMongoHotelStore(client, db.DBNAME)
	hostel := types.Hostel{
		Name: "Serena",
		Location: "Dushanbe",
	}
	room := types.Room{
		Type: types.SingleRoomType,
		BasePrice: 99.9,
	}
	_ = room
	insertedHostel, err := hostelSotre.InsertHostel(ctx, &hostel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHostel)
}