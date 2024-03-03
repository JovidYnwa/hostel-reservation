package db

const MongoDBNameEnvName = "Mongo_DB_NAME"

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
