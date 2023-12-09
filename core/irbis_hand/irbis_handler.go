package irbis_hand

import (
	"fmt"
	"irbis_api/internal/models"
	"log"
	"time"

	"github.com/amironov73/GoIrbis/src/irbis"
)

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

	parameters := irbis.NewSearchParameters()
	parameters.Expression = fmt.Sprintf(`"K=%s$" * "K=%s$"`, user_id, last_name)
	//parameters.Format = irbis.OPTIMIZED_FORMAT
	//Используется свой формат вывода json, получается формируется он на стороне ИРБИС.
	parameters.Format = "@json_dolg2"
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
	return first.Description, nil

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
