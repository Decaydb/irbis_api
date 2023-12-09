package models

type VirtualUserData struct {
	Name      string
	Surname   string
	Family    string
	Birth     string
	Gender    string
	Phone     string
	Email     string
	Postcode  string
	Country   string
	City      string
	Street    string
	House     string
	Apartment string
}

type ErrorMessage struct {
	Error int
}

type ConnectionData struct {
	Version string
	Acces   string
}
