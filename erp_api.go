package main

import (
	"encoding/json"
	"net/http"
)

type OrderStatus int

const (
	Pending OrderStatus = iota
	Completed
	InProgress
)

var orderStatusName = map[OrderStatus]string{
	Pending:    "pending",
	Completed:  "completed",
	InProgress: "in_progress",
}

type Order struct {
	OrderID         int    `json:"id,string"`
	OrderName       string `json:"order_name"`
	Status          string `json:"order_status"`
	ProductName     string `json:"product_name"`
	Quantity        int    `json:"quantity"`
	ActiveOperators int    `json:"active_operators"`
}

type ITData []Order

func fetchITData() (ITData, error) {
	var it ITData
	response, err := http.Get("http://localhost:3000/orders")
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
