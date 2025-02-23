package handlers

import (
	"alice/keramico/internal/redis"
	"alice/keramico/models"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	reds "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func generateJWT(userID int, role string) (string, error) {

	secret, _ := os.LookupEnv("secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 2).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func Register(c *gin.Context, db *sql.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Сохранение в БД
	result, err := db.Exec("INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)",
		user.Username,
		user.Email,
		string(hashedPassword), // Убрано явное преобразование в строку
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
		log.Println(plainPass)
		log.Println([]byte(user.Password))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Генерация JWT
	tokenString, err := generateJWT(user.ID, user.Role)
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

func Logout(c *gin.Context, rds *redis.RedisClient) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDFloat, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	userIDInt := int(userIDFloat)

	log.Printf("Deleting token for user %d", userIDInt)
	err := rds.DeleteToken(strconv.Itoa(userIDInt))
	if err != nil {
		log.Printf("Failed to delete token for user %d: %v", userIDInt, err)

		if errors.Is(err, reds.Nil) {
			c.JSON(http.StatusOK, gin.H{"message": "Token already deleted"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
