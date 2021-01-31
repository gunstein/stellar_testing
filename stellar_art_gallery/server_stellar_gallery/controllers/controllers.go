package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server_stellar_gallery/models"
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


type FindBigFileUrlOutput struct {
	BigFileUrl string `json:"big_file_url" binding:"required"`
	Comment    string `json:"comment"`
}

//  /order/:memo/big_file_url/:key
func FindBigFileUrl(c *gin.Context) {

	memo_param := c.Param("memo")
	key_param := c.Param("key")

	var output FindBigFileUrlOutput

	orderid, err := strconv.ParseUint(memo_param, 10, 32)
	if err != nil {
		output.Comment = "memo is not a number."
		c.JSON(http.StatusNotAcceptable, gin.H{"data": output})
		return
	}

	order, err := models.GetOrderForId(uint(orderid))
	if err != nil {
		output.Comment = "order not found."
		c.JSON(http.StatusNotFound, gin.H{"data": output})
		return
	}

	//Check order is paid and key is correct
	if order.Paid == false {
		output.Comment = "Not paid."
		c.JSON(http.StatusNotAcceptable, gin.H{"data": output})
		return
	}
	
	if order.DownloadKey != key_param{
		output.Comment = "Wrong downloadkey."
		c.JSON(http.StatusUnauthorized, gin.H{"data": output})
		return
	}

	//get art from orderid
	art, err := models.GetArtForId(order.ArtId)
	if err != nil {
		output.Comment = "Artid not found."
		c.JSON(http.StatusNotFound, gin.H{"data": output})
		return
	}

	output.BigFileUrl = art.BigFileUrl
	c.JSON(http.StatusOK, gin.H{"data": output})
}

type CreateOrderInput struct {
	ArtId uint   `json:"artid" binding:"required"`
}

type CreateOrderOutput struct {
	Account string `json:"account" binding:"required"`
	Memo    string `json:"memo" binding:"required"`
	DownloadKey string `json:"download_key" binding:"required"`
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

		//Generate downloadKey
		downloadKey := uuid.New().String()

		// Create order
		order := models.Order{ArtId: input.ArtId, DownloadKey: downloadKey}
		models.DB.Create(&order)

		var output = CreateOrderOutput{Account: account, Memo: strconv.FormatUint(uint64(order.ID), 10), DownloadKey: downloadKey}
		c.JSON(http.StatusOK, gin.H{"data": output})
	}
	return gin.HandlerFunc(fn)
}
