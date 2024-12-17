package db

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

type Pagination struct {
	Page  int64
	Limit int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
