package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-mongo-todos/db"
	"github.com/go-mongo-todos/handlers"
	"github.com/go-mongo-todos/services"
	"github.com/joho/godotenv"
)

type Application struct {
	Models services.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load env file: ", err.Error())
	}
	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		log.Println("failed to connect: ", err.Error())

	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			// panic(err)
			log.Println(err.Error())
		}
	}()

	services.New(mongoClient)

	if err := http.ListenAndServe(":8080", handlers.CreateRouter()); err != nil {
		log.Panicln("servier running error: ", err.Error())
	} else {
		log.Println("Server running in port", 8080)
	}
}
