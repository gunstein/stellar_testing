package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server_stellar_gallery/controllers"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server_stellar_gallery/models"

	"github.com/stellar/go/clients/horizonclient"
)

func main() {
	account_publickey := flag.String("account", "not_set_account", "account public key")
	flag.Parse()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	
		client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
		client.HorizonURL = "https://34.231.194.216/"
		
		opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := client.StreamPayments(ctx, opRequest, controllers.CreatePaymentHandler(c, *account_publickey, client))
		if err != nil {
			fmt.Println(err)
		}
	
		c.Writer.Flush()
		<-c.Writer.CloseNotify()
		// do something after client is gone
		fmt.Println("Client gone")
	})	

	// Run the server
	r.Run(":5000")
}
