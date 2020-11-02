package database

import (
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	
	"../config"
)

// PlacesColl is placese Collection
var  PlacesColl *mongo.Collection
// UserColl is users Collection
var  UserColl   *mongo.Collection
//Client .
var Client *mongo.Client
//ClientDisconnect is disconnect db
var ClientDisconnect func ()
// ConCancel .
var ConCancel func ()


//DbConnect connect mongodb 
func DbConnect () {
	password := config.MongoPassword
	var URL = "mongodb+srv://hikaru:" + password + "@cluster0.zla4m.mongodb.net/placesShare?retryWrites=true&w=majority"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	ConCancel = cancel
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URL))
	if err != nil {
		panic(err)
	}
	Client = client
	PlacesColl = client.Database("placesShare").Collection("places")
	UserColl = client.Database("placesShare").Collection("users")

  ClientDisconnect = func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
  
}