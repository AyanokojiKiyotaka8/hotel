package main

import (
	"context"
	"log"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	err        error
	hotelStore *db.MongoHotelStore
	roomStore  *db.MongoRoomStore
	userStore  *db.MongoUserStore
)

func init() {
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)

	if err := hotelStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := roomStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := userStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "kingsize",
			Price: 199.9,
		},
		{
			Size:  "normal",
			Price: 122.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err = roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(isAdmin bool, fname, lname, email, password string) {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})

	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("Aaa", "Bbb", 1)
	seedHotel("Qqq", "Www", 2)
	seedHotel("Eee", "Rrr", 3)
	seedUser(false, "aaa", "bbb", "aaa@bbb.com", "aaabbb111")
	seedUser(true, "admin", "admin", "admin@admin.com", "admin111")
}
