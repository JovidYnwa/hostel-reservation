package api

import (
	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HostelHandler struct {
	store *db.Store
}

func NewHostelHandler(store *db.Store) *HostelHandler {
	return &HostelHandler{
		store: store,
	}
}

type HostelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HostelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := bson.M{"hostelId": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HostelHandler) HandleGetHostel(c *fiber.Ctx) error {
	id := c.Params("id")
	hostel, err := h.store.Hostel.GetHostelByID(c.Context(), id)
	if err != nil {
		return ErrInvalidID()
	}
	return c.JSON(hostel)
}

type ResourceResp struct {
	Result int `json:"results"`
	Data   any `json:"data"`
	Page   int `json:"page"`
}

func (h *HostelHandler) HandleGetHostels(c *fiber.Ctx) error {
	var pagination db.Pagination
	if err := c.QueryParser(&pagination); err != nil {
		return ErrBadRequest()
	}
	hostels, err := h.store.Hostel.GetHostels(c.Context(), nil, &pagination)
	if err != nil {
		return ErrNotResourceNotFound("hostels")
	}
	resp := ResourceResp{
		Data:   hostels,
		Result: len(hostels),
		Page:   int(pagination.Page),
	}
	return c.JSON(resp)
}
