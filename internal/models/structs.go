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

type ServStatusD struct {
	RegClients   int      `json:"reg_clients"`
	RunNow       int      `json:"run_now"`
	RunNowDetail []string `json:"run_now_detail"`
	TotalComm    int      `json:"total_com"`
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

type ConnParams struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type OnHands struct {
	Book               string `json:"brief"`
	DateOfIssue        string `json:"issue_date"`
	ExpectedReturnDate string `json:"return_date"`
}

type MapOnHands struct {
	Number int      `json:"num"`
	Value  []string `json:"book_info"`
}

type RecordDetails struct {
	Title          string `json:"title"`
	Author         string `json:"author"`
	AnotherAuthors string `json:"another_authors"`
	DocType        string `json:"doc_type"`
	Lang           string `json:"lang"`
	YearOfPubl     string `json:"year_of_publ"`
}
