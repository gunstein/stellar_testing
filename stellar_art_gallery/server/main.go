package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/controllers"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"
)

func main() {
	r := gin.Default()

	// Connect to database
	models.ConnectDatabase()

	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/orders", controllers.CreateOrder)

	// Run the server
	r.Run()
}

/*package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)

func main() {
	client := horizonclient.DefaultTestNetClient
	// all payments
	opRequest := horizonclient.OperationRequest{Cursor: "760209215489"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(op operations.Operation) {
		fmt.Println(op)
	}
	err := client.StreamPayments(ctx, opRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
	//handleRequests()
}
*/
