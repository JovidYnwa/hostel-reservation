package api

import (
	"net/http"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrNotResourceNotFound("booking")
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{Type: "msg", Msg: "updated"})
}

// This should be admin authorized! TODO
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrNotResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

// This need to be user authorized! TODO
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		log.Info(err)
		return ErrNotResourceNotFound("booking")
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	return c.JSON(booking)
}
