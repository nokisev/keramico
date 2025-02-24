package handlers

import (
	"alice/keramico/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context, db *sql.DB) {
	userID := c.Param("id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var profile models.Profile
	err = db.QueryRow("SELECT user_id, fullname, avatar, bio FROM profiles WHERE user_id = ?", id).
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
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userIDFromContext, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Преобразуем userID из float64 в int
	userIDFloat, ok := userIDFromContext.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	currentUserID := int(userIDFloat)
	if id != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}
	currentProfile, err := getCurrentProfile(db, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current profile: " + err.Error()})
		return
	}

	var updatedProfile models.Profile
	if err := c.ShouldBindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if updatedProfile.FullName != "" {
		currentProfile.FullName = updatedProfile.FullName
	}
	if updatedProfile.Avatar != "" {
		currentProfile.Avatar = updatedProfile.Avatar
	}
	if updatedProfile.Bio != "" {
		currentProfile.Bio = updatedProfile.Bio
	}
	_, err = db.Exec("UPDATE profiles SET fullname = ?, avatar = ?, bio = ? WHERE user_id = ?", currentProfile.FullName, currentProfile.Avatar, currentProfile.Bio, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func getCurrentProfile(db *sql.DB, userID int) (models.Profile, error) {
	var profile models.Profile
	err := db.QueryRow("SELECT fullname, avatar, bio FROM profiles WHERE user_id = ?", userID).
		Scan(&profile.FullName, &profile.Avatar, &profile.Bio)
	if err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}
