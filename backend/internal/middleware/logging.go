package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} ${status} ${latency} ${ip} ${error}\n",
		Output: log.Writer(),
	})
}

// ErrorMiddleware handles errors and returns JSON responses
func ErrorMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}
