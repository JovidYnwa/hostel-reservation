package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JovidYnwa/hostel-reservation/api/middleware"
	"github.com/JovidYnwa/hostel-reservation/db/fixtures"
	"github.com/JovidYnwa/hostel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "jova", "edu", false)
		hostel         = fixtures.AddHostel(db.Store, "river hostel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 3.3, hostel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 reponse got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking god %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("excpected %s got %s", booking.UserID, have.UserID)
	}

	//testing non admin can access bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("non 200 reponse got %d", resp.StatusCode)
	}

}

//14
