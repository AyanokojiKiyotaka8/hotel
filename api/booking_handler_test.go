package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AyanokojiKiyotaka8/booking.git/db/fixtures"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(tdb.Store, "qqq", "www", false)
		user           = fixtures.AddUser(tdb.Store, "aaa", "bbb", false)
		hotel          = fixtures.AddHotel(tdb.Store, "qqq", "www", nil, 7)
		room           = fixtures.AddRoom(tdb.Store, "large", 100.00, true, hotel.ID)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, 7, time.Now(), time.Now().AddDate(0, 0, 7))
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		route          = app.Group("/", JWTAuthentication(tdb.Store.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateToken(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d", resp.StatusCode)
	}

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}

	// test non auth user cannot access the booking
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateToken(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 response got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser      = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.Store, "aaa", "bbb", false)
		hotel          = fixtures.AddHotel(tdb.Store, "qqq", "www", nil, 7)
		room           = fixtures.AddRoom(tdb.Store, "large", 100.00, true, hotel.ID)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, 7, time.Now(), time.Now().AddDate(0, 0, 7))
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(tdb.Store.User), AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateToken(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	if bookings[0].ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookings[0].ID)
	}
	if bookings[0].UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookings[0].UserID)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateToken(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected a status unauthorized but got %d", resp.StatusCode)
	}
}
