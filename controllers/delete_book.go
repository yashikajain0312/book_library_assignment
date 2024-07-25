package controllers

import (
    "encoding/csv"
    "net/http"
    "os"
    "strings"
	"path/filepath"

    "github.com/gin-gonic/gin"
	"book-library-assignment/auth"
)

type DeleteBookRequest struct {
    BookName        string `json:"book_name" binding:"required"`
}

func DeleteBookHandler(c *gin.Context) {
    userType, err := auth.GetUserTypeFromToken(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token", "details": err.Error()})
        return
    }

    if userType != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can delete books"})
        return
    }

	var requestBody DeleteBookRequest

    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
        return
    }

    if requestBody.BookName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Book name is required"})
        return
    }

	filePath := filepath.Join("data", "regularUser.csv")

    if err := deleteBookFromFile(filePath, requestBody.BookName); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully!"})
}

func deleteBookFromFile(filePath, bookName string) error {
    file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }

    // Write records back to file excluding the book to be deleted
    var updatedRecords [][]string
    for _, record := range records {
        if strings.EqualFold(record[0], bookName) { // Ignore case while comparing book names
            continue
        }
        updatedRecords = append(updatedRecords, record)
    }

    // Truncate file and write updated records
    if err := file.Truncate(0); err != nil {
        return err
    }

    if _, err := file.Seek(0, 0); err != nil {
        return err
    }

    writer := csv.NewWriter(file)
    defer writer.Flush()

    if err := writer.WriteAll(updatedRecords); err != nil {
        return err
    }

    return nil
}
