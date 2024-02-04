package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/db/fixtures"
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
	bookingStore db.BookingStore
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

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		FromDate: from,
		TillDate: till,
	}
	resp, err := bookingStore.InsertBooking(context.Background(), booking); 
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking:", resp.ID)
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
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err :=client.Database(db.DBNAME).Drop(ctx); err != nil{
		log.Fatal(err)
	}
	hostelSotre := db.NewMongoHostelStore(client)
	store := &db.Store {
		User: db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room: db.NewMongoRoomStore(client, hostelSotre),
		Hostel: hostelSotre,
	}
	user := fixtures.AddUser(store, "test", "testi", false)
	admin := fixtures.AddUser(store, "admin", "testi", true)
	fmt.Println("admin ->", admin)
	hostel := fixtures.AddHostel(store, "some hostel", "casablanka", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hostel.ID)
	booking := fixtures.AddBooking(*store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println(booking)	
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
	bookingStore = db.NewMongoBookingStore(client)

}
