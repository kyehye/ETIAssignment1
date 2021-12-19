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

//Global variable for microservice driver
var drivers Driver
var db *sql.DB

type Driver struct { // map this type to the record in the table
	DriverID     string
	FirstName    string
	LastName     string
	MobileNo     string
	EmailAdd     string
	CarLicenseNo string
	Availability bool
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////							Functions for MySQL Database										////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
// func for Driver's page
//Registering new driver
func CreateNewDriver(db *sql.DB, d Driver) {
	query := fmt.Sprintf("INSERT INTO Drivers VALUES ('%s', '%s', '%s','%s', '%s', '%s', %t)",
		d.DriverID, d.FirstName, d.LastName, d.MobileNo, d.EmailAdd, d.CarLicenseNo, d.Availability)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

//Updating existing driver information
func UpdateDriverInfo(db *sql.DB, d Driver) {
	fmt.Println(d)
	query := fmt.Sprintf("UPDATE Drivers SET FirstName='%s', LastName='%s', MobileNo='%s', EmailAdd='%s', CarLicenseNo='%s', Availability=%t WHERE DriverID='%s'",
		d.FirstName, d.LastName, d.MobileNo, d.EmailAdd, d.CarLicenseNo, d.Availability, d.DriverID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//Driver using mobile phone number to login
func DriverLogin(db *sql.DB, MobileNo string) (Driver, string) {
	query := fmt.Sprintf("SELECT * FROM Drivers WHERE MobileNo = '%s'", MobileNo)

	results := db.QueryRow(query)
	var errMsg string

	switch err := results.Scan(&drivers.DriverID, &drivers.FirstName, &drivers.LastName, &drivers.MobileNo, &drivers.EmailAdd, &drivers.CarLicenseNo, &drivers.Availability); err {
	case sql.ErrNoRows:
		errMsg = "Mobile number not found. Driver login failed."
	case nil:
	default:
		panic(err.Error())
	}

	return drivers, errMsg
}

//Get available driver for trip booking of passenger
func GetAvailableDriver(db *sql.DB) Driver {
	results, err := db.Query("Select * FROM ridesharing_db.Drivers WHERE Availability = true LIMIT 1") //Limit only 1 driver at a time for each trip booking.

	if err != nil {
		panic(err.Error())
	}

	for results.Next() { //Search for an available driver.
		// map this type to the record in the table
		err = results.Scan(&drivers.DriverID, &drivers.FirstName,
			&drivers.LastName, &drivers.MobileNo, &drivers.EmailAdd, &drivers.CarLicenseNo, &drivers.Availability)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(drivers.DriverID, drivers.FirstName,
			drivers.LastName, drivers.MobileNo, drivers.EmailAdd, drivers.CarLicenseNo, drivers.Availability)
	}
	return drivers
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for HTTP											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func driver(w http.ResponseWriter, r *http.Request) {
	var DriverID string

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new driver
		if r.Method == "POST" {
			// read the string sent to the service
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)
				//Check if user fill up the required information for registering Driver's account
				if newDriver.DriverID == "" || newDriver.FirstName == "" || newDriver.LastName == "" || newDriver.MobileNo == "" || newDriver.EmailAdd == "" || newDriver.CarLicenseNo == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver " + "information " + "in JSON format"))
					return
				} else {
					CreateNewDriver(db, newDriver) //Once everything is checked, driver's account will be created
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Successfully created driver's information"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver" + "information" + "in JSON format"))
			}
		}
		//---PUT is for creating or updating existing driver---
		if r.Method == "PUT" {
			//fmt.Sscan(params["DriverID"], &DriverID)
			queryParams := r.URL.Query() //used to resolve the conflict of calling API using the '%s'?Driver='%s' method
			DriverID = queryParams["DriverID"][0]
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &drivers)
				if drivers.FirstName == "" || drivers.LastName == "" || drivers.MobileNo == "" || drivers.EmailAdd == "" || drivers.CarLicenseNo == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply driver " + " information " + "in JSON format"))
				} else {
					drivers.DriverID = DriverID
					UpdateDriverInfo(db, drivers)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Successfully updated driver's information"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " + "driver information " + "in JSON format"))
			}
		}
	}
	//---GET is used to get existing driver's information such as mobile phone number to login
	if r.Method == "GET" {
		MobileNo := r.URL.Query().Get("MobileNo")
		fmt.Println("MobileNo: ", MobileNo)
		drivers, errMsg := DriverLogin(db, MobileNo)

		if errMsg == "Mobile number not found. Driver login failed." {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Driver's account not found"))
		} else {
			fmt.Println(drivers)
			json.NewEncoder(w).Encode(&drivers)
		}
	}
	//---Deny any deletion of driver's account or other driver's information
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - For audit purposes, driver's account cannot be deleted."))
	}
}

//---GET is used retrieve available driver to assign to trip based on driver.Availability
func getavailabledriver(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		driver := GetAvailableDriver(db)
		fmt.Println(driver)
		json.NewEncoder(w).Encode(driver)
	}
}

//Assign driver to trip
func main() {
	// instantiate drivers
	ridesharing_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ridesharing_db")

	db = ridesharing_db
	// handle error
	if err != nil {
		panic(err.Error())
	}
	//handle the API connection across all three microservices, Passengers, Trips and Drivers
	router := mux.NewRouter()
	router.HandleFunc("/drivers", driver).Methods(
		"GET", "PUT", "POST", "DELETE")
	router.HandleFunc("/drivers/available", getavailabledriver).Methods(
		"GET")
	fmt.Println("Driver microservice API --> Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

	defer db.Close()
}

//To add Driver,
//Code: curl -H "Content-Type:application/json" -X POST http://localhost:5000/api/v1/drivers/S8392010D -d "{\"id\":\"S8392010D\", \"firstname\":\"Sha Li\", \"lastname\":\"Kang\", \"mobileno\":\"81234567\", \"emailadd\":\"shali@gmail.com\", \"carlicenseno\":\"S1283930C\"}"
//Field: ID, FirstName, LastName, MobileNo, EmailAdd, CarLicenseNo
