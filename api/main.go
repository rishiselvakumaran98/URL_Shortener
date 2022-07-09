package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rishiselvakumaran98/URL_Shortener/api/routes"
)

func setupRoutes(app *fiber.app){
	app.Get("/:url", routes.resolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main(){
	err := godotenv.Load()
	if err != nil{
		fmt.Println(err)
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}