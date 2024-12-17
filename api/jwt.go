package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnauthorized()
		}

		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}

		exp := claims["expires"].(float64)
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
		if err != nil {
			return err
		}

		filter := bson.M{"_id": userID}
		user, err := userStore.GetUser(c.Context(), filter)
		if err != nil {
			return ErrUnauthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, ErrUnauthorized()
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println("Error in parsing token:", err)
		return nil, ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnauthorized()
	}
	return claims, nil
}
