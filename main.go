package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
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
		experHandler   = api.NewExperHandler()
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
	auth.Get("/periods", experHandler.HandlePeriod)
	auth.Get("/pdf", experHandler.GeneratePDF)
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")

	fileNmae :="logs/" + time.Now().Format("02.01.2006_15.04") + ".log"
	logFile, err := os.OpenFile(fileNmae, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening log file: %v\n", err)

	}
	log.SetOutput(logFile)
	
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
