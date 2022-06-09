package ports

import "github.com/gofiber/fiber/v2"

func Index(c *fiber.Ctx) error {
	const helloMessage = "Hello World!"
	c.Send([]byte(helloMessage))
	return nil
}
