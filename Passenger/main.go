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

//Global variable for microservice passenger
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
//Registering new passenger
func CreateNewPassenger(db *sql.DB, p Passenger) {
	query := fmt.Sprintf("INSERT INTO Passengers VALUES ('%s', '%s', '%s', '%s', '%s')",
		p.PassengerID, p.FirstName, p.LastName, p.MobileNo, p.EmailAdd)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

//Updating existing passenger information
func UpdatePassengerInfo(db *sql.DB, p Passenger) {
	query := fmt.Sprintf("UPDATE Passengers SET FirstName='%s', LastName='%s', MobileNo='%s', EmailAdd='%s' WHERE PassengerID='%s'",
		p.FirstName, p.LastName, p.MobileNo, p.EmailAdd, p.PassengerID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//Passenger using mobile phone number to login
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

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for HTTP											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func passenger(w http.ResponseWriter, r *http.Request) {

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
				//Check if user fill up the required information for registering Passenger's account
				if newPassenger.PassengerID == "" || newPassenger.FirstName == "" || newPassenger.LastName == "" || newPassenger.MobileNo == "" || newPassenger.EmailAdd == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + "information " + "in JSON format"))
					return
				} else {
					CreateNewPassenger(db, newPassenger) //Once everything is checked, passenger's account will be created
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
		//---PUT is for creating or updating existing passenger---
		if r.Method == "PUT" {
			queryParams := r.URL.Query() //used to resolve the conflict of calling API using the '%s'?PassengerID='%s' method
			PassengerID = queryParams["PassengerID"][0]
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &passengers)
				//Check if user fill up the required information for updating Passenger's account information
				if passengers.FirstName == "" || passengers.LastName == "" || passengers.MobileNo == "" || passengers.EmailAdd == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply passenger " + " information " + "in JSON format"))
				} else {
					passengers.PassengerID = PassengerID
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
	//---GET is used to get existing passenger's information such as mobile phone number to login
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
	//---Deny any deletion of passenger's account or other passenger's information
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - For audit purposes, passenger's account cannot be deleted."))
	}
}

func main() {
	// instantiate passengers
	ridesharing_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ridesharing_db")

	db = ridesharing_db
	// handle error
	if err != nil {
		panic(err.Error())
	}
	//handle the API connection across all three microservices, Passengers, Trips and Drivers
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
