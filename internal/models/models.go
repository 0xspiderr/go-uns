package models

import "encoding/json"

// unified namespace data from the IT and OT level
type UNSData struct {
	TimeStamp string          `json:"timestamp"`
	Topic     string          `json:"topic"`
	OT        json.RawMessage `json:"ot_data"`
	IT        ITData          `json:"it_data"`
}

type Order struct {
	OrderID     int    `json:"id,int"`
	OrderName   string `json:"order_name"`
	Status      string `json:"order_status"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	//	ActiveOperators int    `json:"active_operators"`
}

type ITData []Order
