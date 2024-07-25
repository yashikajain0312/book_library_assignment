// main.go

package main

import (
    "github.com/gin-gonic/gin"
    "book-library-assignment/controllers" 
)

func main() {
    router := gin.Default()

    router.POST("/login", controllers.LoginHandler)
	router.GET("/home", controllers.HomeHandler)
	router.POST("/addBook", controllers.AddBookHandler)
	router.DELETE("/deleteBook", controllers.DeleteBookHandler)

    router.Run(":6000")
}
