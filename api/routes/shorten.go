package routes

import (
	"os"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/go-redis/redis/v8"
	"github.com/rishiselvakumaran98/URL_Shortener/api/database"
	"github.com/rishiselvakumaran98/URL_Shortener/api/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

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
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil{
		_= r2.Set(database.Ctx, c.IP, os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else{
		val, _ := r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP().Result())
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"rate_limit_rest": limit/time.Nanosecond/time.Minute,
			})
		}
	}
	// Check if the input is an actual URL 
	if !govalidator.IsURL(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}
	// Check for domain error
	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON()
	}
	// enforce https, SSL, 
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == ""{
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()
	val, err = r.Get(database.Ctx, id).Result()

	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}


	err = r.Set(database.Ctx, id, body.URL, body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect to server",
		})
	}

	resp := response{
		URL: 			body.URL,
		CustomShort:	"",
		Expiry:			body.Expiry,
		RateRemaining:	10,
		RateLimitReset:	30,
	}
	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.RateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.RateLimitReset = ttl/ time.Nanosecond/time.Minute
	resp.CustomShort = os.Getenv("Domain") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}