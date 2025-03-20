package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite" // تشغيل SQLite بدون CGO
)

// ✅ تعريف نموذج المستخدم
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
}

// ✅ تعريف متغير قاعدة البيانات
var db *gorm.DB

func main() {
	// ✅ الحصول على المنفذ من متغير البيئة (لتشغيله على Railway)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // القيمة الافتراضية إذا لم يتم تحديدها
	}

	var err error
	db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("فشل في الاتصال بقاعدة البيانات:", err)
	}

	// ✅ إنشاء جدول المستخدمين
	db.AutoMigrate(&User{})

	// ✅ إعداد الـ API باستخدام Gin
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.POST("/signup", signupHandler)

	// ✅ تشغيل السيرفر
	log.Println("🚀 Running on port:", port)
	r.Run(":" + port)
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

	// ✅ مقارنة كلمة المرور المشفرة
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

	// ✅ تشفير كلمة المرور قبل الحفظ
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	userInput.Password = string(hashedPassword)

	// ✅ إضافة المستخدم إلى قاعدة البيانات
	result := db.Create(&userInput)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
