package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"
)

// GET /art
type FindArtOutput struct {
	ArtId        uint   `json:"artid" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	SmallFileUrl string `json:"small_file_url" binding:"required"`
}

// Find all art
func FindArt(c *gin.Context) {
	var art []models.Art

	//Think its better to do above.
	models.DB.Select("ID", "Name", "Description", "SmallFileName").Find(&art)

	//Convert from Artarray to outputarray
	var output []FindArtOutput
	for _, b := range art {
		element := FindArtOutput{ArtId: b.ID, Name: b.Name, Description: b.Description, SmallFileUrl: b.SmallFileName}
		output = append(output, element)
	}

	c.JSON(http.StatusOK, gin.H{"data": output})
}

type CreateOrderInput struct {
	Email string `json:"email" binding:"required"`
	ArtId uint   `json:"artid" binding:"required"`
}

// POST /orders
// Create new order
func CreateOrder(c *gin.Context) {
	// Validate input
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create order
	order := models.Order{Email: input.Email, ArtId: input.ArtId}
	models.DB.Create(&order)

	c.JSON(http.StatusOK, gin.H{"data": order})
}
