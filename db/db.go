package db

const (
	DBNAME     = "hostel-reservation"
	TestDBName = "hostel-reservation-test"
	DBURI      = "mongodb://localhost:27017"
)

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hostel  HostelStore
	Room    RoomStore
	Booking BookingStore
}
