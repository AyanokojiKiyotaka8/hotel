package api

import (
	"errors"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	filter := bson.M{}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("booking")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	filter := bson.M{"_id": oid}
	booking, err := h.store.Booking.GetBooking(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error:": "not found"})
		}
		return ErrResourceNotFound("booking")
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"canceled": true,
		},
	}

	booking, err := h.store.Booking.GetBooking(c.Context(), filter)
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	err = h.store.Booking.UpdateBooking(c.Context(), filter, update)
	if err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "canceled",
	})
}
