package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_golang_rest")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", createCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomerByID).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("server is starting at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	rows, err := db.Query("SELECT * FROM customers")

	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	var customers []customer

	for rows.Next() {
		var customer customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address)

		if err != nil {
			panic(err.Error())
		}

		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO customers(name, email, phone, address) VALUES(?, ?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error)
	}

	keyVal := make(map[string]interface{})
	json.Unmarshal(body, &keyVal)

	name := keyVal["name"]
	email := keyVal["email"]
	phone := keyVal["phone"]
	address := keyVal["address"]

	_, err = stmt.Exec(name, email, phone, address)

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New customer has been added.")
}

func getCustomerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM customers WHERE id = ?", params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var customer customer

	for result.Next() {
		err := result.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address)

		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(customer)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE customers SET name = ?, email = ?, phone = ?, address = ? WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]interface{})
	json.Unmarshal(body, &keyVal)

	name := keyVal["name"]
	email := keyVal["email"]
	phone := keyVal["phone"]
	address := keyVal["address"]

	_, err = stmt.Exec(name, email, phone, address, params["id"])

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Customer with ID %s has been updated.", params["id"])
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM customers WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Customer with ID %s has been deleted.", params["id"])
}
