package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ✅ ملف تخزين المستخدمين
const usersFile = "users.json"

// ✅ نموذج المستخدم
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ✅ قائمة المستخدمين في الذاكرة
var users []User
var mutex sync.Mutex // منع التعديل المتزامن على الملف

// ✅ تحميل المستخدمين من JSON عند بدء التشغيل
func loadUsers() error {
	file, err := os.Open(usersFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // إذا لم يكن الملف موجودًا، لا توجد مشكلة
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&users)
}

// ✅ حفظ المستخدمين إلى JSON
func saveUsers() error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Create(usersFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // جعل JSON أكثر قابلية للقراءة
	return encoder.Encode(users)
}

func main() {
	// ✅ تحميل البيانات عند بدء التشغيل
	if err := loadUsers(); err != nil {
		log.Fatal("❌ فشل في تحميل المستخدمين:", err)
	}

	// ✅ الحصول على المنفذ من البيئة
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ✅ إعداد API باستخدام Gin
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.POST("/signup", signupHandler)

	// ✅ تشغيل الخادم
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

	// ✅ البحث عن المستخدم في JSON
	for _, user := range users {
		if user.Email == userInput.Email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
				return
			}
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
}

// ✅ دالة إنشاء حساب
func signupHandler(c *gin.Context) {
	var userInput User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// ✅ التحقق من وجود البريد مسبقًا
	for _, user := range users {
		if user.Email == userInput.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	// ✅ تشفير كلمة المرور قبل الحفظ
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	userInput.Password = string(hashedPassword)

	// ✅ إضافة المستخدم الجديد
	users = append(users, userInput)
	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
