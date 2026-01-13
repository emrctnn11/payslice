package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// module Path to get product struct
	"github.com/emrctnn11/sezzle-payslice-backend/models"
)

type ProductHandler struct {
	DB *sql.DB
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query("SELECT id, name, price_cents, inventory FROM products")
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.PriceCents, &p.Inventory); err != nil {
			continue
		}
		products = append(products, p)
	}

	json.NewEncoder(w).Encode(products)
}
