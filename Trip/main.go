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

var trips Trip
var db *sql.DB

type Trip struct {
	TripID       string
	TripStatus   string //(Processing, Ongoing, End)
	PassengerID  string
	DriverID     string
	PickUpPoint  string
	DropOffPoint string
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////							Functions for MySQL Database										////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func CreateNewTrip(db *sql.DB, t Trip) {
	query := fmt.Sprintf("INSERT INTO Trips VALUES ('%s', '%s', '%s', '%s', '%s', '%s')",
		t.TripID, t.TripStatus, t.PassengerID, t.DriverID, t.PickUpPoint, t.DropOffPoint)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func GetTripInfo(db *sql.DB, TripID string) []Trip {
	query := fmt.Sprintf("SELECT * FROM Trips where TripID = '%s'", TripID)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	var trips []Trip
	for results.Next() {
		var newTrip Trip
		err = results.Scan(&newTrip.TripID, &newTrip.TripStatus, &newTrip.PassengerID, &newTrip.DriverID, &newTrip.PickUpPoint, &newTrip.DropOffPoint)
		if err != nil {
			panic(err.Error())
		}

		trips = append(trips, newTrip)
	}
	return trips
}

func UpdateTripInfo(db *sql.DB, t Trip) {
	query := fmt.Sprintf("UPDATE Trips SET TripStatus='%s', PassengerID='%s', DriverID='%s', PickUpPoint='%s', DropOffPoint='%s' WHERE TripID='%s'",
		t.TripStatus, t.PassengerID, t.DriverID, t.PickUpPoint, t.DropOffPoint, t.TripID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for HTTP											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func trip(w http.ResponseWriter, r *http.Request) {

	var TripID string

	if r.Header.Get("Content-type") == "application/json" {

		//Create a new trip
		if r.Method == "POST" {

			// read the string sent to the service
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newTrip)

				if newTrip.TripID == "" || newTrip.PassengerID == "" || newTrip.DriverID == "" || newTrip.PickUpPoint == "" || newTrip.DropOffPoint == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply trip " + "information " + "in JSON format"))
					return
				} else {
					//fmt.Println("newTrip: ", newTrip)
					newTrip.TripStatus = "Processing"
					CreateNewTrip(db, newTrip)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Successfully created trip"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply trip information " +
					"in JSON format"))
			}
		}
		if r.Method == "PUT" {
			queryParams := r.URL.Query()
			TripID = queryParams["TripID"][0]
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &trips)

				if trips.TripStatus == "" || trips.DriverID == "" || trips.PassengerID == "" || trips.PickUpPoint == "" || trips.DropOffPoint == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply trip " + "information " + "in JSON format"))
					return
				} else {
					trips.TripID = TripID
					UpdateTripInfo(db, trips)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Successfully updated trip's information"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply trip information " +
					"in JSON format"))
			}
		}
	}
	//Get trip's information for passenger
	if r.Method == "GET" { //its working
		DriverID := r.URL.Query().Get("DriverID")
		fmt.Println("DriverID: ", DriverID)
		trips := GetTripInfo(db, DriverID)

		json.NewEncoder(w).Encode(&trips)
	}
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - For audit purposes, trip's information cannot be deleted."))
	}
}

func main() {
	// instantiate trips
	fmt.Println("GonGrab MySQL!")
	ridesharing_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ridesharing_db")

	db = ridesharing_db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/trips", trip).Methods(
		"GET", "PUT", "POST", "DELETE")
	fmt.Println("Trip microservice API --> Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))

	defer db.Close()
}
