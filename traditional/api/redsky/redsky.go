package redsky

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type client struct {
	key string
	api string
}

// Client represents a RedSky client.
var Client *client

func init() {
	env := os.Getenv("REDSKY_ENV")
	if env == "" {
		log.Fatal("REDSKY_ENV is required")
	}

	client, err := findRedSky(env)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
}

// Note the following keys are NOT secret, so relax.
func findRedSky(env string) (*client, error) {
	var key, api string

	switch env {
	case "uat":
		api = "https://redsky-uat.perf.target.com"
		key = "3yUxt7WltYG7MFKPp7uyELi1K40ad2ys"
	case "prod":
		api = "https://redsky.target.com"
		key = "ff457966e64d5e877fdbad070f276d18ecec4a01"
	default:
		return nil, fmt.Errorf("Unrecognized RedSky environment: %q", env)
	}

	return &client{key, api}, nil
}

// GetProductInfoSmall hits the stripped down "case study" product API endpoint.
// This endpoint doesn't include pricing information.
func (c *client) GetProductInfoSmall(id string) (*Product, error) {
	return c.get("/redsky_aggregations/v1/redsky/case_study_v1", map[string]string{"tcin": id})
}

func (c *client) get(endpoint string, m map[string]string) (*Product, error) {
	url := fmt.Sprintf("%s%s", c.api, endpoint)

	args := fiber.AcquireArgs()
	defer fiber.ReleaseArgs(args)
	for k, v := range m {
		args.Set(k, v)
	}
	args.Set("key", c.key)

	req := fiber.Get(url).UserAgent("myRetail").QueryStringBytes(args.QueryString())

	code, b, errs := req.Bytes()
	if len(errs) > 0 {
		return nil, fmt.Errorf("%s", errs)
	}

	if code/100 != 2 {
		return nil, fmt.Errorf("Non-2xx status code from RedSky %d: %s", code, b)
	}

	product := &Product{}
	if err := json.Unmarshal(b, product); err != nil {
		return nil, err
	}

	return product, nil
}
