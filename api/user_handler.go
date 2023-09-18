package api

import (
	"context"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/types"
	"github.com/gofiber/fiber/v2"
)


type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler{
	return &UserHandler{
		userStore: userStore,
	}
}


func (h *UserHandler) HandleGetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil{
		return err
	}
	return c.JSON(user)
}

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FistName: "James",
		LastName: "Yo",
	}
	return c.JSON(u)
}