package main

import (
	"context"
	"fmt"
	"flag"
	"io"
	//"time"
	//"strconv"
		
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/controllers"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"

	"github.com/stellar/go/clients/horizonclient"
)
/*
func streamer(c *gin.Context){
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")

	c.Stream(func(w io.Writer) bool {
		c.Writer.Write([]byte(time.Now().Format(time.RFC3339Nano) + "\n"))
		c.Writer.Flush()
		return false
	})
}

func streamer2(c *gin.Context){
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Stream(func(w io.Writer) bool {
		c.Writer.Write([]byte("message:gvtest"+ "\n"))
		c.Writer.Flush()
		return true
	})
}

func streamer3(c *gin.Context){
	chanStream := make(chan string)
	go func() {
		defer close(chanStream)
		for i := 0; i < 100; i++ {
			chanStream <- strconv.Itoa(i)
			time.Sleep(time.Second * 1)
		}
	}()
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Stream(func(w io.Writer)bool{
		if msg, ok := <-chanStream; ok {
			c.Writer.Write([]byte("message:"+msg+ "\n"))
			c.Writer.Flush()
		}
		return true
	})
}
*/

func main() {
	account_publickey := flag.String("account", "not_set_account", "account public key")
	flag.Parse()

	// Build and connect to database
	models.ConnectDatabase()

	//used to inform consumers about incoming payments
	//payments := make(chan string)
	//gvtest := make(chan string)

	//Handling all received payments
	client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
	opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
	ctx := context.Background()
	/*
	go func() {
		err := client.StreamPayments(ctx, opRequest, controllers.CreatePaymentHandler(payments, *account_publickey, client))
		if err != nil {
			fmt.Println(err)
		}
	}()
	*/

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	
	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/order", controllers.CreateOrderHandler(*account_publickey))
	r.GET("/order/:memo/big_file_url/:key", controllers.FindBigFileUrl)

	r.GET("/stream", func(c *gin.Context) {
		payments := make(chan string)
		go func() {
			err := client.StreamPayments(ctx, opRequest, controllers.CreatePaymentHandler(payments, *account_publickey, client))
			if err != nil {
				fmt.Println(err)
			}
		}()
		//c.Writer.Header().Set("Content-Type", "text/event-stream")
		//c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Stream(func(w io.Writer)bool{
			if msg, ok := <-payments; ok {
				fmt.Println("write message to stream " + msg)
				c.SSEvent("message", msg)
				//c.Writer.Write([]byte("message:"+msg+ "\n"))
				//c.Writer.Flush()
			}
			return true
		})		
	})	
	/*
	r.GET("/stream", func(c *gin.Context) {
		go func() {
			err := client.StreamPayments(ctx, opRequest, controllers.CreatePaymentHandler(payments, *account_publickey, client))
			if err != nil {
				fmt.Println(err)
			}
		}()
		c.Stream(func(w io.Writer) bool {
			if message, ok := <-payments; ok {
				fmt.Println("before c.SSEvent")
				c.SSEvent("message", message.Memo)
				//return true	
			}
			return false
		})
	})
	*/
/*
	r.GET("/test", func(c *gin.Context) {
		go func() {
			defer close(payments)
			for i := 0; i < 10; i++ {
				gvtest <- strconv.Itoa(i)
				time.Sleep(time.Second * 1)
			}
		}()
	})	

	r.GET("/stream2", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Stream(func(w io.Writer)bool{
			if msg, ok := <-gvtest; ok {
				c.Writer.Write([]byte("message:"+msg+ "\n"))
				c.Writer.Flush()
			}
			return true
		})		
	})
	*/
	//r.GET("/streamtest", streamer3)

	// Run the server
	r.Run(":8080")
}
