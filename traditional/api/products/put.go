package products

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andreykaipov/target-case-study/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var myRetailCreds string

func init() {
	myRetailCreds = os.Getenv("MYRETAIL_CREDS")

	if myRetailCreds == "" {
		log.Println("MYRETAIL_CREDS is unset. PUT /products/:id will be unauthenticated!")
	}
}

// Put returns our Fiber handlers for our PUT endpoint.
func Put() []fiber.Handler {
	return []fiber.Handler{
		handleAuth(),
		writePriceInfo,
	}
}

func handleAuth() fiber.Handler {
	if myRetailCreds == "" {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	creds := strings.Split(myRetailCreds, ":")

	return basicauth.New(basicauth.Config{
		Users: map[string]string{creds[0]: creds[1]},
	})
}

func writePriceInfo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	// Always expect JSON. Request would fail anyways if it's not JSON; it
	// just makes error handling easier for us.
	c.Request().Header.Add("Content-Type", "application/json")

	body := &ProductRequest{}
	if err := c.BodyParser(body); err != nil {
		return err
	}

	if err := validator.New().Struct(body); err != nil {
		return err
	}

	// This makes the ID optional in the body; we default to the path
	if body.ID == 0 {
		body.ID = id
	}

	if id != body.ID {
		return c.Status(422).SendString("Mismatched IDs in path and in PUT body")
	}

	products := mongo.Client.Database("myretail").Collection("products")

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": body}
	if _, err := products.UpdateOne(context.Background(), filter, update, opts); err != nil {
		return err
	}

	return c.Status(204).JSON(body)
}
