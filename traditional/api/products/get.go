package products

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/andreykaipov/target-case-study/api/redsky"
	"github.com/andreykaipov/target-case-study/mongo"
	"github.com/andreykaipov/target-case-study/redis"
	redisdriver "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// Get returns our Fiber handlers for our GET endpoint in the necessary order
// for whether Redis is enabled or not.
func Get() []fiber.Handler {
	noRedis := func(_ *fiber.Ctx) bool { return redis.URI == "" }

	return []fiber.Handler{
		prepProduct,
		skip.New(checkCache, noRedis),
		askRedSky,
		skip.New(writeToCache, noRedis),
		queryPriceInfo,
		end,
	}
}

func prepProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	c.Locals("product", &ProductResponse{ID: id})

	return c.Next()
}

func checkCache(c *fiber.Ctx) error {
	log.Println("Checking cache to see if product info already exists...")

	product := c.Locals("product").(*ProductResponse)

	cmd := redis.Client.HGet(context.Background(), c.Params("id"), "name")

	var err error
	if product.Name, err = cmd.Result(); err == redisdriver.Nil {
		return c.Next()
	} else if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return queryPriceInfo(c)
}

func askRedSky(c *fiber.Ctx) error {
	log.Println("Hitting RedSky for general product info...")

	product := c.Locals("product").(*ProductResponse)

	resp, err := redsky.Client.GetProductInfoSmall(c.Params("id"))
	if err != nil {
		return err
	}

	product.Name = resp.Data.Product.Item.ProductDescription.Title

	return c.Next()
}

func writeToCache(c *fiber.Ctx) error {
	log.Println("Writing product info to cache...")

	id := c.Params("id")
	product := c.Locals("product").(*ProductResponse)

	if _, err := redis.Client.HSet(context.Background(), id, "name", product.Name).Result(); err != nil {
		return err
	}

	if _, err := redis.Client.Expire(context.Background(), id, 48*time.Hour).Result(); err != nil {
		return err
	}

	return c.Next()
}

func queryPriceInfo(c *fiber.Ctx) error {
	log.Println("Querying database for product price info...")

	product := c.Locals("product").(*ProductResponse)

	products := mongo.Client.Database("myretail").Collection("products")
	result := products.FindOne(context.Background(), bson.M{"_id": product.ID})

	if err := result.Decode(product); err == mongodriver.ErrNoDocuments {
		log.Printf("Pricing info doesn't exist for item %d yet :(", product.ID)
		return end(c)
	} else if err != nil {
		return err
	}

	return end(c)
}

func end(c *fiber.Ctx) error {
	log.Println("Returning product response...")

	return c.JSON(c.Locals("product"))
}
