package api

import (
	"errors"
	"fmt"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}	
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserById(c.Context(), id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.JSON(map[string]string{"error": "not found"})
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		//values bson.M
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	filter := db.Map{"_id": userID}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return ErrNotResourceNotFound("user")
	}
	return c.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandler) HandlerTest(c *fiber.Ctx) error {
	return fmt.Errorf("bad stuff going on")
	// user := c.Context().UserValue("user")
	// return c.JSON(user)
}
