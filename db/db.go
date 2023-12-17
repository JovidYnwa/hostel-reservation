package db


const (
	DBNAME = "hostel-reservation"
	TestDBName = "hostel-reservation-test"
	DBURI = "mongodb://localhost:27017"
)


type Store struct {
	User    UserStore
	Hostel  HostelStore
	Room    RoomStore
}