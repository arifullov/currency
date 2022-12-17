package main

import (
	"github.com/gofiber/fiber/v2"
)

func healthcheckHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"status": "available",
	})
}
