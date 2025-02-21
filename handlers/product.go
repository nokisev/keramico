package handlers

import (
	"alice/keramico/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, description, price, seller_id, created_at FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.SellerID, &p.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context, db *sql.DB) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO products (name, description, price, seller_id) VALUES (?,?,?,?)", p.Name, p.Description, p.Price, p.SellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	c.JSON(http.StatusCreated, p)
}
