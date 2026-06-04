package models

import "encoding/json"

// unified namespace data from the IT and OT level
type UNSData struct {
	TimeStamp string          `json:"timestamp"`
	Topic     string          `json:"topic"`
	OT        json.RawMessage `json:"ot_data"`
	IT        Order           `json:"it_data"`
}

type Order struct {
	OrderID      int    `json:"id"`
	OrderName    string `json:"name"`    // Changed from "order_name"
	Status       string `json:"status"`  // Changed from "order_status"
	ProductName  string `json:"product"` // Changed from "product_name"
	Quantity     int    `json:"quantity"`
	CreatedAt    string `json:"createdAt"`    // New field
	ConveyorID   int    `json:"conveyorId"`   // New field
	ConveyorName string `json:"conveyorName"` // New field
}

type ITData []Order
