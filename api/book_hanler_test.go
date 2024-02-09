package api

import (
	"encoding/json"
	"fmt"
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
		user           = fixtures.AddUser(db.Store, "jova", "edu", false)
		hostel         = fixtures.AddHostel(db.Store, "river hostel", "a", 4, nil)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(*db.Store, user.ID, hostel.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		BookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking

	admin.Get("/", BookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookings []types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookings)
}


//6:03