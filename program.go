package main

import (
	"gobasic/model"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Set up database connection
	dsn := "user=tanny password=Unw9tg2JXJ2GKslRdb63jyeDgUInNMKG dbname=gologin port=5432 host=dpg-ct8m81m8ii6s73ccndc0-a.singapore-postgres.render.com"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(cors.Default())

	// POST /login route (for authentication)
	r.POST("/login", func(c *gin.Context) {
		var json model.Login

		// Bind incoming JSON to the Login struct
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user exists in the database
		var user model.User
		if err := db.Where(`"user" = ?`, json.User).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Compare the provided password with the hashed password from the database
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password))
		if err != nil {
			// Password mismatch
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Successful login
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Login Success",
			"userID":  user.ID,
		})
	})

	// POST /register route (for user registration)
	r.POST("/register", func(c *gin.Context) {
		var json model.Login

		// Bind incoming JSON to the Login struct
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user already exists
		var userExist model.User
		if err := db.Where("user = ?", json.User).First(&userExist).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		// Hash the password before saving
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
			return
		}

		// Create new user in the database
		user := model.User{
			User:     json.User,
			Password: string(encryptedPassword),
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
			return
		}

		// Return a success response
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "User registered successfully",
			"userID":  user.ID,
		})
	})

	// Run the server on localhost:8080
	r.Run(":8081")
}
