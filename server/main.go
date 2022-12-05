package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	DateOfBirth string `json:"dateOfBirth"`
}

func SaveImage(c *gin.Context) {
	file, uploadedFile, err := c.Request.FormFile("file")

	// fmt.Println(file, uploadedFile)

	if err != nil {
		log.Println("image upload error -->", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	fileName := strings.Split(uploadedFile.Filename, ".")[0]

	fileExt := strings.Split(uploadedFile.Filename, ".")[1]

	image, err := os.Create("./images/"+fmt.Sprintf("%s.%s", fileName, fileExt))

	if err != nil {
		log.Println("image save error -->", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	io.Copy(image, file)

	imageUrl := "http://localhost:8080/images/" + fmt.Sprintf("%s.%s", fileName, fileExt)

	data := map[string] string {
		"status": "success",
		"imageUrl": imageUrl,
		"header": uploadedFile.Header.Get(fileName),
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

var users []User

func main() {
	fmt.Println("hello world")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "gin is working",
		})
	})

	r.POST("/photo", SaveImage)


	r.GET("/user", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": &users,
		})
	})

	r.POST("/user", func(c *gin.Context) {
		var input User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := User{
			FirstName: input.FirstName,
			LastName: input.LastName,
			PhoneNumber: input.PhoneNumber,
			DateOfBirth: input.DateOfBirth,
		}
		users = append(users, user)
		c.JSON(http.StatusOK, gin.H{"data": user})
	})
	r.Run()
}
