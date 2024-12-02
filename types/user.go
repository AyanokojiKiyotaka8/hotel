package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minFirstNameLen = 2
	minLastNameLen = 2
	minPasswordLen = 7
)

type User struct {
	ID 			primitive.ObjectID 	`bson:"_id,omitempty" json:"id,omitempty"`
	FirstName 	string 				`bson:"firstName" json:"firstName"`
	LastName 	string 				`bson:"lastName" json:"lastName"`
	Email 		string 				`bson:"email" json:"email"`
	EncPassword string				`bson:"EncPassword" json:"-"`
}

type GetUserFromParams struct {
	FirstName 	string 	`json:"firstName"`
	LastName 	string 	`json:"lastName"`
	Email 		string 	`json:"email"`
	Password	string	`json:"password"`
}

func NewUserFromParams(u *GetUserFromParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName: u.FirstName,
		LastName: u.LastName,
		Email: u.Email,
		EncPassword: string(encpw),
	}, nil
}

func (u *GetUserFromParams) Validate() []string {
	var errors []string 
	if len(u.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen))
	}
	if len(u.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen))
	}
	if len(u.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
	}
	if !isEmailValid(u.Email) {
		errors = append(errors, fmt.Sprintf("email is invalid"))
	}

	return errors
}

func isEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}