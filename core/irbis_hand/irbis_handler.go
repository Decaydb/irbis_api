package irbis_hand

import (
	"encoding/json"
	"fmt"
	"irbis_api/internal/models"
	"log"
	"strconv"
	"strings"
	"sync"
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
	}

	resp := strings.Split(first.Description, "\n")
	if len(resp[0]) < 1 {
		resp = []string{}
	}
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

func IrbStatus(login, password, detail string) (int, int, int, []string, error) {
	con := irbis.NewConnection()
	con.Host = "irbis"
	con.Port = 6666
	con.Workstation = "C"
	con.Username = login
	con.Password = password

	if !con.Connect() {
		log.Println("Не удалось подключиться и получить данные от сервера")
		return 0, 0, 0, []string{}, fmt.Errorf("Не удалось подключиться и получить данные от сервера")
	}
	defer con.Disconnect()
	res := con.GetServerStat()
	connected := res.ClientCount
	comm := res.TotalCommandCount
	run := len(res.RunningClients)
	cowSlice := []string{}
	if detail != "" {
		for _, coworker := range res.RunningClients {
			cowSlice = append(cowSlice, coworker.Name)
		}
		return connected, run, comm, cowSlice, nil
	}
	return connected, run, comm, cowSlice, nil
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

func IrbBooksDetail(login, password, user_id, last_name string) (string, error) {

	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = "RDR"
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return "", fmt.Errorf("Ошибка подключения")
	}
	defer conn.Disconnect()
	var (
		brief      []string
		issueDate  []string
		returnDate []string
		wg         sync.WaitGroup
	)
	wg.Add(3)
	go func(co *irbis.Connection) {
		parameters := irbis.NewSearchParameters()
		parameters.Expression = fmt.Sprintf(`"K=%s$" * "K=%s$"`, user_id, last_name)
		parameters.Format = "@rdrw_raw_books" //Заменить на вывод только книг!!!!
		parameters.NumberOfRecords = 1
		found := co.SearchEx(parameters)
		var first irbis.FoundLine
		if len(found) == 0 {
			println("Не нашли")
		} else {
			first = found[0]
		}

		brief = strings.Split(first.Description, "\n")
		if len(brief[0]) < 1 {
			brief = []string{}
		}
		wg.Done()

	}(conn)

	go func(co *irbis.Connection) {
		parameters := irbis.NewSearchParameters()
		parameters.Expression = fmt.Sprintf(`"K=%s$" * "K=%s$"`, user_id, last_name)
		parameters.Format = "@rdrw_raw_date_issue" //Формировать вывод ТОЛЬКО дат выдачи!!!!!
		parameters.NumberOfRecords = 1
		found := co.SearchEx(parameters)
		var first irbis.FoundLine
		if len(found) == 0 {
			println("Не нашли")
		} else {
			first = found[0]
		}

		issueDate = strings.Split(first.Description, "\n")
		if len(issueDate[0]) < 1 {
			issueDate = []string{}
		}
		wg.Done()
	}(conn)

	go func(co *irbis.Connection) {
		parameters := irbis.NewSearchParameters()
		parameters.Expression = fmt.Sprintf(`"K=%s$" * "K=%s$"`, user_id, last_name)
		parameters.Format = "@rdrw_raw_date_return" //Сюда ТОЛЬКО даты предполагаемого возврата книг!
		parameters.NumberOfRecords = 1
		found := co.SearchEx(parameters)
		var first irbis.FoundLine
		if len(found) == 0 {
			println("Не нашли")
		} else {
			first = found[0]
		}

		returnDate = strings.Split(first.Description, "\n")
		if len(returnDate[0]) < 1 {
			issueDate = []string{}
		}
		wg.Done()
	}(conn)
	println(len(brief), len(issueDate), len(returnDate))
	wg.Wait()

	//ma := make(map[int]models.OnHands)

	/*
		dic:= func(books, dateOfIssue, dateOfReturn []string) map[int][3]string{
			m := make(map[int][3]string)
			for i:=0; i <len(books);i++{
				m[i] = [3]string{books[i],dateOfIssue[i],dateOfReturn[i]}
			}
			return m
		}
	*/
	dic := func() map[int]models.OnHands {
		m := make(map[int]models.OnHands)
		for i := 0; i < len(brief); i++ {
			m[i] = models.OnHands{
				Book:               brief[i],
				DateOfIssue:        issueDate[i],
				ExpectedReturnDate: returnDate[i],
			}
		}
		return m
	}()

	jsonData, err := json.Marshal(&dic)
	if err != nil {
		log.Println("Ошибка упаковки карты")
	}

	return string(jsonData), nil

}

