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
const tripGetDriverIDUrl = "http://localhost:5002/trips/driver"

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

//calling for Main Menu
func main() {
	MainMenu()
}

//Main menu when user launch GrabnGo's Console Application
func MainMenu() {
	for {
		fmt.Println("Welcome to GrabnGo!")
		fmt.Println("[1] Login as Passenger")
		fmt.Println("[2] Login as Driver")
		fmt.Println("[3] Register Passenger Account")
		fmt.Println("[4] Register Driver Account")
		fmt.Println("[0] Quit")

		fmt.Print("\nEnter your option: ")
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
//						Main menu features for user accessing GrabnGo as Passenger					  //																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func passengerLogin() {
	fmt.Print("Please enter your mobile number: ") //Passenger login using their mobile number
	var mobileno string
	fmt.Scanln(&mobileno)

	//Look for registered driver with the given mobile number
	passenger := getPassengerMobileNo(mobileno)
	if mobileno != passenger.MobileNo {
		fmt.Println("\nInvalid mobile number. Login failed.") //If mobile number does not exist, login fail.
	} else {
		fmt.Printf("\nWelcome to GrabnGo Passenger, %s %s!\n", passenger.FirstName, passenger.LastName)
		passengerMainMenu(passenger)
	}
}

func passengerRegister() {
	//Prompt user to enter all the required information to register as a passenger
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

////////////////////////////////////////////////////////////////////////////////////////////////////////
///				  						Passenger's menu features									 ///
////////////////////////////////////////////////////////////////////////////////////////////////////////
func passengerMainMenu(passenger Passenger) {
	for {
		fmt.Println("[1] Book trip")
		fmt.Println("[2] View all trip")
		fmt.Println("[3] Update personal information")
		fmt.Println("[0] Logout")

		fmt.Print("\nEnter your option: ")
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			passengerBookTrip(passenger)
		} else if userInput == "2" {
			allPassengerTrips(passenger)
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
func allPassengerTrips(passenger Passenger) {
	trips := viewPassengerTrips(passenger.PassengerID)

	for t := len(trips) - 1; t >= 0; t-- { //Look for
		trip := trips[t]
		fmt.Println()
		fmt.Println("Trip ID: ", trip.TripID)
		fmt.Println("Trip Status: ", trip.TripStatus)
		fmt.Println("Driver ID: ", trip.DriverID)
		fmt.Println("Passenger ID: ", trip.PassengerID)
		fmt.Println("Pick Up Postal Code: ", trip.PickUpPoint)
		fmt.Println("Drop Off Postal Code: ", trip.DropOffPoint)
		fmt.Println()
	}
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
//						Main menu features for user accessing GrabnGo as Driver						  //																						////
////////////////////////////////////////////////////////////////////////////////////////////////////////
func driverLogin() {
	for {
		//Driver login using their mobile number
		fmt.Print("Please enter your mobile number: ")
		var mobileno string
		fmt.Scanln(&mobileno)

		//Look for registered driver with the given mobile number
		driver := getDriverMobileNo(mobileno)
		if mobileno != driver.MobileNo { //If mobile number does not exist, login fail.
			fmt.Println("\nInvalid mobile number. Login failed.")
			break
		} else {
			fmt.Printf("\nWelcome to GrabnGo Driver, %s %s!\n", driver.FirstName, driver.LastName)
			driverMainMenu(driver)
			break
		}
	}
}

func driverRegister() {
	//Prompt user to enter all the required information to register as a driver
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

////////////////////////////////////////////////////////////////////////////////////////////////////////
///				  						Driver's menu features										 ///
////////////////////////////////////////////////////////////////////////////////////////////////////////
func driverMainMenu(driver Driver) {
	for {
		fmt.Println("[1] Start trip")
		fmt.Println("[2] End trip")
		fmt.Println("[3] Update personal information")
		fmt.Println("[0] Logout")

		fmt.Print("\nEnter your option: ")
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "1" {
			driverStartTrip(driver)
		} else if userInput == "2" {
			driverEndTrip(driver)
		} else if userInput == "3" {
			driverUpdateInformation(driver)
		} else if userInput == "0" {
			break
		} else {
			fmt.Println("\nInvalid Option")
			driverMainMenu(driver)
		}
	}
}

//driverStartTrip and driverEndTrip is to set trip status from "Processing" to
//"Ongoing" when driver is assigned to trip, and
//"Ended" when driver finished his trip
func driverStartTrip(driver Driver) {
	driverTrips := viewDriverTrips(driver.DriverID)
	//fmt.Println(driverTrips)
	var initTripStatus Trip

	for _, t := range driverTrips { //Retrieve the values only that consists of "Processing" status of the trip for the specific driver
		if t.TripStatus == "Processing" {
			initTripStatus = t //Stored in the list to update the Trip Status
		}
	}
	if (initTripStatus == Trip{}) { //Check the list, if previously the list did not store any Trip Status consisting of "Processing", it will print the following message.
		fmt.Println("No processing trips found")
	} else {
		initTripStatus.TripStatus = "Ongoing" //update Trip Status from "Processing" to "Ongoing"
		UpdateTripInfo(initTripStatus)
		fmt.Println("Trip ongoing")
	}
}

func driverEndTrip(driver Driver) {
	driverTrips := viewDriverTrips(driver.DriverID)
	var initTripStatus Trip
	for _, t := range driverTrips { //Retrieve the values only that consists of "Ongoing" status of the trip for the specific driver
		if t.TripStatus == "Ongoing" {
			initTripStatus = t //Stored in the list to update the Trip Status
		}
	}
	if (initTripStatus == Trip{}) { //Check the list, if previously the list did not store any Trip Status consisting of "Ongoing", it will print the following message.
		fmt.Println("No ongoing trips found")
	} else {
		initTripStatus.TripStatus = "Ended" //update Trip Status from "Ongoing" to "Ended"
		UpdateTripInfo(initTripStatus)
		fmt.Println("Trip ended")
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

func viewPassengerTrips(PassengerID string) []Trip {
	var trips []Trip

	url := fmt.Sprintf("%s?PassengerID=%s", tripUrl, PassengerID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Passenger is not assigned to any trips")
		return trips
	}

	json.NewDecoder(resp.Body).Decode(&trips)
	return trips
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

func viewDriverTrips(DriverID string) []Trip {
	var trips []Trip

	url := fmt.Sprintf("%s?DriverID=%s", tripGetDriverIDUrl, DriverID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Driver is not assigned to any trips")
		return trips
	}

	json.NewDecoder(resp.Body).Decode(&trips)
	return trips
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

func UpdateTripInfo(newTripInfo Trip) error {
	fmt.Println()
	url := fmt.Sprintf("%s?TripID=%s", tripUrl, newTripInfo.TripID)

	_, err := httpPut(url, newTripInfo)
	return err
}
