package irbis_hand

import (
	"encoding/json"
	"fmt"
	"irbis_api/internal/models"
	"log"
	"strings"
	"time"

	"github.com/amironov73/GoIrbis/src/irbis"
)

func UserBooksOnHands(login, password, user_id, last_name string) (string, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = "RDR"
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return "", fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()

	parameters := irbis.NewSearchParameters()
	parameters.Expression = fmt.Sprintf(`"K=%s$" * "K=%s$"`, user_id, last_name)
	//parameters.Format = irbis.OPTIMIZED_FORMAT
	//Используется свой формат вывода json, получается формируется он на стороне ИРБИС.
	parameters.Format = "@rdrw_raw_dolg"
	parameters.NumberOfRecords = 1
	found := conn.SearchEx(parameters)
	var first irbis.FoundLine
	if len(found) == 0 {
		println("Не нашли")
	} else {
		// в found находится слайс структур FoundLine
		first = found[0]
		//fmt.Println("MFN:", first.Mfn, "DESCRIPTION:", first.Description)
	}
	resp := strings.Split(first.Description, "\n")

	respond := models.Books{
		Books: resp,
	}
	jsonData, err := json.Marshal(respond)
	if err != nil {
		return "", fmt.Errorf("Ошибка упаковки json")
	}
	return string(jsonData), nil

}

func UserProfile(login, password, user_id, last_name string) (string, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = "RDR"
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return "", fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()

	found := conn.Search(fmt.Sprintf(`"K=%s$" * "K=%s$"`, user_id, last_name))
	if len(found) == 0 {
		return "", fmt.Errorf("Ничего не найдено")
	}

	var (
		userName     string
		userLastName string
		userSurname  string
		userCategory string
		userReg      string
	)

	for _, mfn := range found {
		record := conn.ReadRecord(mfn)
		userLastName = record.FM(10)
		userName = record.FM(11)
		userSurname = record.FM(12)
		userCategory = record.FSM(50, 'A')
		userReg = record.FM(51)
	}

	usermodel := models.UserInfo{
		UserName: fmt.Sprintf("%s.%s.%s", userName[0:2], userSurname[0:2], userLastName),
		Category: userCategory,
		RegDate:  fmt.Sprintf("%s.%s.%s", userReg[6:], userReg[4:6], userReg[0:4]),
	}

	jsonData, err := json.Marshal(usermodel)
	if err != nil {
		return "", fmt.Errorf("Ошибка упаковки json")
	}

	return string(jsonData), nil

}

func CreateVirUser(v *models.VirtualUserData) (string, error) {
	connection := irbis.NewConnection()
	connection.Host = "irbis"
	connection.Username = "amogus"
	connection.Password = "test14"
	if !connection.Connect() {
		println("Не удалось подключиться")
		return "", fmt.Errorf("Can't connect...")
	}

	defer connection.Disconnect()
	connection.Database = "RDRV"
	now := time.Now()
	fDate := now.Format("02.01.2006")
	fTime := now.Format("15:04")

	record := irbis.NewMarcRecord()
	record.Add(10, v.Family)
	record.Add(11, v.Name)
	record.Add(12, v.Surname)
	record.Add(21, v.Birth)
	record.Add(23, v.Gender)
	record.Add(17, v.Phone)
	record.Add(51, "").
		Add('6', fDate+","+fTime)
	record.Add(33, "Почтовый индекс: "+v.Postcode+", адрес:"+v.Country+","+v.City+", ул."+v.Street+", д."+v.House+",кв."+v.Apartment)

	record.Add(920, "RDR")
	connection.WriteRecord(record)

	log.Println(record)

	return fmt.Sprintf("%s", record), nil
}

func CoworkerProfile(clogin, cpassword string) (string, string, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = clogin
	conn.Password = cpassword
	//conn.Database = "RDR"
	if !conn.Connect() {
		println("Не удалось подключиться и получить данные работника")
		return "", "", fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()
	server_ver := conn.ServerVersion

	ini := conn.Ini
	dbaccess := ini.GetValue("Main", "DBNNAMECAT", "no_access")
	dbn := conn.ReadMenuFile(dbaccess)
	log.Println("Чтение файла: ")
	a := dbn.Entries
	for _, bases := range a {
		log.Println(bases)
	}
	log.Println(dbn)
	return server_ver, dbaccess, nil
}

func IrbStatus(login, password string) (int, int, int, error) {
	con := irbis.NewConnection()
	con.Host = "irbis"
	con.Port = 6666
	con.Workstation = "C"
	con.Username = login
	con.Password = password

	if !con.Connect() {
		log.Println("Не удалось подключиться и получить данные от сервера")
		return 0, 0, 0, fmt.Errorf("Не удалось подключиться и получить данные от сервера")
	}
	defer con.Disconnect()
	res := con.GetServerStat()
	connected := res.ClientCount
	comm := res.TotalCommandCount
	run := len(res.RunningClients)
	return connected, run, comm, nil
}

func Reload(login, password string) error {
	con := irbis.NewConnection()
	con.Host = "irbis"
	con.Port = 6666
	con.Workstation = "C"
	con.Username = login
	con.Password = password
	if !con.Connect() {
		log.Println("Не удалось перезапустить сервер")
		return fmt.Errorf("Не удалось подключиться к сереверу для перезапуска")
	}
	var response bool
	response = con.RestartServer()
	if response != true {
		return fmt.Errorf("Не удалось перезапустить сервер")
	}
	return nil
}
