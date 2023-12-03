package main

import (
	"context"
	"flag"
	"log"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/api/middleware"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Customizing error handling
var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error": err.Error()})
		    
    },
}


func main(){

	listenAddr :=flag.String("listenAddr", ":5000", "The listen address of the Api server")
	flag.Parse()


	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	//Handler initialization
	var(
	userHandler    = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME)) 
	hostelStore    = db.NewMongoHostelStore(client)
	roomStore      = db.NewMongoRoomStore(client, hostelStore)
	hostelHandeler = api.NewHostelHandler(hostelStore, roomStore)
	app   = fiber.New(config)
	auth  = app.Group("/api")
	apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)//cheking the authentication
	)
	
	auth.Post("/auth", userHandler.HandleAuthenticate)
	
	//User handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)

	//hostel handlers
	apiv1.Get("/hostel", hostelHandeler.HandleGetHostels)
	app.Listen(*listenAddr)

	//For testing
	apiv1.Get("/test", userHandler.HandlerTest)
	app.Listen(*listenAddr)
}
//27