package main

import (
	"context"
	"fmt"
	"flag"
	"io"
		
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

	//used to inform consumers about incoming payments
	payments := make(chan controllers.Message)

	//Handling all received payments
	client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
	opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
	ctx := context.Background()
	go func() {
		err := client.StreamPayments(ctx, opRequest, controllers.CreatePaymentHandler(payments, *account_publickey, client))
		if err != nil {
			fmt.Println(err)
		}
    }()

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	
	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/order", controllers.CreateOrderHandler(*account_publickey))

	r.GET("/stream/:memo", func(c *gin.Context) {
		memo_param := c.Param("memo")
		c.Stream(func(w io.Writer) bool {
			if message, ok := <-payments; ok {
				if message.Memo == memo_param{
					c.SSEvent("message", message.Url)
					return true	
				}
			}
			return false
		})
	})

	// Run the server
	r.Run()
}
