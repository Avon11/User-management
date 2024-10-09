package main

import (
	"log"
	"os"

	"kenshi/config"
	"kenshi/database"
	"kenshi/handlers"
	"kenshi/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.New()

	// Connect to MongoDB
	client, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(nil)

	// Create Fiber app
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	// Routes
	api := app.Group("/api")

	// User routes
	api.Post("/signup", handlers.SignUp)
	api.Post("/signin", handlers.SignIn)
	api.Post("/signout", handlers.SignOut)

	// Task routes (protected)
	tasks := api.Group("/tasks", middleware.AuthRequired())
	tasks.Post("/", handlers.CreateTask)
	tasks.Get("/", handlers.GetAllTasks)
	tasks.Get("/:id", handlers.GetTask)
	tasks.Put("/:id", handlers.UpdateTask)
	tasks.Delete("/:id", handlers.DeleteTask)

	// get token
	token := api.Group("/token", middleware.AuthRequired())
	token.Get("/", handlers.RefreshToken) // This is the API to get access toekn if it is expired. client need to send refresh token in this api to get new access token

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
