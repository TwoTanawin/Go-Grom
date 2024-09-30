package main

import (
	"fmt"

	"log"
	"os"
	"time"

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

	// newBook := &Book{
	// 	Name:        "ML Engineer",
	// 	Autor:       "Two",
	// 	Description: "Test",
	// 	Price:       500,
	// }

	// createBook(db, newBook)

	currentBook := getBook(db, 5)

	fmt.Println(currentBook)

	currentBook.Name = "New DL book"
	currentBook.Price = 6000

	updateBook(db, currentBook)

	// deleteBook(db, 1)

	// currentBook = searchBook(db, "ML Engineer")
	// fmt.Println(currentBook)

	currentBooks := searchBooks(db, "ML Engineer")
	// fmt.Println(currentBooks)
	for _, book := range currentBooks {
		fmt.Println(book.ID, book.Name, book.Autor, book.Description, book.Price)
	}

}
