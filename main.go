package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Estructura y lista de libros (igual que antes)
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

var books = []Book{
	{ID: "1", Title: "1984", Author: "George Orwell", ISBN: "0451524934"},
	{ID: "2", Title: "The Catcher in the Rye", Author: "J.D. Salinger", ISBN: "0316769487"},
}

// Simulación de un proceso costoso
func generateRecommendations(books []Book, ch chan<- []string) {
	// Simular un tiempo de procesamiento largo
	time.Sleep(10 * time.Second)

	recommendations := []string{}
	for _, book := range books {
		recommendations = append(recommendations, fmt.Sprintf("Recommended: %s by %s", book.Title, book.Author))
	}

	ch <- recommendations
}

func main() {
	app := fiber.New()

	// Endpoint para generar recomendaciones
	app.Get("/books/recommend", func(c *fiber.Ctx) error {
		recommendationsCh := make(chan []string)

		// Lanzar la goroutine para procesar las recomendaciones
		go generateRecommendations(books, recommendationsCh)

		log.Println("Working on other stuff :)")
		// Realizar otras tareas si es necesario
		// Aquí podrías manejar otras solicitudes, logs, etc.

		// Esperar a que la goroutine termine y obtener el resultado
		recommendations := <-recommendationsCh

		return c.JSON(recommendations)
	})

	// Otros endpoints como antes...
	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(books)
	})

	app.Post("/books", func(c *fiber.Ctx) error {
		var book Book
		if err := c.BodyParser(&book); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		books = append(books, book)
		return c.JSON(book)
	})

	app.Put("/books/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var updatedBook Book
		if err := c.BodyParser(&updatedBook); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		for i, book := range books {
			if book.ID == id {
				books[i] = updatedBook
				books[i].ID = id
				return c.JSON(books[i])
			}
		}
		return c.Status(fiber.StatusNotFound).SendString("Book not found")
	})

	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, book := range books {
			if book.ID == id {
				books = append(books[:i], books[i+1:]...)
				return c.SendString("Book deleted")
			}
		}
		return c.Status(fiber.StatusNotFound).SendString("Book not found")
	})

	app.Listen(":8000")
}
