package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, from , till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID: uid,
		RoomID: rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking); 
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}


func AddRoom(store *db.Store, size string, ss bool, price float64, hid primitive.ObjectID) *types.Room{
	room := &types.Room{
		Size: size,
		Seaside: ss,
		Price: price,
		HostelID: hid,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHostel(store *db.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hostel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hostel := types.Hostel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHostel, err := store.Hostel.InsertHostel(context.TODO(), &hostel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHostel
}

func AddUser(store *db.Store, fn string, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	inserteUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return inserteUser
}
