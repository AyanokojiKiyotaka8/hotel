package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "wertyuio",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)
	fmt.Println(resp.Status)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	fmt.Println(user)

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
