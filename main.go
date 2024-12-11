package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func main() {

	s := http.NewServeMux()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// 1. Initiate Snap client
		var s = snap.Client{}
		s.New("YOUR-SERVER-KEY", midtrans.Sandbox)
		// Use to midtrans.Production if you want Production Environment (accept real transaction).

		// 2. Initiate Snap request param
		req := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  "YOUR-ORDER-ID-12345",
				GrossAmt: 100000,
			},
			CreditCard: &snap.CreditCardDetails{
				Secure: true,
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: "John",
				LName: "Doe",
				Email: "john@doe.com",
				Phone: "081234567890",
			},
		}

		// 3. Execute request create Snap transaction to Midtrans Snap API
		snapResp, _ := s.CreateTransaction(req)
		fmt.Fprintln(w, snapResp)
	})
	log.Println("server on running")
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		log.Fatal(err)
	}

}