func GuidSearchRecord(login, password, guid string) (string, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = "IKNBU"
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return "", fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()

	docT := map[string]string{
		"a":  "текстовые материалы, кроме рукописных; в том числе печатные текстовые материалы, микроформы, а также электронные текстовые материалы.",
		"a1": "Препринт",
		"a2": "ксерокопия",
		"b":  "текстовые материалы, рукописные; в том числе микроформы и оцифрованные рукописные текстовые материалы.",
		"c":  "музыкальные партитуры, кроме рукописных; в том числе печатные музыкальные партитуры, микроформы, а также электронные музыкальные партитуры.",
		"d":  "музыкальные партитуры, рукописные; в том числе микроформы и оцифрованные рукописные музыкальные партитуры.",
		"e":  "картографические материалы, кроме рукописных; в том числе географические карты, атласы, глобусы, цифровые географические карты, а также другие картографические материалы.",
		"f":  "картографические материалы, рукописные; в том числе микроформы и оцифрованные рукописные географические карты.",
		"g":  "кинофильмы, проекционные и видеоматериалы (без уточнения)",
		"g1": "видеозаписи",
		"g2": "кинофильмы",
		"g3": `проекционные материалы (диафильмы, слайды, пленки, НО для непроекционной двухмерной графики-см. "k" ниже).`,
		"i":  "звукозаписи, немузыкальные",
		"j":  "звукозаписи, музыкальные",
		"k":  `двухмерная графика (изоматериал, иллюстрации, чертежи и т. п.-графики, схемы, коллажи, компьют.графика, рисунки, "duplication masters", живописные изображения, фотонегативы, фотоотпечатки, почтовые открытки, плакаты, эстампы, "spirit masters", технические чертежи, фотомеханические репродукции, а также репродукции перечисленных выше материалов.`,
		"l":  `компьютерные файлы. Включает следующие классы электронных ресурсов: программное обеспечение (в том числе программы, игры, шрифты), числовые данные, мультимедиа, онлайновые системы или службы. Для этих классов материалов, если существует важный аспект, требующий отнесения материала к иной категории, определяемой значением кода поз.6 маркера, вместо кода l используется код, соответствующий данному аспекту (например, картографические векторные изображения кодируются как картографические материалы, а не как числовые данные). Другие классы электронных ресурсов кодируются в соответствии с наиболее важными аспектами ресурса (например, текстовый материал, графика, картографический материал, музыкальная или не-музыкальная звукозапись, движущееся изображение). В случаях, если наиболее важный аспект не может быть определен однозначно, документ кодируется как "компьютерный файл".`,
		"l1": "сайт",
		"l2": "электронный журнал",
		"m":  "мультимедиа (документ, содержащий компоненты двух или более видов; ни один из компонентов не является основным в наборе)",
		"m2": "комплект (документ, содержащий информацию на разных носителях, включая текстовый документ)",
		"r":  "трехмерные художественные объекты и реалии. Включает искусственные объекты, такие как: модели, диорамы, игры, головоломки, макеты, скульптуры и другие трехмерные художественные объекты и их репродукции, экспонаты, устройства, предметы одежды, игрушки, а также естественные объекты, например: препараты для микроскопа и другие предметы, смонтированные для визуального изучения.",
		"1":  "Шрифт Брайля И Муна",
		"2":  "Шрифт Брайля (РТШ - рельефно-точечный шрифт)",
		"3":  "Шрифт Муна",
		"4":  "УШ (укрупненный шрифт)",
		"5":  "Тактильное издание",
		"6":  "РГП (рельефно-графическое пособие)"}

	found := conn.Search(fmt.Sprintf(`"SKGUID=%s$"`, guid))
	if len(found) == 0 {
		return "", fmt.Errorf("Ничего не найдено")
	}
	var rec models.RecordDetails
	for _, mfn := range found {
		record := conn.ReadRecord(mfn)
		rec.Title = record.FSM(200, 'A') + record.FSM(200, 'E') + ", " + record.FSM(200, 'F')
		rec.Author = record.FSM(700, 'A') + ", " + record.FSM(700, 'G')
		rec.AnotherAuthors = record.FSM(701, 'A') + ", " + record.FSM(701, 'G')
		rec.DocType = docT[record.FSM(900, 'T')]
		rec.Lang = record.FM(101)
		rec.YearOfPubl = record.FSM(210, 'A') + ", " + record.FSM(210, 'D') + ", " + record.FSM(210, 'C')
	}

	jsonData, err := json.Marshal(&rec)
	if err != nil {
		return "", fmt.Errorf("Ошибка упаковки json")
	}

	return string(jsonData), nil
}

