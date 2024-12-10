package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type AuthResponse struct {
	User 	*types.User `json:"user"`
	Token 	string		`json:"token"`
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{"email": params.Email}
	user, err := h.userStore.GetUser(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("wrong credentials")
		}
		return err
	}

	if !types.IsValidPassword(user.EncPassword, params.Password) {
		return fmt.Errorf("wrong credentials")
	}

	resp := AuthResponse{
		User: user,
		Token: createToken(user),
	}

	return c.JSON(resp)
}

func createToken(user *types.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"email": user.Email,
		"expires": time.Now().Add(3 * time.Hour).Unix(),
	})
	
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
	}
	return tokenStr
}