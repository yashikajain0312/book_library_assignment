package controllers

import (
	"encoding/csv"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"book-library-assignment/auth"
)

func HomeHandler(c *gin.Context) {
	userType, err := auth.GetUserTypeFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token", "details": err.Error()})
		return
	}

	filePaths := getBookFilePaths(userType)

	// Reading books from CSV files.
	var allBooks []string
	for _, filePath := range filePaths {
		books, err := readBooksFromFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read books from file", "details": err.Error()})
			return
		}
		allBooks = append(allBooks, books...)
	}

	c.JSON(http.StatusOK, gin.H{"books": allBooks})
}

func getBookFilePaths(userType string) []string {
	dataDir := "data"
	if userType == "admin" {
		return []string{filepath.Join(dataDir, "adminUser.csv"), filepath.Join(dataDir, "regularUser.csv")}
	}
	return []string{filepath.Join(dataDir, "regularUser.csv")}
}

func readBooksFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var books []string
	for _, record := range records {
		books = append(books, record[0]) 
	}
	return books, nil
}
