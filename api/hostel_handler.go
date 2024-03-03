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

type HostelQueryParams struct {
	db.Pagination
	Rating int 
}

func (h *HostelHandler) HandleGetHostels(c *fiber.Ctx) error {
	var params HostelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": params.Rating,
	}
	hostels, err := h.store.Hostel.GetHostels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrNotResourceNotFound("hostels")
	}
	resp := ResourceResp{
		Data:   hostels,
		Result: len(hostels),
		Page:   int(params.Page),
	}
	return c.JSON(resp)
}