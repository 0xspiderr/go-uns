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

// Get mock data from the HTTP server
func erpDataHandler(w http.ResponseWriter, r *http.Request) {
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

func startERPServer() {
	http.HandleFunc("/api/orders", erpDataHandler)
	fmt.Println("Started HTTP server running on http://localhost:8089/api/orders")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatalf("Starting HTTP server failed: %v", err)
	}
}
