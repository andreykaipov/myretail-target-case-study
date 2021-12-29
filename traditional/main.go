package main

import (
	"log"
	"os"

	"github.com/andreykaipov/target-case-study/api/products"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var bind string

func init() {
	if bind = os.Getenv("MYRETAIL_BIND"); bind == "" {
		bind = ":3000"
	}
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Get("/products/:id", products.Get()...)
	app.Put("/products/:id", products.Put()...)
	log.Fatal(app.Listen(bind))
}
