package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var passengers Passenger
var db *sql.DB

type Passenger struct { // map this type to the record in the table
	PassengerID string
	FirstName   string
	LastName    string
	MobileNo    string
	EmailAdd    string
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////							Functions for MySQL Database										////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
//Edit existing passenger information in DB
func CreateNewPassenger(db *sql.DB, p Passenger) {
	query := fmt.Sprintf("INSERT INTO Passengers VALUES ('%s', '%s', '%s', '%s', '%s')",
		p.PassengerID, p.FirstName, p.LastName, p.MobileNo, p.EmailAdd)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func UpdatePassengerInfo(db *sql.DB, p Passenger) {
	query := fmt.Sprintf("UPDATE Passengers SET FirstName='%s', LastName='%s', MobileNo='%s', EmailAdd='%s' WHERE PassengerID='%s'",
		p.FirstName, p.LastName, p.MobileNo, p.EmailAdd, p.PassengerID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func PassengerLogin(db *sql.DB, MobileNo string) (Passenger, string) {
	query := fmt.Sprintf("SELECT * FROM Passengers WHERE MobileNo = '%s'", MobileNo)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&passengers.PassengerID, &passengers.FirstName, &passengers.LastName, &passengers.MobileNo, &passengers.EmailAdd); err {
	case sql.ErrNoRows:
		errMsg = "Mobile number not found. Passenger login failed."
	case nil:
	default:
		panic(err.Error())
	}

	return passengers, errMsg
}

//Get passenger's information
func GetPassengerInfo(db *sql.DB) {
	results, err := db.Query("Select * FROM ridesharing_db.Passengers")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var passenger Passenger
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName,
			&passenger.LastName, &passenger.MobileNo, &passenger.EmailAdd)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(passenger.PassengerID, passenger.FirstName,
			passenger.LastName, passenger.MobileNo, passenger.EmailAdd)
	}
}

//Get passenger's trip history record (Should this be done in trips?)
func GetPassengerTripRecords(db *sql.DB) {
	results, err := db.Query("Select * FROM ridesharing_db.Passengers")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var passenger Passenger
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName,
			&passenger.LastName, &passenger.MobileNo, &passenger.EmailAdd)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(passenger.PassengerID, passenger.FirstName,
			passenger.LastName, passenger.MobileNo, passenger.EmailAdd)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for HTTP											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func passenger(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var PassengerID string

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new passenger
		if r.Method == "POST" {
			// read the string sent to the service
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.PassengerID == "" || newPassenger.FirstName == "" || newPassenger.LastName == "" || newPassenger.MobileNo == "" || newPassenger.EmailAdd == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
					return
				} else {
					CreateNewPassenger(db, newPassenger)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Successfully created passenger's account"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " +
					"in JSON format"))
			}
		}
		//---PUT is for creating or updating
		// existing passenger's info---
		if r.Method == "PUT" {

			fmt.Sscan(params["PassengerID"], &PassengerID)
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &passengers)

				if passengers.FirstName == "" || passengers.LastName == "" || passengers.MobileNo == "" || passengers.EmailAdd == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + " information " + "in JSON format"))
				} else {
					UpdatePassengerInfo(db, passengers)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Successfully updated passenger's information"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " + "passenger information " + "in JSON format"))
			}
		}

	}
	//Get passenger's information
	if r.Method == "GET" { //its working
		MobileNo := r.URL.Query().Get("MobileNo")
		fmt.Println("MobileNo: ", MobileNo)
		passengers, errMsg := PassengerLogin(db, MobileNo)

		if errMsg == "Mobile number not found. Passenger login failed." {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Passenger's account not found"))
		} else {
			fmt.Println(passengers)
			json.NewEncoder(w).Encode(&passengers)
		}
	}
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - For audit purposes, passenger account cannot be deleted."))
	}
}

/*
func passengertriprequest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") == "application/json" {
		// POST is for passenger to request for a trip
		if r.Method == "PUT" {
			// read the string sent to the service
			var newPassengerTripRequest Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassengerTripRequest)

				if newPassengerTripRequest.PickUpPoint == "" || newPassengerTripRequest.DropOffPoint == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
					return
				} else {
					PassengerRequestTrip(db, passengers)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("201 - Successfully updated passenger's trip"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " +
					"in JSON format"))
			}
		}
	}
}
*/

func main() {
	// instantiate passengers
	ridesharing_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ridesharing_db")

	db = ridesharing_db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()

	router.HandleFunc("/passengers", passenger).Methods(
		"GET", "POST", "PUT", "DELETE")
	fmt.Println("Passenger microservice API --> Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))

	defer db.Close()
}

//To add Passenger,
//Code: curl -H "Content-Type:application/json" -X POST http://localhost:5001/api/v1/passengers/2 -d "{\"id\":\"2\", \"firstname\":\"Sha Li\", \"lastname\":\"Kang\", \"mobileno\":\"81234567\", \"emailadd\":\"shali@gmail.com\"}"
//Field: ID, FirstName, LastName, MobileNo, EmailAdd
