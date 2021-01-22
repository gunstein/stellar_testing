package controllers

import (
	"fmt"
		
	"strconv"
	"strings"

	"github.com/gunstein/stellar_testing/stellar_art_gallery/server/models"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)

type Message struct {
    Memo string
    Url  string
}

//Producer
func CreatePaymentHandler(payments chan<- Message, account string, client *horizonclient.Client ) func(operations.Operation){
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
		if strings.EqualFold(payment.From, account) {
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
		//Inform consumers
		payments <- Message{Memo: memo, Url: art.BigFileUrl}
		/*
		sent, err := SendEmailSMTP(order.Email, *emailFrom, *emailHost, *emailPassword, *emailPort, art.BigFileUrl)
		if !sent || err != nil {
			fmt.Println("SendEmailSMTP failed.")
			return
		}
		*/
		/*
		//update order to emailsent
		order, err = models.UpdateOrderToEmailSent(order.ID)
		if err != nil {
			fmt.Println("UpdateOrderToEmailSent failed.")
			return
		}	*/	
	}

	return paymentHandler
}
