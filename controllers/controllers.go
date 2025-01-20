package controllers

import (
	"context"
	"fmt"
	"go-jwt/database"
	"go-jwt/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	coll := database.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))

	//Getting user input from body
	var doc models.User
	if err := c.Bind(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Hashing the password given by the user
	hash, err := bcrypt.GenerateFromPassword([]byte(doc.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	//Creating empty user object and giving corresponding values
	user := models.User{
		Email:    doc.Email,
		Password: string(hash),
		Name:     doc.Name,
	}

	//Storing the populated user object along with hashed password into the db
	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
