package api

import (
	"fmt"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
)


type HostelHandler struct {
	hostelStore db.HostelStore 
	roomStore   db.RoomStore
}

func NewHostelHandler(hs db.HostelStore, rs db.RoomStore) *HostelHandler{
	return &HostelHandler{
		hostelStore: hs,
		roomStore: rs, 
	} 
}

type HostelQueryParams struct {
	Rooms bool
	Rating int
}

func (h *HostelHandler) HandleGetHostels(c *fiber.Ctx) error {
	var qparams HostelQueryParams
	if err := c.QueryParser(&qparams); err !=nil {
		return err
	}
	fmt.Println(qparams)

	hostels, err := h.hostelStore.GetHostels(c.Context(),nil)
	if err != nil {
		return err
	}
	return c.JSON(hostels)
}