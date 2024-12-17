package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AyanokojiKiyotaka8/booking.git/api"
	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	bookingStore := db.NewMongoBookingStore(client, db.DBNAME)
	store := &db.Store{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	}

	if err := hotelStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := roomStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := userStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := bookingStore.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	user := fixtures.AddUser(store, "aaa", "bbb", false)
	fmt.Println(user.FirstName, " -> ", api.CreateToken(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println(admin.FirstName, " -> ", api.CreateToken(admin))

	hotel := fixtures.AddHotel(store, "qqq", "www", nil, 7)
	fmt.Println(hotel)

	room := fixtures.AddRoom(store, "large", 100.00, true, hotel.ID)
	fmt.Println(room)

	booking := fixtures.AddBooking(store, user.ID, room.ID, 7, time.Now(), time.Now().AddDate(0, 0, 7))
	fmt.Println(booking)

	for i := 0; i < 100; i++ {
		_ = fixtures.AddHotel(store, fmt.Sprintf("Hotel %d", i), "www", nil, rand.Intn(10))
	}
}

func init() {
	db.LoadConfig()
}
