package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "aaa@bbb.com",
		FirstName: "Aaa",
		LastName:  "Bbb",
		Password:  "aaabbb111",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncPassword) > 0 {
		t.Errorf("expecting a EncryptedPassword not to be included in json")
	}
	if params.FirstName != user.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if params.LastName != user.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if params.Email != user.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

func init() {
	db.LoadConfig()
}
