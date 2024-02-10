package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hostelSotre  db.HostelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func main() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hostelSotre := db.NewMongoHostelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hostelSotre),
		Hostel:  hostelSotre,
	}
	user := fixtures.AddUser(store, "test", "testi", false)
	fmt.Println("test ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "testi", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hostel := fixtures.AddHostel(store, "some hostel", "casablanka", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hostel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println(booking)
}


