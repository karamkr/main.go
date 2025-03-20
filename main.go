package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // SQLite بدون CGO
)


type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("file:users.db?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		panic("فشل في الاتصال بقاعدة البيانات")
	}

	db.AutoMigrate(&User{})

	r := gin.Default()

	r.POST("/login", loginHandler)
	r.POST("/signup", signupHandler)

	r.Run(":8080")
}

func loginHandler(c *gin.Context) {
	var userInput User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user User
	result := db.Where("email = ?", userInput.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 🔹 مقارنة كلمة المرور المُدخلة بالمحفوظة (المشفرة)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func signupHandler(c *gin.Context) {
	var userInput User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	userInput.Password = string(hashedPassword)

	result := db.Create(&userInput)

	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}
