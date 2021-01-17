package main

import (
	"context"
	"flag"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/controllers"
	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)

//TODO: Use more secure email
func SendEmailSMTP(emailTo string, emailFrom string, emailHost string, emailPassword string, emailPort string,
	big_image_url string) (bool, error) {

	emailAuth := smtp.PlainAuth("", emailFrom, emailPassword, emailHost)

	emailBody := "TODO: build a useful message with url"

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Test Email" + "!\n"
	msg := []byte(subject + mime + "\n" + emailBody)
	addr := fmt.Sprintf("%s:%s", emailHost, emailPort)

	emailToArray := []string{emailTo}
	if err := smtp.SendMail(addr, emailAuth, emailFrom, emailToArray, msg); err != nil {
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

	// Build and connect to database
	models.ConnectDatabase()

	client := horizonclient.DefaultTestNetClient //DefaultPublicNetClient
	// payments for an account
	opRequest := horizonclient.OperationRequest{ForAccount: *account_publickey}
	ctx := context.Background()
	paymentHandler := func(op operations.Operation) {
		transaction, err := client.TransactionDetail(op.GetTransactionHash())
		if err != nil {
			fmt.Println(err)
			return
		}
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
		art, err := models.GetArtForId(order.ArtId)
		if err != nil {
			fmt.Println(err)
			return
		}
		//Check payment amount
		payment, ok := op.(operations.Payment)
		if !ok {
			fmt.Println("Not payment type operation")
			return
		}
		if strings.EqualFold(payment.From, *account_publickey) {
			//Not interested in outgoing payments
			return
		}

		value, err := strconv.ParseFloat(payment.Amount, 32)
		if err != nil {
			fmt.Println("Conversion to float failed.")
			return
		}

		if float32(value) < art.Price {
			fmt.Println("Paid to little.")
			return
		}

		SendEmailSMTP(order.Email, *emailFrom, *emailHost, *emailPassword, *emailPort, art.BigFileUrl)
	}
	err := client.StreamPayments(ctx, opRequest, paymentHandler)
	if err != nil {
		fmt.Println(err)
	}

	r := gin.Default()
	// Routes
	r.GET("/art", controllers.FindArt)
	r.POST("/orders", controllers.CreateOrderHandler(*account_publickey))

	// Run the server
	r.Run()
}
