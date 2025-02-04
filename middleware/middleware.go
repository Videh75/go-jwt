package middleware

import (
	"context"
	"fmt"
	"go-jwt/database"
	"go-jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func RequireAuth(c *gin.Context) {
	// Get cookie from request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Expired"})
		}
		coll := database.Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("DB_COLLECTION"))
		var user models.User
		filter := bson.D{{Key: "email", Value: claims["sub"]}}
		err := coll.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		c.Set("user", user.Name)
		c.Next()

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired Token"})
		return
	}
}
