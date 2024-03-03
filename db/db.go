package db

import "os"

var DBNAME string

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

func init() {
	DBNAME = os.Getenv("MONGO_DB_NAME")
}
