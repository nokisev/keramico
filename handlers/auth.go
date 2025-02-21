package handlers

import (
	"alice/keramico/internal/redis"
	"alice/keramico/models"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context, db *sql.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(user.Password)), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Сохранение в БД
	result, err := db.Exec("INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)",
		strings.TrimSpace(user.Username),
		strings.TrimSpace(user.Email),
		hashedPassword, // Убрано явное преобразование в строку
		"user",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)
	user.Password = "" // Не возвращаем пароль
	c.JSON(http.StatusCreated, user)
}

func Login(c *gin.Context, db *sql.DB, redisClient *redis.RedisClient) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Очистка пробелов
	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	var user models.User
	err := db.QueryRow("SELECT id, username, email, password, role FROM users WHERE email = ?", input.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	plainPass := []byte(input.Password)

	// Сравнение пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), plainPass); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret123"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	err = redisClient.StoreToken(strconv.Itoa(user.ID), tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token in Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
