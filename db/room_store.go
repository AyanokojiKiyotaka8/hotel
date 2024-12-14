package db

import (
	"context"
	"fmt"

	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, dbname string, hs HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection("rooms"),
		HotelStore: hs,
	}
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("---- dropping the room collection ----")
	return s.coll.Drop(ctx)
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": room.HotelID}
	update := bson.M{
		"$push": bson.M{"rooms": room.ID},
	}
	if err := s.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := curr.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
