package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// handler package import
	"github.com/emrctnn11/sezzle-payslice-backend/handlers"
	"github.com/emrctnn11/sezzle-payslice-backend/middleware"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	// 1. Connect to MYSQL
	// Format: username:password@tcp(host:port)/dbname
	mysqlDSN := "sezzle_user:sezzle_password@tcp(localhost:3306)/product_db"
	productDB, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatal("Error opening MySQL: ", err)
	}

	postgresDSN := "user=sezzle_admin password=adminpassword dbname=ledger_db sslmode=disable host=localhost port=5432"
	orderDB, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to verify
	if err := productDB.Ping(); err != nil {
		log.Fatal("Cannot connect to MySQL. Is Docker Running? ", err)
	}
	if err := orderDB.Ping(); err != nil {
		log.Fatal("Postgres Down:", err)
	}
	fmt.Println("Connected to MYSQL(PRODUCT CATALOG) & Postgres")

	ph := &handlers.ProductHandler{DB: productDB}
	oh := &handlers.OrderHandler{OrderDB: orderDB, ProductDB: productDB}

	// defining routes
	http.HandleFunc("/products", middleware.CORSMiddleware(ph.GetProducts))
	http.HandleFunc("/orders", middleware.CORSMiddleware(oh.CreateOrder))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// API Handlers:
