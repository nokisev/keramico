package handlers

import (
	"alice/keramico/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context, db *sql.DB) {
	userID := c.MustGet("user_id").(int)

	var profile models.Profile

	err := db.QueryRow("SELECT user_id, fullname, avatar, bio FROM profiles WHERE user_id = ?", userID).
		Scan(&profile.UserID, &profile.FullName, &profile.Avatar, &profile.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context, db *sql.DB) {
	userID := c.MustGet("user_id").(int)

	var profile models.Profile

	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := db.Exec("UPDATE profiles SET fullname = ?, avatar = ?, bio = ? WHERE user_id = ?", profile.FullName, profile.Avatar, profile.Bio, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
