package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"
)

// GET /art
type FindArtOutput struct {
	ArtId        uint    `json:"artid" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Artist       string  `json:"artist" binding:"required"`
	Price        float32 `json:"price" binding:"required"`
	SmallFileUrl string  `json:"small_file_url" binding:"required"`
}

// Find all art
func FindArt(c *gin.Context) {
	var art []models.Art

	models.DB.Select("ID", "Title", "Description", "Artist", "Price", "SmallFileUrl").Find(&art)

	//Convert from Artarray to outputarray
	var output []FindArtOutput
	for _, b := range art {
		element := FindArtOutput{ArtId: b.ID, Title: b.Title, Description: b.Description, Artist: b.Artist, Price: b.Price, SmallFileUrl: b.SmallFileUrl}
		output = append(output, element)
	}

	c.JSON(http.StatusOK, gin.H{"data": output})
}

type CreateOrderInput struct {
	ArtId uint   `json:"artid" binding:"required"`
}

type CreateOrderOutput struct {
	Account string `json:"account" binding:"required"`
	Memo    string `json:"memo" binding:"required"`
}

// POST /orders
// Create new order
//This is one way to get account into Handlerfunc as parameter
//https://stackoverflow.com/questions/34046194/how-to-pass-arguments-to-router-handlers-in-golang-using-gin-web-framework
func CreateOrderHandler(account string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Validate input
		var input CreateOrderInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create order
		order := models.Order{ArtId: input.ArtId}
		models.DB.Create(&order)

		var output = CreateOrderOutput{Account: account, Memo: strconv.FormatUint(uint64(order.ID), 10)}
		c.JSON(http.StatusOK, gin.H{"data": output})
	}
	return gin.HandlerFunc(fn)
}
