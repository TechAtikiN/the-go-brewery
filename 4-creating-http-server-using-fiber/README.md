# Building a RESTful API with Go and Fiber ⚡️
Go ecosystem has various frameworks like **Gin**, **Echo**, **Chi**, **Fiber** and more. Among these the most popular one in terms of **Simplicity** and **Performance** is **Fiber**.

Fiber is inspired by **Express.js** (Node.js framework) due to which Node.js developers have an "_Ah-ha_!" moment when they get started with Fiber. Let's explore Fiber in more detail and get our hands dirty by building a **RESTful API** with it.

> In the last blog, we explored building a RESTful API using the Gin framework. If you haven't checked it out yet, you can find it [here](../3-creating-http-server-using-gin/README.md).

## About `Fiber` ⚡️
- Uses **`fasthttp`**  instead of the regular `net/http` package.
    > 💡 **`fasthttp`** package is a high-performance, low-memory HTTP server and client for Go, designed 
for efficiency. It offers significantly improved "throughput" and lower "latency" compared to Go's standard net/http package. [Check more details](https://pkg.go.dev/github.com/valyala/fasthttp)

- **Express inspired structure**, making it intuitive and easy to learn for Node.js developers.
- **Middleware support** for logging, CORS, rate limiting and more.
- Built-in support for JSON parsing, request parsing and serving static files.
- Out of the box support for **WebSockets**, **template engines**, and **testing**.

## Specific to Fiber~
- **`fiber.New()`**: Creates a new Fiber app instance. You can pass configuration options to customize it (e.g. custom error handling)
- **`c.BodyParser(&struct)`**: Parses request bodies into Go structs automatically
- **`c.Params("key")`**: Grabs URL parameters (like `/books/:id`)
- **`c.Query("key")`**: Gets query string values (like `?search=golang`)
- **`app.Listen(":8080")`**: Starts the HTTP server on the specified port
- **Route handlers**: Defines a route for a specific HTTP method and path. For example: `app.Get()`, `app.Post()`

## Building Your Server 🚀

We'll build a RESTful API to manage a simple library of books stored in memory. The API will support the following endpoints:

- **GET /books** → Get all books (with search and filtering)
- **POST /books** → Add a new book  
- **GET /books/:id** → Get details of a specific book
- **PUT /books/:id** → Update book information
- **DELETE /books/:id** → Remove a specific book

Let's dive in!

Before writing code, ensure you have Go installed on your system. 
Create a new directory for your project, initialize a Go module, and install the Fiber package.

```bash
~ $ mkdir books-api

~ $ cd books-api

~/books-api $ go mod init books-api

~/books-api $ go get github.com/gofiber/fiber/v2
```

### 1. Start simple
Let's start with a minimal Fiber server to test the setup. Create a file named `main.go`:

```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New() // Create a new Fiber app

    // Define a simple route
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Welcome to our Books API! 📚",
        })
    })

    log.Fatal(app.Listen(":8080")) // Start server on port 8080
}
```

Run the server using the following command:

```bash
~/books-api $ go run main.go
```

> You can setup `air` for live reloading during development. Check out 
[air](https://github.com/air-verse/air)

Open your browser and navigate to http://localhost:8080. You should see the following output:

```json
{
    "message": "Welcome to our Books API! 📚"
}
```

### 2. Define the Book Model and In-Memory Store
> 💡 This example uses an in-memory store for simplicity. In production, you’d likely use a database (e.g. PostgreSQL, MySQL, MongoDB) with an ORM like GORM.

```go
package main

import (
    "errors"
    "log"
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

// In-memory store for books
var library = []Book{
    {ID: 1, Title: "The Go Programming Language", Author: "Alan Donovan", Genre: "Technology", Pages: 380, Available: true, PublishedYear: 2015},
    {ID: 2, Title: "Clean Code", Author: "Robert Martin", Genre: "Technology", Pages: 464, Available: false, PublishedYear: 2008},
    {ID: 3, Title: "The Pragmatic Programmer", Author: "Dave Thomas", Genre: "Technology", Pages: 352, Available: true, PublishedYear: 1999},
}

var nextBookID = 4 // Next ID for new books
```

### 3. Writing the API Endpoints

#### 3.1. GET /books
This endpoint retrieves all books, with optional search and filtering by genre and availability.

```go
func getAllBooks(c *fiber.Ctx) error {
    result := library
    
    // for search and filter logics refer to the complete code

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "books": result,
        "total": len(result),
    })
}
```

#### 3.2. GET /books/:id

This endpoint retrieves a specific book by its ID.


```go
func findBookByID(id int) (*Book, int, error) {
    for i, book := range library {
        if book.ID == id {
            return &book, i, nil
        }
    }
    return nil, -1, errors.New("book not found")
}

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
```

#### 3.3. Adding New Books

```go
func createBook(c *fiber.Ctx) error {
    var newBook Book
    
    // Parse the JSON body into our Book struct
    if err := c.BodyParser(&newBook); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Cannot parse JSON",
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
    newBook.ID = nextBookID
    nextBookID++
    library = append(library, newBook)
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Book added successfully",
        "book": newBook,
    })
}
```

To check the **Update** and **Delete** endpoints, you can refer to the [complete code](./main.go).

### 4. Wiring Everything Together 🔌

Now let's put it all together with a global error handler, middleware, and route definitions in the `main.go` file:

```go
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

    // Routes
    app.Get("/books", getAllBooks)
    app.Post("/books", createBook)
    app.Get("/books/:id", getBook)
    app.Put("/books/:id", updateBook)
    app.Delete("/books/:id", deleteBook)

    // Start server
    log.Fatal(app.Listen(":8080"))
}
```

> For in-depth testing, consider using tools like [Postman](https://www.postman.com/) or [Thunder client](https://www.thunderclient.com/).

## TL;DR
- Fiber is built on top of **fasthttp** for high performance and low memory usage.
- Fiber is a great choice if you want **Express.js-like** simplicity with Go’s speed and efficiency.
- Fiber provides built-in middleware for logging, CORS and simplifies request parsing and define routing in a clean and easy way.
- Custom error handling is easy by utilizing the global ErrorHandler config or using middlewares.
- Native support for WebSockets, Static file serving, Template engines, and Multipart form parsing.
- Learn more about Fiber from the [official documentation](https://docs.gofiber.io/).