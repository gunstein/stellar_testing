package controllers

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/gunstein/stellar_testing/stellar_art_gallery/server_stellar_gallery/models"

	broadcast "github.com/dustin/go-broadcast"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
)


func CreatePaymentHandler(broadcaster broadcast.Broadcaster, account string, client *horizonclient.Client ) func(operations.Operation){
	paymentHandler := func(op operations.Operation) {
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
		fmt.Println("Payment received. memo: ", memo)

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

		broadcaster.Submit(memo)
	}

	return paymentHandler
}
