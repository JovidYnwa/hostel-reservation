package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JovidYnwa/hostel-reservation/api"
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

func seedUser(isAdmin bool, fname string, lname string, email string, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: password,
	})
	if err != nil{
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	inserteUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return inserteUser
}

func seedHostel(name string, location string, rating int) *types.Hostel {
	hostel := types.Hostel{
		Name:     name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}

	insertedHostel, err := hostelSotre.InsertHostel(ctx, &hostel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHostel
}

func seedRoom(size string, ss bool, price float64, hostelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size: size,
		Seaside: ss,
		Price: price,
		HostelID: hostelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func main() {
	seedUser(true, "jova", "admin", "jova@admin.com", "adminpass1234")
	seedUser(false, "vova", "notadmin", "jova@jova.com", "supersecurepass")
	hostel := seedHostel("Serena", "Tajikistan", 5)
	seedRoom("small", true, 99.99, hostel.ID)
	seedRoom("medium", true, 199.99, hostel.ID)
	seedRoom("medium", false, 199.99, hostel.ID)
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
