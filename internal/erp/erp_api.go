package erp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/0xspiderr/go-uns/internal/models"
)

// might remove this
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

func FetchITData() ([]models.Order, error) {
	var orders []models.Order
	response, err := erpClient.Get("http://192.168.1.13:8083/orders")
	if err != nil {
		return orders, err
	}
	// debug response
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&orders)

	return orders, err
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
