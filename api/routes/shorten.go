package routes

import (
	"time"
	"github.com/gofiber/fiber/v2"
)

// Define the request format for the shortener
type request struct {
	URL				string			`json:"url"`
	CustomShort		string			`json:"short"`
	Expiry			time.Duration	`json:"expiry"`
}

// Define the response format for the shortener 
type response struct {
	URL				string			`json:"url"`
	CustomShort		string			`json:"short"`
	Expiry			time.Duration	`json:"expiry"`
	RateRemaining	int 			`json:"rate_limit"`
	RateLimitReset	time.Duration	`json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(&body); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse Json"})
	} 
	// Implement rate limiting
	if !govalidator.IsURL(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}
	// Check if the input is an actual URL
	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON()
	}
	// Check for domain error

	// enforce https, SSL, 
	body.URL = helpers.EnforceHTTP(body.URL)
}