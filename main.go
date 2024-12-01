package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello World")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("Request Method:", c.Method())
		fmt.Println("Request URL:", c.OriginalURL())
		fmt.Println("Request Body:", string(c.Body()))
		return c.Next()
	})

	todos := []Todo{}

	// get Todo list
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// create Todo item
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is empty"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		// var x int = 5   //0x00001
		// var p *int = &x //0x00001
		// fmt.Println(p)  //0x00001
		// fmt.Println(*p) //5

		return c.Status(201).JSON(todo)
	})

	// update Todo item
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true

				return c.Status(200).JSON(todo)
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// delete Todo item
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)

				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":" + PORT))

}
