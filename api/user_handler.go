package api

import (
	"github.com/JovidYnwa/hostel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FistName: "James",
		LastName: "Yo",
	}
	return c.JSON(u)
}

func HandleGetUserByID(c *fiber.Ctx) error {
	return c.JSON("James")
}
