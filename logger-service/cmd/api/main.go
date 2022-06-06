package main

import (
	"context"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	//connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context to disconnect from mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//start web server
	//go app.serve()
	log.Println("Web server started on port", webPort)
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr:    ":" + webPort,
// 		Handler: app.routes(),
// 	}
// 	if err := srv.ListenAndServe(); err != nil {
// 		log.Panic(err)
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
	//create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	con, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo:", err)
		return nil, err
	}

	log.Println("Connected to mongo")

	return con, nil
}
