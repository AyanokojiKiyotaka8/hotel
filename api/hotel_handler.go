package api

import (
	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type ResultResp struct {
	Data    any   `json:"data"`
	Results int   `json:"results"`
	Page    int64 `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := bson.M{"rating": params.Rating}
	opts := &options.FindOptions{}
	opts.SetSkip((params.Page - 1) * params.Limit)
	opts.SetLimit(params.Limit)

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, opts)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}

	res := ResultResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    params.Page,
	}
	return c.JSON(res)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	filter := bson.M{"_id": oid}
	hotel, err := h.store.Hotel.GetHotel(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("room")
	}
	return c.JSON(rooms)
}
