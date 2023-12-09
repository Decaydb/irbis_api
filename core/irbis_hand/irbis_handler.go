package irbis_hand

import (
	"fmt"
	"irbis_api/internal/models"
	"log"

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
	log.Println("Имя:", v.Name, "Фамилия:", v.Family, "Отчество:", v.Surname, "Дата рождения:", v.Birth, "Телефон:", v.Phone)
	return "", nil
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
