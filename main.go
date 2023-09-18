package main

import (
	"context"
	"flag"
	"log"

	"github.com/JovidYnwa/hostel-reservation/api"
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hostel-reservation"
const userColl = "users"


// Customizing error handling
var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error": err.Error()})
		    
    },
}


func main(){

	listenAddr :=flag.String("listenAddr", ":5000", "The listen address of the Api server")
	flag.Parse()


	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	//Handler initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	app.Listen(*listenAddr)
}