func SoloMfn(base, login, password string) (int, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Workstation = irbis.ADMINISTRATOR
	conn.Username = "avt"
	conn.Password = "hedge"
	conn.Database = base
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return 0, fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()
	based := conn.GetDatabaseInfo(base)
	return based.MaxMfn, nil
}

func GenRecords(base, login, password string, start, end int) (string, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = base
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return "", fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()

	resSlice := models.Records{}

	record := func(num int) models.Record {
		r := conn.ReadRecord(num)
		var booksCont = models.Record{
			GUID:   r.FM(1119),
			Author: r.FSM(700, 'A') + " " + r.FSM(700, 'B'),
			Title:  r.FSM(200, 'A'),
			Year:   r.FSM(210, 'D'),
		}
		return booksCont
	}
	if start == 0 {
		start = 1
	}
	for i := start; i <= end; i++ {
		resSlice.Books = append(resSlice.Books, record(i))
	}

	jsonData, err := json.Marshal(resSlice)
	if err != nil {
		return "", fmt.Errorf("Ошибка упаковки json")
	}

	return string(jsonData), nil

}

func FindBlock(base, login, password string) (string, error) {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = base
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return "", fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()

	db := conn.GetDatabaseInfo(base)
	dbBlocks := db.LockedRecords
	group := models.MfnBlocks{
		MFNs: dbBlocks,
	}
	jsonData, err := json.Marshal(group)
	if err != nil {
		return "", fmt.Errorf("Ошибка упаковки json")
	}
	return string(jsonData), nil
}

func UnlockMfns(mfn, base, login, password string) error {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = login
	conn.Password = password
	conn.Database = base
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()
	mfnList := []int{}
	if strings.Contains(mfn, ",") {
		strMfn := strings.Split(mfn, ",")
		for _, num := range strMfn {
			MfnInt, err := strconv.Atoi(num)
			if err != nil {
				return fmt.Errorf("В список MFN для разблокировки передано не число")
			}
			mfnList = append(mfnList, MfnInt)
		}
	} else {
		a, err := strconv.Atoi(mfn)
		if err != nil {
			return fmt.Errorf("В качестве номера MFN передано не число")
		}
		mfnList = append(mfnList, a)
	}
	conn.UnlockRecords(base, mfnList)

	return nil

}

func Gbl(gbl *models.GCor) error {
	conn := irbis.NewConnection()
	conn.Host = "irbis"
	conn.Port = 6666
	conn.Username = gbl.Login
	conn.Password = gbl.Password
	conn.Database = gbl.Base
	if !conn.Connect() {
		println("Не удалось подключиться для получения данных пользователя")
		return fmt.Errorf("{Error %v", "Не удалось подключиться к IRBIS")
	}
	defer conn.Disconnect()

}
