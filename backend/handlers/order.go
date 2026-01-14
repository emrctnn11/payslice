package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emrctnn11/sezzle-payslice-backend/models"
	"github.com/google/uuid"
)

type OrderHandler struct {
	OrderDB   *sql.DB // Postgres
	ProductDB *sql.DB //MYSQL
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the request
	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	txProduct, err := h.ProductDB.Begin()
	if err != nil {
		http.Error(w, "Inventory System error", 500)
		return
	}
	// always rollback if we dont commit explicitly
	defer txProduct.Rollback()

	// check inventory ()
	var inventory int
	var priceCents int64
	// we query standard price and inventory
	err = txProduct.QueryRow("SELECT price_cents, inventory FROM products WHERE id = ?", req.ProductID).Scan(&priceCents, &inventory)
	if err != nil {
		http.Error(w, "Product not found", 404)
		return
	}
	if inventory <= 0 {
		http.Error(w, "Product Out of Stock", 409)
		return
	}
	// Decrement Inventory
	if _, err := txProduct.Exec("UPDATE products SET inventory = inventory - 1 WHERE id = ?", req.ProductID); err != nil {
		http.Error(w, "Update Failed", 500)
		return
	}

	//logic create the order object.
	// TODO: we should listen the order from product.
	txOrder, err := h.OrderDB.Begin()
	if err != nil {
		http.Error(w, "Ledge Error", 500)
		return
	}
	defer txOrder.Rollback()
	// insert into postgres
	orderID := uuid.New().String()
	_, err = txOrder.Exec(`INSERT INTO orders (id, user_id, total_amount_cents, status) VALUES ($1, $2, $3, 'PENDING')`,
		orderID, req.UserID, priceCents)
	if err != nil {
		log.Println("Order Insert Failed: ", err)
		http.Error(w, "Order Failed", 500)
		return
	}

	// calculate & insert 4 installments
	splits := splitAmount(priceCents)
	for i, amount := range splits {
		installmentID := uuid.New().String()
		// Due date: Today, +2 weeks, +4 weeks, +6 weeks
		dueDate := time.Now().AddDate(0, 0, i*14)

		_, err := txOrder.Exec(`INSERT INTO installments (id, order_id, amount_cents, due_date, status) VALUES ($1, $2, $3, $4, 'PENDING')`, installmentID, orderID, amount, dueDate)
		if err != nil {
			log.Println("Installment Insert Failed:", err)
			http.Error(w, "Installment Failed", 500)
			return
		}
	}

	// commit mysql transaction
	if err := txOrder.Commit(); err != nil {
		http.Error(w, "Transaction Commit Failed", 500)
		return
	}

	if err := txProduct.Commit(); err != nil {
		log.Println(w, "CRITICAL: INVENTORY NOT DECREMENTED but Order Created")
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": orderID, "status": "CREATED"})
}

func splitAmount(totalCents int64) []int64 {
	splits := make([]int64, 4)
	baseAmount := totalCents / 4
	remainder := totalCents % 4

	for i := 0; i < 4; i++ {
		splits[i] = baseAmount
	}

	// Distribute the remainder pennies to the first installments
	for i := 0; i < int(remainder); i++ {
		splits[i]++
	}

	return splits
}
