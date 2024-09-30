package main

import (
	"fmt"

	"log"
	"os"
	"time"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // Docker service name of the PostgreSQL container
	port     = 5433         // Internal PostgreSQL port inside the container
	user     = "admin"      // From environment POSTGRES_USER
	password = "password"   // From environment POSTGRES_PASSWORD
	dbname   = "mydatabase" // From environment POSTGRES_DB
)

func main() {
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Output for verification
	fmt.Println(db)
	// db.Migrator().DropColumn(&Book{}, "publisher")
	db.AutoMigrate(&Book{})
	fmt.Println("Migrate successful!")

	//Setup Fiber
	app := fiber.New()

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(getBooks(db))
	})

	app.Get("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book := getBook(db, id)

		return c.JSON(book)
	})

	app.Post("/books", func(c *fiber.Ctx) error {
		book := new(Book)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err := createBook(db, book)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Create Book Successful!",
		})
	})

	app.Put("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book := new(Book)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book.ID = uint(id)

		err = updateBook(db, book)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Update Book Successful!",
		})

	})

	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = deleteBook(db, id)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Delete Book Successful!",
		})
	})

	app.Listen(":8082")

}
