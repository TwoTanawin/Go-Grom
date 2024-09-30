package main

import (
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name  string `json:"name"`
	Autor string `json:"authen"`
	// Publisher   string
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(&book)

	if result.Error != nil {
		return result.Error
	}

	// fmt.Println("Create Book Successful!")
	return nil
}

func getBook(db *gorm.DB, id int) *Book {
	var book Book
	result := db.First(&book, id)

	if result.Error != nil {
		log.Fatal("Error geting book: %v", result.Error)
	}

	return &book
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)

	if result.Error != nil {
		log.Fatal("Error geting books: %v", result.Error)
	}

	return books
}

func updateBook(db *gorm.DB, book *Book) error {
	result := db.Model(&book).Updates(book)

	if result.Error != nil {
		return result.Error
	}

	// fmt.Println("Update Book Successful!")
	return nil

}

func deleteBook(db *gorm.DB, id int) error {
	var book Book
	result := db.Delete(&book, id)

	if result.Error != nil {
		return result.Error
	}

	// fmt.Println("Delete Book Successful!")
	return nil
}

func searchBook(db *gorm.DB, bookName string) *Book {
	var book Book

	result := db.Where("name = ?", bookName).First(&book)

	if result.Error != nil {
		log.Fatal("Error Search book: %v", result.Error)
	}

	return &book
}

func searchBooks(db *gorm.DB, bookName string) []Book {
	var books []Book

	result := db.Where("name = ?", bookName).Order("price").Find(&books)

	if result.Error != nil {
		log.Fatal("Error Search book: %v", result.Error)
	}

	return books
}
