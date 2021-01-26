package main

import (
	"context"
	"fmt"
	"flag"
		
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/controllers"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"

	"github.com/stellar/go/clients/horizonclient"
)

func main() {
	account_publickey := flag.String("account", "not_set_account", "account public key")
	flag.Parse()

	// Build and connect to database
	models.ConnectDatabase()

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	
	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/order", controllers.CreateOrderHandler(*account_publickey))
	r.GET("/order/:memo/big_file_url/:key", controllers.FindBigFileUrl)

	r.GET("/stream", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
	
		client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
		opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
		ctx, cancel := context.WithCancel(context.Background())
		err := client.StreamPayments(ctx, opRequest, controllers.CreatePaymentHandler(c, *account_publickey, client))
		if err != nil {
			fmt.Println(err)
		}
	
		c.Writer.Flush()
		<-c.Writer.CloseNotify()
		// do something after client is gone
		cancel()
	})	

	// Run the server
	r.Run(":8080")
}
