package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Order struct {
	OrderID     int    `json:"id,int"`
	OrderName   string `json:"order_name"`
	Status      string `json:"order_status"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	//	ActiveOperators int    `json:"active_operators"`
}

var erpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// Pentru comunicare ERP -> UNS -> OT
// structura cu:
// number of parts
// cmd start/stop
// conveyor number
// type ERPCommand {
//
// }
type ITData []Order

func fetchITData() (ITData, error) {
	var it ITData
	response, err := erpClient.Get("http://localhost:8083/orders")
	if err != nil {
		return it, err
	}
	// debug response
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&it)

	return it, err
}

/* IT mock data
func itDataHandler(w http.ResponseWriter, r *http.Request) {
	data := Order{
		OrderID:         1,
		OrderName:       "mock_order_name",
		Status:          "Pending",
		ProductName:     "mock_product_name",
		Quantity:        10,
		ActiveOperators: 2,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
*/

/* IT mock server
func startITServer() {
	http.HandleFunc("/api/orders", itDataHandler)
	fmt.Println("Started HTTP server running on http://localhost:8089/api/orders")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatalf("Starting HTTP server failed: %v", err)
	}
}
*/
