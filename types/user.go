package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswoerdLen = 7

)

type CreateUserParams struct {
	FirstName 	string `json:"firstName"`
	LastName 	string `json:"lastName"`
	Email    	string `json:"email"`
	Password    string `json:"password"`
}

func (params CreateUserParams) Validate() error {
	if len(params.FirstName) < minFirstNameLen {
		return fmt.Errorf("firstNmae length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		return fmt.Errorf("lastNmae length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswoerdLen {
		return fmt.Errorf("password length should be at least %d characters", minPasswoerdLen)
	}
	if !isEmailValid(params.Email) {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

type User struct{
	ID       		  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName 		  string             `bson:"firstName" json:"firstName"`
	LastName 		  string 			 `bson:"lastName" json:"email"`
	Email    		  string 			 `bson:"email" json:"lastName"`
	EncryptedPassword string 		     `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: 			params.FirstName,
		LastName:  			params.LastName,
		Email:     			params.Email,
		EncryptedPassword:  string(encpw),
	},nil
}