package controllers

import (
    "encoding/csv"
    "os"
    "path/filepath"
    "strconv"
    "net/http"

    "github.com/gin-gonic/gin"
	"book-library-assignment/auth"
)

type AddBookRequest struct {
    BookName        string `json:"book_name" binding:"required"`
    Author          string `json:"author" binding:"required"`
    PublicationYear int    `json:"publication_year" binding:"required"`
}

func AddBookHandler(c *gin.Context) {
    userType, err := auth.GetUserTypeFromToken(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token", "details": err.Error()})
        return
    }

    if userType != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can add books"})
        return
    }

	var requestBody AddBookRequest

    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
        return
    }

    if !isValidYear(requestBody.PublicationYear) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication year"})
        return
    }

    // adding book to regularUser.csv file.
    filePath := filepath.Join("data", "regularUser.csv")

    if err := addBookToFile(filePath, requestBody); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book added successfully!"})
}

func addBookToFile(filePath string, book AddBookRequest) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // writing the new book data in the CSV file.
    if err := writer.Write([]string{book.BookName, book.Author, strconv.Itoa(book.PublicationYear)}); err != nil {
        return err
    }

    return nil
}

func isValidYear(year int) bool {
    return year >= 0 && year <= 9999
}
