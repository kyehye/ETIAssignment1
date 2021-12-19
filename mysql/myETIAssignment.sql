CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost';

CREATE DATABASE ridesharing_db;

USE ridesharing_db;

---'string' format is used across the Tables created in the database. Thus, all the microservices is querying based on "string" for values/index.
CREATE TABLE Passengers (PassengerID VARCHAR(5) NOT NULL PRIMARY KEY, FirstName VARCHAR(30), LastName VARCHAR(30), MobileNo VARCHAR(8), EmailAdd VARCHAR(100)); 
CREATE TABLE Drivers (DriverID VARCHAR(5) NOT NULL PRIMARY KEY, FirstName VARCHAR(30), LastName VARCHAR(30), MobileNo VARCHAR(8), EmailAdd VARCHAR(100), CarLicenseNo VARCHAR(10), Availability BOOL);  
CREATE TABLE Trips (TripID VARCHAR(5) NOT NULL PRIMARY KEY, TripStatus VARCHAR(20), PassengerID VARCHAR(5), DriverID VARCHAR(5), PickUpPoint VARCHAR(10), DropOffPoint VARCHAR(10)); 

--Below are sql dumps that can be used ONLY FOR TESTING PURPOSES using the GET, PUT method.
--All the required testing should be done via Postman.
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNo, EmailAdd) VALUES ("0001", "Jake", "Lee", "81234567", "jakelee@gmail.com");
INSERT INTO Passengers (PassengerID, FirstName, LastName, MobileNo, EmailAdd) VALUES ("0002", "Celine", "Tay", "87654321", "celinetay@gmail.com");

INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNo, EmailAdd, CarLicenseNo) VALUES ("0001", "Celest", "Teo", "91234567", "celteo@gmail.com","S1234567A");
INSERT INTO Drivers (DriverID, FirstName, LastName, MobileNo, EmailAdd, CarLicenseNo) VALUES ("0002", "Ang Kor", "Ngo", "97654321", "ngoak@gmail.com","S7654321B");

INSERT INTO Trips (TripID, TripStatus, PassengerID, DriverID, PickUpPoint, DropOffPoint) VALUES ("0001", "Ongoing", "0001", "0002", "650406","670302");
INSERT INTO Trips (TripID, TripStatus, PassengerID, DriverID, PickUpPoint, DropOffPoint) VALUES ("0002", "Ended", "0001", "0004", "650433","750002");
INSERT INTO Trips (TripID, TripStatus, PassengerID, DriverID, PickUpPoint, DropOffPoint) VALUES ("0003", "Ongoing", "0002", "0003", "698402","984722");