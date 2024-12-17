package api

import (
	"context"
	"log"
	"testing"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	Store  *db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TestDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.TestDBNAME)
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, db.TestDBNAME, hotelStore),
			User:    db.NewMongoUserStore(client, db.TestDBNAME),
			Booking: db.NewMongoBookingStore(client, db.TestDBNAME),
		},
	}
}

func init() {
	db.LoadConfig()
}
