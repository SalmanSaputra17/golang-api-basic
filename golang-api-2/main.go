package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var dbErr error

type booking struct {
	gorm.Model
	User    string `json:"user"`
	Members int    `json:"members"`
}

func main() {
	db, dbErr = gorm.Open("mysql", "root:@/db_golang_rest?charset=utf8&parseTime=true&loc=Local")

	if dbErr != nil {
		log.Println("Cannot open a connection")
	} else {
		log.Println("Server established")
	}

	defer db.Close()

	db.AutoMigrate(&booking{})
	handleRequests()
}

func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:8080")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/bookings", getBookings).Methods("GET")
	router.HandleFunc("/bookings", createBooking).Methods("POST")
	router.HandleFunc("/bookings/{id}", getBookingByID).Methods("GET")
	router.HandleFunc("/bookings/{id}", updateBooking).Methods("PUT")
	router.HandleFunc("/bookings/{id}", deleteBooking).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bookings []booking
	db.Find(&bookings)

	json.NewEncoder(w).Encode(bookings)
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	var booking booking
	json.Unmarshal(body, &booking)

	db.Create(&booking)

	fmt.Fprintf(w, "New booking has been created.")
}

func getBookingByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var booking booking
	db.First(&booking, params["id"])

	json.NewEncoder(w).Encode(booking)
}

func updateBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	var booking booking
	var keyVal = make(map[string]interface{})

	db.First(&booking, params["id"])
	json.Unmarshal(body, &keyVal)

	db.Model(&booking).Update(keyVal)

	fmt.Fprintf(w, "Booking with ID %s has been updated.", params["id"])
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var booking booking
	db.First(&booking, params["id"])
	db.Delete(&booking)

	fmt.Fprintf(w, "Booking with ID %s has been deleted.", params["id"])
}
