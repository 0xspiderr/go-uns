package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type ITData struct {
	OrderID         int         `json:"order_id"`
	OrderName       string      `json:"order_name"`
	Status          OrderStatus `json:"order_status"`
	ProductName     string      `json:"product_name"`
	Quantity        int         `json:"quantity"`
	ActiveOperators int         `json:"active_operators"`
}

func fetchITData() (ITData, error) {
	var it ITData
	response, err := http.Get("http://localhost:8089/api/orders")
	if err != nil {
		return it, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&it)
	return it, err
}

func itDataHandler(w http.ResponseWriter, r *http.Request) {
	data := ITData{
		OrderID:         1,
		OrderName:       "mock_order_name",
		Status:          InProgress,
		ProductName:     "mock_product_name",
		Quantity:        10,
		ActiveOperators: 2,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func startITServer() {
	http.HandleFunc("/api/orders", itDataHandler)
	fmt.Println("Started HTTP server running on http://localhost:8089/api/orders")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatalf("Starting HTTP server failed: %v", err)
	}
}
