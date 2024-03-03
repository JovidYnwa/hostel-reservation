package main

import (
	"context"
	"log"
	"os"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Customizing error handling
var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBNAME))
	if err != nil {
		log.Fatal(err)
	}

	//Handler initialization
	var (
		hostelStore  = db.NewMongoHostelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hostelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hostel:  hostelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hostelHandeler = api.NewHostelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore)) //cheking the authentication
		admin          = apiv1.Group("/admin", api.AdminAuth)
	)

	auth.Post("/auth", authHandler.HandleAuthenticate)

	//User handlers
	auth.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)

	//hostel handlers
	apiv1.Get("/hostel", hostelHandeler.HandleGetHostels)

	//room handlers
	apiv1.Get("room/", roomHandler.HandleGetRooms)
	apiv1.Post("room/:id/book", roomHandler.HandleBookRoom)

	//booking handlers
	apiv1.Get("booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("booking/:id/cancel", bookingHandler.HandleGetCancelBooking)

	//admin can do only
	admin.Get("booking", bookingHandler.HandleGetBookings)

	//cancel booking TODO

	//For testing
	apiv1.Get("/test", userHandler.HandlerTest)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}
