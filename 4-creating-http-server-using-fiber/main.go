package main

import (
	"errors"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Book represents a book in our library
type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Genre         string `json:"genre"`
	Pages         int    `json:"pages"`
	Available     bool   `json:"available"`
	PublishedYear int    `json:"published_year"`
}

// Our in-memory book storage
var library = []Book{
	{ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan", Genre: "Technology", Pages: 380, Available: true, PublishedYear: 2015},
	{ID: 2, Title: "Clean Code", Author: "Robert Martin", Genre: "Technology", Pages: 464, Available: false, PublishedYear: 2008},
	{ID: 3, Title: "The Pragmatic Programmer", Author: "Dave Thomas", Genre: "Technology", Pages: 352, Available: true, PublishedYear: 1999},
	{ID: 4, Title: "Design Patterns", Author: "Gang of Four", Genre: "Technology", Pages: 395, Available: true, PublishedYear: 1994},
}

var nextID = 5 // Simple counter for generating new IDs

// containsIgnoreCase checks if a string contains another string (case-insensitive)
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// bookMatchesSearch checks if a book matches the search term in title or author
func bookMatchesSearch(book Book, searchTerm string) bool {
	return containsIgnoreCase(book.Title, searchTerm) || containsIgnoreCase(book.Author, searchTerm)
}

// findBookByID searches for a book by its ID and returns the book and its index
func findBookByID(id int) (*Book, int, error) {
	for i, book := range library {
		if book.ID == id {
			return &book, i, nil
		}
	}
	return nil, -1, errors.New("book not found")
}

// getAllBooks handles GET /books - retrieves all books with optional filtering
// supports query parameters: ?search=, ?genre=, ?available=
func getAllBooks(c *fiber.Ctx) error {
	search := c.Query("search")
	genre := c.Query("genre")
	availableStr := c.Query("available")

	var filtered []Book

	for _, book := range library {
		// Filter by search term if provided
		if search != "" && !bookMatchesSearch(book, search) {
			continue
		}

		// Filter by genre if provided
		if genre != "" && !strings.EqualFold(book.Genre, genre) {
			continue
		}

		// Filter by availability if provided
		if availableStr != "" {
			available, err := strconv.ParseBool(availableStr)
			if err == nil && book.Available != available {
				continue
			}
		}

		filtered = append(filtered, book)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"books": filtered,
		"total": len(filtered),
	})
}

// createBook handles POST /books - adds a new book to the library
func createBook(c *fiber.Ctx) error {
	var newBook Book

	if err := c.BodyParser(&newBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
		})
	}

	// Basic validation
	if newBook.Title == "" || newBook.Author == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title and Author are required",
		})
	}

	// Assign ID and add to library
	newBook.ID = nextID
	nextID++
	library = append(library, newBook)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book added successfully",
		"book":    newBook,
	})
}

// getBook handles GET /books/:id - retrieves a specific book by ID
func getBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	book, _, err := findBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

// updateBook handles PUT /books/:id - updates book information
func updateBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	_, index, err := findBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	var updates Book
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
		})
	}

	// Keep the original ID and update the book
	updates.ID = id
	library[index] = updates

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book updated successfully",
		"book":    updates,
	})
}

// toggleBookAvailability handles PATCH /books/:id - toggles the availability status
func toggleBookAvailability(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	book, index, err := findBookByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	// Toggle availability status
	library[index].Available = !book.Available

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book availability updated",
		"book":    library[index],
	})
}

// deleteBook handles DELETE /books/:id - removes a book from the library
func deleteBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid book ID",
		})
	}

	for i, book := range library {
		if book.ID == id {
			library = slices.Delete(library, i, i+1)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Book deleted successfully",
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Book not found",
	})
}

// main function sets up the Fiber app and defines routes
func main() {
	// Create Fiber app with custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Library API is running! 📚",
			"version": "1.0.0",
			"endpoints": []string{
				"GET /books - Get all books",
				"POST /books - Add new book",
				"GET /books/:id - Get book by ID",
				"PUT /books/:id - Update book",
				"PATCH /books/:id - Toggle availability",
				"DELETE /books/:id - Delete book",
			},
		})
	})

	// Book routes
	app.Get("/books", getAllBooks)
	app.Post("/books", createBook)
	app.Get("/books/:id", getBook)
	app.Put("/books/:id", updateBook)
	app.Patch("/books/:id", toggleBookAvailability)
	app.Delete("/books/:id", deleteBook)

	// Start server
	log.Fatal(app.Listen(":8080"))
}
