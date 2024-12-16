package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/AyanokojiKiyotaka8/booking.git/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "aaa", "bbb", false)
	hotel := fixtures.AddHotel(tdb.Store, "qqq", "www", nil, 7)
	room := fixtures.AddRoom(tdb.Store, "large", 100.00, true, hotel.ID)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, 7, time.Now(), time.Now().AddDate(0, 0, 7))
	fmt.Println(booking)
}
