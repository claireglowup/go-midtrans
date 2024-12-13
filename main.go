package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/example"
	"github.com/midtrans/midtrans-go/snap"
)

var s snap.Client

func initializeSnapClient() {
	s.New("SB-Mid-server-jffebVOicyta3cHnr9FU_uge", midtrans.Sandbox)
}

func createTransaction() string {
	// Optional : here is how if you want to set append payment notification for this request
	s.Options.SetPaymentAppendNotification("https://example.com/append")

	// Optional : here is how if you want to set override payment notification for this request
	s.Options.SetPaymentOverrideNotification("https://example.com/override")
	// Send request to Midtrans Snap API

	resp, err := s.CreateTransaction(GenerateSnapReq())
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	return resp.Token
}

func main() {

	fmt.Println("================ Request with Snap Client ================")
	initializeSnapClient()

	http.HandleFunc("/getPayment", handler)
	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	// HTML Content
	token := createTransaction()
	htmlContent := fmt.Sprintf(`
	<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <!-- @TODO: replace SET_YOUR_CLIENT_KEY_HERE with your client key -->
    <script
      type="text/javascript"
      src="https://app.midtrans.com/snap/snap.js"
      data-client-key="SB-Mid-client-xq1cgQGMbKKo6r5U"
    ></script>
    <!-- Note: replace with src="https://app.midtrans.com/snap/snap.js" for Production environment -->
  </head>

  <body>
    <button id="pay-button" style="color: blue">Pay!</button>

    <!-- @TODO: You can add the desired ID as a reference for the embedId parameter. -->
    <div id="snap-container"></div>

    <script type="text/javascript">
      // For example trigger on button clicked, or any time you need
      var payButton = document.getElementById("pay-button");
      payButton.addEventListener("click", function () {
        // Trigger snap popup. @TODO: Replace TRANSACTION_TOKEN_HERE with your transaction token.
        // Also, use the embedId that you defined in the div above, here.
        window.snap.embed("%s", {
          embedId: "snap-container",
          onSuccess: function (result) {
            /* You may add your own implementation here */
            alert("payment success!");
            console.log(result);
          },
          onPending: function (result) {
            /* You may add your own implementation here */
            alert("wating your payment!");
            console.log(result);
          },
          onError: function (result) {
            /* You may add your own implementation here */
            alert("payment failed!");
            console.log(result);
          },
          onClose: function () {
            /* You may add your own implementation here */
            alert("you closed the popup without finishing the payment");
          },
        });
      });
    </script>
  </body>
</html>`, token)

	// Menulis HTML ke Response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

func GenerateSnapReq() *snap.Request {

	// Initiate Customer address
	custAddress := &midtrans.CustomerAddress{
		FName:       "John",
		LName:       "Doe",
		Phone:       "081234567890",
		Address:     "Baker Street 97th",
		City:        "Jakarta",
		Postcode:    "16000",
		CountryCode: "IDN",
	}

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-GO-ID-" + example.Random(),
			GrossAmt: 200000,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    "John",
			LName:    "Doe",
			Email:    "john@doe.com",
			Phone:    "081234567890",
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "ITEM1",
				Price: 200000,
				Qty:   1,
				Name:  "Someitem",
			},
		},
	}
	return snapReq
}
