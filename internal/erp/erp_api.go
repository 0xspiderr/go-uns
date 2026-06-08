package erp

import (
	"encoding/json"
	"log"
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
	// The ip address of the ERP and endpoint to fetch the order data
	response, err := erpClient.Get("http://25.10.202.120:8083/orders")
	if err != nil {
		return orders, err
	}
	// debug response
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&orders)

	return orders, err
}

func HandleNewOrder(w http.ResponseWriter, r *http.Request) {
	// ensure POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, ERP should send POST", http.StatusMethodNotAllowed)
		return
	}

	// decode incoming order
	var newOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		log.Printf("Failed to decode ERP payload: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received new order from the ERP: [%s] for %d units of %s", newOrder.OrderName, newOrder.Quantity, newOrder.ProductName)

	// send response back to ERP so they know we got the order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "order successfully received by the UNS",
		"status":  "received",
	})
}

// listening server for the ERP
func StartServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders", HandleNewOrder)

	log.Printf("Starting UNS HTTP Server on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("HTTP server crashed: %v", err)
	}
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
