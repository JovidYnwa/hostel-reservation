package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		err           error
		mongoEndpoint = os.Getenv(db.MongoDBNameEnvName)
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoEndpoint).Drop(ctx); err != nil {
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

	for i := 0; i <= 100; i++ {
		name := fmt.Sprintf("hostel name %d", i)
		location := fmt.Sprintf("hostel location %d", i)
		fixtures.AddHostel(store, name, location, rand.Intn(5)+1, nil)

	}
}
