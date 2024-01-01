package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/JovidYnwa/hostel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
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
	var params AuthParams
	if err := c.BodyParser(&params); err != nil{ 
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")		
		}
		return err
	}	
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("Invalid credantials")
	}
	return nil
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * 4)
	claims := jwt.MapClaims{
		"id": user.ID,
		"email": user.Email,
		"validTill": validTill,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		fmt.Println(err)
	}
	return tokenStr
}