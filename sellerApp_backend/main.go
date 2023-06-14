package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type OrderPayload struct {
	ID           string      `json:"id"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items"`
	Total        float64     `json:"total"`
	CurrencyUnit string      `json:"currencyUnit"`
}

type OrderItem struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:15761070aA@@tcp(127.0.0.1:3306)/sys")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func storeOrder(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	var orders []OrderPayload
	err := json.NewDecoder(r.Body).Decode(&orders)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received orders:", orders)

	db, err := getDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, order := range orders {
		stmt, err := tx.Prepare("INSERT INTO orders (id, status, total, currency_unit) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(order.ID, order.Status, order.Total, order.CurrencyUnit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, item := range order.Items {
			stmt, err = tx.Prepare("INSERT INTO order_items (order_id, item_id, description, price, quantity) VALUES (?, ?, ?, ?, ?)")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(order.ID, item.ID, item.Description, item.Price, item.Quantity)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Order(s) stored successfully")
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	orderID := r.URL.Query().Get("id")

	db, err := getDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, status, total, currency_unit FROM orders WHERE id = ?", orderID)

	var order OrderPayload
	err = row.Scan(&order.ID, &order.Status, &order.Total, &order.CurrencyUnit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT item_id, description, price, quantity FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItem
		err = rows.Scan(&item.ID, &item.Description, &item.Price, &item.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		order.Items = append(order.Items, item)
	}

	jsonResponse, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sortBy := r.URL.Query().Get("sort")
	filterStatus := r.URL.Query().Get("status")
	filterDescription := r.URL.Query().Get("description")
	filterPrice := r.URL.Query().Get("price")
	filterQuantity := r.URL.Query().Get("quantity")

	// Set default values if not provided
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	// Convert page and limit to integers
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	// Calculate the offset based on page and limit
	offset := (pageNum - 1) * limitNum

	db, err := getDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Build the SQL query
	query := "SELECT o.id, o.status, o.total, o.currency_unit, i.item_id, i.description, i.price, i.quantity FROM orders o JOIN order_items i ON o.id = i.order_id WHERE 1=1"
	params := []interface{}{}

	// Add filters
	if filterStatus != "" {
		query += " AND o.status = ?"
		params = append(params, filterStatus)
	}
	if filterDescription != "" {
		query += " AND i.description = ?"
		params = append(params, filterDescription)
	}
	if filterPrice != "" {
		query += " AND i.price = ?"
		params = append(params, filterPrice)
	}
	if filterQuantity != "" {
		query += " AND i.quantity = ?"
		params = append(params, filterQuantity)
	}

	// Add sorting
	if sortBy != "" {
		query += " ORDER BY " + sortBy
	}

	// Add the pagination
	query += " LIMIT ? OFFSET ?"
	params = append(params, limitNum, offset)

	rows, err := db.Query(query, params...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []OrderPayload
	for rows.Next() {
		var order OrderPayload
		var item OrderItem
		err = rows.Scan(&order.ID, &order.Status, &order.Total, &order.CurrencyUnit, &item.ID, &item.Description, &item.Price, &item.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the order already exists in the orders slice
		existingOrder := findOrderByID(orders, order.ID)
		if existingOrder == nil {
			// If not, add it to the orders slice
			order.Items = append(order.Items, item)
			orders = append(orders, order)
		} else {
			// If yes, add the item to the existing order's items slice
			existingOrder.Items = append(existingOrder.Items, item)
		}
	}

	jsonResponse, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func findOrderByID(orders []OrderPayload, orderID string) *OrderPayload {
	for i := range orders {
		if orders[i].ID == orderID {
			return &orders[i]
		}
	}
	return nil
}

type UpdateOrderStatusPayload struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	var payload UpdateOrderStatusPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE orders SET status = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(payload.Status, payload.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Order status updated successfully")
}

func main() {
	http.HandleFunc("/storeOrder", storeOrder)
	http.HandleFunc("/getOrder", getOrder)
	http.HandleFunc("/getAllOrders", getAllOrders)
	http.HandleFunc("/updateOrderStatus", updateOrderStatus)

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
