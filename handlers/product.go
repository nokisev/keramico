package handlers

import (
	"alice/keramico/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, description, image, rating, price FROM products")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Rating, &p.Price)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var product models.Product
	err := db.QueryRow("SELECT id, name, description, image, rating, price FROM products WHERE id = ?", id).Scan(&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Rating,
		&product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		}

		return
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context, db *sql.DB) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO products (name, description, image, rating, price) VALUES (?,?,?,?,?)", p.Name, p.Description, p.Image, p.Rating, p.Price)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	c.JSON(http.StatusCreated, p)
}
