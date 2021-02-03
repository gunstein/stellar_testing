package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"

	broadcast "github.com/dustin/go-broadcast"
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

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	
	broadcaster := broadcast.NewBroadcaster(100)
	defer broadcaster.Close()
	//r.Use(stream.ServeHTTP())

	go func(){

		for{
			//Start listening for payments on the shops account
			client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
			opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
			err := client.StreamPayments(context.Background(), opRequest, controllers.CreatePaymentHandler(broadcaster, *account_publickey, client))
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/order", controllers.CreateOrderHandler(*account_publickey))
	r.GET("/order/:memo/big_file_url/:key", controllers.FindBigFileUrl)

	r.GET("/stream", func(c *gin.Context) {
		ch := make(chan interface{})
		broadcaster.Register(ch)
		defer broadcaster.Unregister(ch)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Flush()
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-ch; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})	

	// Run the server
	r.Run(":5000")
}
