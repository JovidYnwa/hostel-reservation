package api

import (

	"github.com/gofiber/fiber/v2"
)

type resp struct {
	resp_code int
	resp_text string
}

func ExperHandlerGet(c *fiber.Ctx) error {
	return c.JSON("good")
}