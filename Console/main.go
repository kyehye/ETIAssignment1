package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for HTTP											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func httpPost(url string, data interface{}) (*http.Response, error) {
	jsonData, _ := json.Marshal(data)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return response, err
}

func httpPut(url string, data interface{}) (*http.Response, error) {
	jsonData, _ := json.Marshal(data)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	return response, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

const driverUrl = "http://localhost:5000/drivers"
const driverAvailabilityUrl = "http://localhost:5000/drivers/available"
const passengerUrl = "http://localhost:5001/passengers"
const tripUrl = "http://localhost:5002/trips"

type Passenger struct {
	PassengerID string
	FirstName   string
	LastName    string
	MobileNo    string
	EmailAdd    string
}

type Driver struct {
	DriverID     string
	FirstName    string
	LastName     string
	MobileNo     string
	EmailAdd     string
	CarLicenseNo string
	Availability bool
}

type Trip struct {
	TripID       string
	TripStatus   string
	PassengerID  string
	DriverID     string
	PickUpPoint  string
	DropOffPoint string
}

func main() {
	MainMenu()
}

func MainMenu() {
	for {
		fmt.Println("Welcome to GrabnGo!")
		fmt.Println("[1] Login as Passenger")
		fmt.Println("[2] Login as Driver")
		fmt.Println("[3] Register Passenger Account")
		fmt.Println("[4] Register Driver Account")
		fmt.Println("[0] Quit")

		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			passengerLogin()
		} else if userInput == "2" {
			driverLogin()
		} else if userInput == "3" {
			passengerRegister()
		} else if userInput == "4" {
			driverRegister()
		} else if userInput == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			break
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////[1] Passenger/////////////////////////////////////////////////																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func passengerLogin() {
	fmt.Print("Please enter your mobile number: ")
	var mobileno string
	fmt.Scanln(&mobileno)

	passenger := getPassengerMobileNo(mobileno)
	if mobileno != passenger.MobileNo {
		fmt.Println("\nInvalid mobile number. Login failed.")
	} else {
		fmt.Printf("\nWelcome to GrabnGo Passenger, %s %s!\n", passenger.FirstName, passenger.LastName)
		passengerMainMenu(passenger)
	}
}

func passengerRegister() {
	fmt.Print("Passenger ID: ")
	var passengerid string
	fmt.Scanln(&passengerid)

	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobileno string
	fmt.Scanln(&mobileno)

	fmt.Print("Email Address: ")
	var emailadd string
	fmt.Scanln(&emailadd)

	newPassenger := Passenger{
		PassengerID: passengerid,
		FirstName:   firstname,
		LastName:    lastname,
		MobileNo:    mobileno,
		EmailAdd:    emailadd,
	}

	err := CreateNewPassenger(newPassenger)
	if err != nil {
		fmt.Println("Passenger's account creation fail")
	} else {
		fmt.Println("Successfully created passenger's account")
	}
}

func passengerMainMenu(passenger Passenger) {
	for {
		fmt.Println("[1] Book trip")
		fmt.Println("[2] View trip's history")
		fmt.Println("[3] Update personal information")
		fmt.Println("[0] Logout")

		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			passengerBookTrip(passenger)
		} else if userInput == "2" {
			passengerTrips(passenger)
		} else if userInput == "3" {
			passengerUpdateInformation(passenger)
		} else if userInput == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			passengerMainMenu(passenger)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
///				  					Passenger's main menu features									 ///
////////////////////////////////////////////////////////////////////////////////////////////////////////
//Passenger book trip
func passengerBookTrip(passenger Passenger) {
	fmt.Print("Trip ID: ")
	var tripid string
	fmt.Scanln(&tripid)

	fmt.Print("Pick Up Point(Postal Code): ")
	var pickuppoint string
	fmt.Scanln(&pickuppoint)

	fmt.Print("Drop Off Point(Postal Code): ")
	var dropoffpoint string
	fmt.Scanln(&dropoffpoint)

	driver := getAvailableDriver()
	if driver.Availability != true {
		fmt.Println("Drivers are not available at the current moment")
	}

	err := CreateNewTrip(tripid, driver.DriverID, passenger.PassengerID, pickuppoint, dropoffpoint)
	if err != nil {
		fmt.Println("Trip booking failed")
	} else {
		fmt.Println("Trip booked successfully!")
	}

	driver.Availability = false
	UpdateDriverInfo(driver)
}

//View passenger trip
func passengerTrips(passenger Passenger) {
	trips := viewPassengerTrips(passenger.PassengerID)

	for x := len(trips) - 1; x >= 0; x-- {
		trip := trips[x]
		fmt.Println()
		fmt.Println("Trip ID: ", trip.TripID)
		fmt.Println("Driver ID: ", trip.DriverID)
		fmt.Println("Passenger ID: ", trip.PassengerID)
		fmt.Println("Pick Up Postal Code: ", trip.PickUpPoint)
		fmt.Println("Drop Off Postal Code: ", trip.DropOffPoint)
		fmt.Println()
	}
	/*
		x := len(trips)
		for x >= 0 {
			trip := trips[x]
			fmt.Println("Trip ID: ", trip.TripID)
			fmt.Println("Driver ID: ", trip.DriverID)
			fmt.Println("Passenger ID: ", trip.PassengerID)
			fmt.Println("Pick Up Postal Code: ", trip.PickUpPoint)
			fmt.Println("Drop Off Postal Code: ", trip.DropOffPoint)
			fmt.Println()
		}
	*/
}

func viewPassengerTrips(PassengerID string) []Trip {
	var trips []Trip

	url := fmt.Sprintf("%s?PassengerID=%s", tripUrl, PassengerID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No trips found")
		return trips
	}

	json.NewDecoder(resp.Body).Decode(&trips)
	return trips
}

//Update passenger's information
func passengerUpdateInformation(passenger Passenger) {
	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobileno string
	fmt.Scanln(&mobileno)

	fmt.Print("Email Address: ")
	var emailadd string
	fmt.Scanln(&emailadd)

	newPassengerInfo := Passenger{
		PassengerID: passenger.PassengerID,
		FirstName:   firstname,
		LastName:    lastname,
		MobileNo:    mobileno,
		EmailAdd:    emailadd,
	}
	err := UpdatePassengerInfo(newPassengerInfo)
	if err != nil {
		fmt.Println("Passenger information update failed")
	} else {
		fmt.Println("Successfully updated passenger's information")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////Driver////////////////////////////////////////////////////																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func driverLogin() {
	for {
		fmt.Print("Please enter your mobile number: ")
		var mobileno string
		fmt.Scanln(&mobileno)

		driver := getDriverMobileNo(mobileno)
		if mobileno != driver.MobileNo {
			fmt.Println("\nInvalid mobile number. Login failed.")
			break
		} else {
			fmt.Printf("\nWelcome to GrabnGo Driver, %s %s!\n", driver.FirstName, driver.LastName)
			DriverMainMenu(driver)
			break
		}
	}
}

func driverRegister() {
	fmt.Print("Driver ID: ")
	var driverid string
	fmt.Scanln(&driverid)

	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobileno string
	fmt.Scanln(&mobileno)

	fmt.Print("Email Address: ")
	var emailadd string
	fmt.Scanln(&emailadd)

	fmt.Print("Car License Number: ")
	var carlicenseno string
	fmt.Scanln(&carlicenseno)

	newDriver := Driver{
		DriverID:     driverid,
		FirstName:    firstname,
		LastName:     lastname,
		MobileNo:     mobileno,
		EmailAdd:     emailadd,
		CarLicenseNo: carlicenseno,
	}
	//Set driver availability to true, when driver account is first created.
	newDriver.Availability = true
	err := CreateNewDriver(newDriver)
	if err != nil {
		fmt.Println("Driver's account creation fail")
	} else {
		fmt.Println("Successfully created driver's account")
	}
}

func DriverMainMenu(driver Driver) {
	for {
		fmt.Println("[1] Start trip")
		fmt.Println("[2] End trip")
		fmt.Println("[3] Update personal information")
		fmt.Println("[0] Logout")

		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			//driverStartTrip()
		} else if userInput == "2" {
			//driverEndTrip()
		} else if userInput == "3" {
			driverUpdateInformation(driver)
		} else if userInput == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			DriverMainMenu(driver)
		}
	}
}

func driverUpdateInformation(driver Driver) {
	fmt.Print("First Name: ")
	var firstname string
	fmt.Scanln(&firstname)

	fmt.Print("Last Name: ")
	var lastname string
	fmt.Scanln(&lastname)

	fmt.Print("Mobile Number: ")
	var mobileno string
	fmt.Scanln(&mobileno)

	fmt.Print("Email Address: ")
	var emailadd string
	fmt.Scanln(&emailadd)

	fmt.Print("Car License Number: ")
	var carlicenseno string
	fmt.Scanln(&carlicenseno)

	newDriverInfo := Driver{
		DriverID:     driver.DriverID,
		FirstName:    firstname,
		LastName:     lastname,
		MobileNo:     mobileno,
		EmailAdd:     emailadd,
		CarLicenseNo: carlicenseno,
	}
	err := UpdateDriverInfo(newDriverInfo)
	if err != nil {
		fmt.Println("Driver information update failed")
	} else {
		fmt.Println("Successfully updated driver's information")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
////																								////
////									Functions for API											////
////																								////
////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////Passenger//////////////////////////////////////////////////																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////

func CreateNewPassenger(newPassenger Passenger) error {
	url := passengerUrl
	_, err := httpPost(url, newPassenger)
	return err
}

func UpdatePassengerInfo(newPassengerInfo Passenger) error {
	url := fmt.Sprintf("%s?PassengerID=%s", passengerUrl, newPassengerInfo.PassengerID)

	_, err := httpPut(url, newPassengerInfo)
	return err
}

func getPassengerMobileNo(MobileNo string) Passenger {
	var passenger Passenger

	url := fmt.Sprintf("%s?MobileNo=%s", passengerUrl, MobileNo)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return passenger
	}

	json.NewDecoder(resp.Body).Decode(&passenger)
	return passenger
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////Driver////////////////////////////////////////////////////																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////

func CreateNewDriver(newDriver Driver) error {
	url := driverUrl
	_, err := httpPost(url, newDriver)
	return err
}

func UpdateDriverInfo(newDriverInfo Driver) error {
	url := fmt.Sprintf("%s?DriverID=%s", driverUrl, newDriverInfo.DriverID)

	_, err := httpPut(url, newDriverInfo)
	return err
}

func getDriverMobileNo(MobileNo string) Driver {
	var driver Driver

	url := fmt.Sprintf("%s?MobileNo=%s", driverUrl, MobileNo)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return driver
	}

	json.NewDecoder(resp.Body).Decode(&driver)
	return driver
}

//rip
func getAvailableDriver() Driver {
	var driver Driver

	url := fmt.Sprintf("%s/available?Availability=%t", driverUrl, true)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return driver
	}

	json.NewDecoder(resp.Body).Decode(&driver)
	return driver
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////Trip/////////////////////////////////////////////////////																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////

func CreateNewTrip(tripid string, driverid string, passengerid string, pickuppoint string, dropoffpoint string) error {

	url := tripUrl

	var newTrip Trip = Trip{
		TripID:       tripid,
		DriverID:     driverid,
		PassengerID:  passengerid,
		PickUpPoint:  pickuppoint,
		DropOffPoint: dropoffpoint,
	}

	_, err := httpPost(url, newTrip)
	return err
}
