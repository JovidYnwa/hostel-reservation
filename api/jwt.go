package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/JovidYnwa/hostel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("JWT Authenticating")
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnauthorized()
		}
		claims, err := validateToken(token)
		fmt.Println(claims)
		if err != nil {
			return err
		}
		fmt.Println(claims)
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		//Check token expiration
		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserById(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}
		//Set the current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method ", token.Header["alg"])
			return nil, fmt.Errorf("unauthorize")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to pars JWT token: ", err)
		return nil, ErrUnauthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrUnauthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil,ErrUnauthorized()
	}
	return claims, nil
}
