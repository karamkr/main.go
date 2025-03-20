package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // تشغيل SQLite بدون CGO
)

// ✅ نموذج المستخدم
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}

// ✅ قاعدة البيانات
var db *gorm.DB

func main() {
	// ✅ تحديد المنفذ
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ✅ الاتصال بقاعدة البيانات
	var err error
	db, err = gorm.Open(sqlite.Dialector{DSN: "users.db"}, &gorm.Config{})
	if err != nil {
		log.Fatal("فشل في الاتصال بقاعدة البيانات:", err)
	}

	// ✅ إنشاء الجدول إذا لم يكن موجودًا
	if !db.Migrator().HasTable(&User{}) {
		db.AutoMigrate(&User{})
	}

	// ✅ إعداد API
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.POST("/signup", signupHandler)

	// ✅ تشغيل الخادم
	log.Println("🚀 Running on port:", port)
	r.Run()
}

// ✅ دالة تسجيل الدخول
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

	// ✅ مقارنة كلمة المرور
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// ✅ دالة إنشاء حساب
func signupHandler(c *gin.Context) {
	var userInput User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ التحقق من البريد الإلكتروني
	var existingUser User
	if db.Where("email = ?", userInput.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// ✅ تشفير كلمة المرور
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	userInput.Password = string(hashedPassword)

	// ✅ حفظ المستخدم الجديد
	db.Create(&userInput)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
