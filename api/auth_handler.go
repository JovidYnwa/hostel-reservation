package api

import (
	"fmt"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler{
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email string `json:"email"`
	Password string `json:"password"`

}

func (h *AuthHandler) HandleAuthenticate (c *fiber.Ctx) error {
	var AuthParams AuthParams
	if err := c.BodyParser(&AuthParams); err != nil{ 
		return err
	}
	fmt.Println(AuthParams)
	return nil
}