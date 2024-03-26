package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Customizing error handling
var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
	AppName:      "Hostel-Resevation",
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println(client)
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
	fmt.Println(userStore)
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

	app.Use(requestid.New())

	// Logging middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} - ${ip} - ${ua} - ${method} ${path} - ${status} ${resBody}\n", // Customize logging format
		TimeFormat: "02/Jan/2006:15:04:05 -0700",                                           // Customize time format
		Output:     os.Stdout,                                                              // Log output to stdout
	}))
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}


