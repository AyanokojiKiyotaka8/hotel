package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn, ln string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(store *db.Store, name, loc string, rooms []primitive.ObjectID, rating int) *types.Hotel {
	roomIDs := rooms
	if roomIDs == nil {
		roomIDs = []primitive.ObjectID{}
	}

	hotel := &types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.Background(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(store *db.Store, size string, price float64, ss bool, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		SeaSide: ss,
		HotelID: hid,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, np int, fd, td time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:     uid,
		RoomID:     rid,
		NumPersons: np,
		FromDate:   fd,
		TillDate:   td,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
