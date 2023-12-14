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

type ServStatus struct {
	RegClients int
	RunNow     int
	TotalComm  int
}

type User struct {
	UserName string `json:"user_name"`
	Category string `json:"category"`
	RegDate  string `json:"reg_date"`
	Books    []Book `json:"books"`
}

type Book struct {
	Brief              string `json:"brief"`
	DateOfIssue        string `json:"date_of_issue"`
	ExpectedReturnDate string `json:"expected_return_date"`
	PlaceOfIssue       string `json:"place_of_issue"`
}

type UserInfo struct {
	UserName string `json:"user_name"`
	Category string `json:"category"`
	RegDate  string `json:"reg_date"`
}

type Books struct {
	Books []string `json:"books"`
}
