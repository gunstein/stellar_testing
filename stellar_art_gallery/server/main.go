package main

import (
	"context"
	"flag"
	"fmt"
	"net/smtp"
		
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/controllers"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)

//TODO: Use more secure email
func SendEmailSMTP(emailTo string, emailFrom string, emailHost string, emailPassword string, emailPort string,
	big_image_url string) (bool, error) {

	emailAuth := smtp.PlainAuth("", emailFrom, emailPassword, emailHost)

	emailBody := fmt.Sprintf(`<html><body><h1>Downloadlink from StellarGallery</h1>
	<a href="%s">%s</a>
	</body></html>`, big_image_url, big_image_url)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Download your purchased art" + "!\n"
	msg := []byte(subject + mime + "\n" + emailBody)
	addr := fmt.Sprintf("%s:%s", emailHost, emailPort)

	emailToArray := []string{emailTo}
	if err := smtp.SendMail(addr, emailAuth, emailFrom, emailToArray, msg); err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func main() {
	account_publickey := flag.String("account", "not_set_account", "account public key")
	emailFrom := flag.String("email_from", "not_set_email_from", "email from")
	emailHost := flag.String("email_host", "not_set_email_host", "email host")
	emailPassword := flag.String("email_password", "not_set_email_password", "email password")
	emailPort := flag.String("email_port", "not_set_email_port", "email port")
	flag.Parse();
	// Build and connect to database
	models.ConnectDatabase()

	client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
	// payments for an account
	fmt.Println(*account_publickey)
	opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
	//opRequest := horizonclient.OperationRequest{ForAccount: "GBGJFGCDZHQ3LXJOUK7EOZB77OR2GMES3FVQRK4M724THUDLZLP7J6A7"}
	//opRequest := horizonclient.OperationRequest{Cursor: "760209215489"}
	ctx := context.Background()

	paymentHandler := func(op operations.Operation) {
		fmt.Println("Payment received.")

		transaction, err := client.TransactionDetail(op.GetTransactionHash())
		if err != nil {
			fmt.Println(err)
			return
		}
		//confirm type is payment
		payment, ok := op.(operations.Payment)
		if !ok {
			fmt.Println("Not payment type operation")
			return
		}
		if strings.EqualFold(payment.From, *account_publickey) {
			//Not interested in outgoing payments
			return
		}
		//Get orderid from memo		
		memo := transaction.Memo
		u64, err := strconv.ParseUint(memo, 10, 32)
		if err != nil {
			fmt.Println(err)
			return
		}
		order, err := models.GetOrderForId(uint(u64))
		if err != nil {
			fmt.Println(err)
			return
		}
		//get art from orderid
		art, err := models.GetArtForId(order.ArtId)
		if err != nil {
			fmt.Println(err)
			return
		}


		//Check payment amount
		value, err := strconv.ParseFloat(payment.Amount, 32)
		if err != nil {
			fmt.Println("Conversion to float failed.")
			return
		}

		if float32(value) < art.Price {
			fmt.Println("Paid to little.")
			return
		}
		//update order to paid
		order, err = models.UpdateOrderToPaid(order.ID)
		if err != nil {
			fmt.Println("UpdateOrderToPaid failed.")
			return
		}

		sent, err := SendEmailSMTP(order.Email, *emailFrom, *emailHost, *emailPassword, *emailPort, art.BigFileUrl)
		if !sent || err != nil {
			fmt.Println("SendEmailSMTP failed.")
			return
		}

		//update order to emailsent
		order, err = models.UpdateOrderToEmailSent(order.ID)
		if err != nil {
			fmt.Println("UpdateOrderToEmailSent failed.")
			return
		}	
		
	}
	fmt.Println("Before streaming.")
	go func() {
		err := client.StreamPayments(ctx, opRequest, paymentHandler)
		if err != nil {
			fmt.Println(err)
		}
    }()

	fmt.Println("After streaming.")

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	
	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/order", controllers.CreateOrderHandler(*account_publickey))

	// Run the server
	r.Run()
}
