package main

import (
	"context"
	"log"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	roomStore db.RoomStore
	hostelSotre db.HostelStore
	userStore db.UserStore
	ctx = context.Background()
)

func seedUser(fname string, lname string, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: "securepasswordhaha",
	})
	if err != nil{
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil{
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {
	hostel := types.Hostel{
		Name:     name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
			Size:      "small",
			Price: 99.9,
		},
		{
			Size:      "normal",
			Price: 199.9,
		},
		{
			Size:      "kingsize",
			Price: 129.9,
		},
	}
	insertedHostel, err := hostelSotre.InsertHostel(ctx, &hostel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room :=range rooms{
	room.HostelID = insertedHostel.ID
	_, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}
  }
}

func main() {
	seedHotel("Serena", "Tajikistan", 5)
	seedHotel("Moscwa", "Russia", 1)
	seedUser("jova", "jo", "jova@jova.com")
}

func init(){
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err :=client.Database(db.DBNAME).Drop(ctx); err != nil{
		log.Fatal(err)
	}
	hostelSotre = db.NewMongoHostelStore(client)
	roomStore = db.NewMongoRoomStore(client, hostelSotre)
	userStore = db.NewMongoUserStore(client)

}
