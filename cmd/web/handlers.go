package web

import (
	"encoding/json"
	"fmt"
	"irbis_api/core/caching"
	irb "irbis_api/core/irbis_hand"
	"irbis_api/core/validation"
	"irbis_api/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var cache *caching.Cache

func GetUser(w http.ResponseWriter, r *http.Request) {
	//v := r.URL.Values()
	//v.Add("last_name=")
	w.Header().Set("Contetn-Type", "application/json")
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	id := q.Get("id")
	lastName := q.Get("last_name")

	if len(login) > 2 && len(password) > 2 {
		if len(id) > 1 && len(lastName) >= 2 {
			description, errd := irb.UserProfile(login, password, id, lastName)
			if errd == nil {
				fmt.Fprint(w, description)
			} else {
				group := models.ErrorMessage{
					Error: 1,
				}
				b, _ := json.Marshal(group)
				fmt.Fprint(w, string(b))
			}
		} else {
			group := models.ErrorMessage{
				Error: 2,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}
	} else {
		group := models.ErrorMessage{
			Error: 3,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	}

}

func WorkerLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetn-Type", "application/json")
	q := r.URL.Query()
	clogin := q.Get("login")
	cpassword := q.Get("password")
	ver, dbacces, err := irb.CoworkerProfile(clogin, cpassword)
	if err != nil {
		group := models.ErrorMessage{
			Error: 2,
		}
		b, err := json.Marshal(group)
		if err != nil {
			log.Println("error:", err)
		}
		fmt.Fprint(w, b)
	} else {
		info := models.ConnectionData{
			Version: ver,
			Acces:   dbacces,
		}
		b, _ := json.Marshal(info)
		fmt.Fprint(w, string(b))
	}
}

// TODO: Переработать регистрацию. Должна быть проверка существует ли такой юзер
// А json запросом + передавать не параметрами
func CreateVirtual(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetn-Type", "application/json")
	q := r.URL.Query()
	virtual := models.VirtualUserData{
		Name:      q.Get("name"),
		Surname:   q.Get("surname"),
		Family:    q.Get("family_name"),
		Birth:     q.Get("birth_date"),
		Gender:    q.Get("gender"),
		Phone:     q.Get("phone"),
		Email:     q.Get("email"),
		Postcode:  q.Get("postcode"),
		Country:   q.Get("country"),
		City:      q.Get("city"),
		Street:    q.Get("street"),
		House:     q.Get("house"),
		Apartment: q.Get("apartment"),
	}
	if err := validation.NamesValidation(virtual.Name, virtual.Family, virtual.Surname); err == nil {
		virtual.Phone, err = validation.PhoneValidation(virtual.Phone)
		if err == nil {
			e := strings.NewReplacer(" ", ".", ",", "\n", "\r", "\t")
			virtual.Name = e.Replace(virtual.Name)
			virtual.Surname = e.Replace(virtual.Surname)
			virtual.Family = e.Replace(virtual.Family)
			_, err := irb.CreateVirUser(&virtual)
			if err != nil {
				log.Println("Ошибка регистрации пользователя", virtual.Name)
			}

		} else {
			group := models.ErrorMessage{
				Error: 15,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}

	} else {
		group := models.ErrorMessage{
			Error: 14,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	}

}

// TODO: Убрать краткий вариант, только полная информация
func ServerStatus(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	detail := q.Get("detail")
	if len(login) > 1 && len(password) > 1 {

		conn, run, comm, runDetail, err := irb.IrbStatus(login, password, detail)
		if err == nil {
			if detail == "detail" {
				group := models.ServStatusD{
					RegClients:   conn,
					RunNow:       run,
					RunNowDetail: runDetail,
					TotalComm:    comm,
				}
				b, _ := json.Marshal(group)
				fmt.Fprint(w, string(b))

			} else {
				group := models.ServStatus{
					RegClients: conn,
					RunNow:     run,
					TotalComm:  comm,
				}
				b, _ := json.Marshal(group)
				fmt.Fprint(w, string(b))
			}
		} else {
			group := models.ErrorMessage{
				Error: 111,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}
	} else {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	}
}

func ReloadIrbis(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	if len(login) <= 1 && len(password) <= 1 {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	} else {
		err := irb.Reload(login, password)
		if err == nil {
			group := models.ErrorMessage{
				Error: 0,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		} else {
			group := models.ErrorMessage{
				Error: 100,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}

	}

}

func OnHands(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	id := q.Get("id")
	lastName := q.Get("last_name")
	if len(password) > 1 || len(login) > 1 {
		if len(id) > 4 && len(lastName) > 2 {
			resp, err := irb.UserBooksOnHands(login, password, id, lastName)
			if err != nil {
				group := models.ErrorMessage{
					Error: 4,
				}
				b, _ := json.Marshal(group)
				fmt.Fprint(w, string(b))
			} else {
				fmt.Fprint(w, resp)
			}

		} else {
			group := models.ErrorMessage{
				Error: 1,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}
	}
}

func OnHandsDetail(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	id := q.Get("id")
	lastName := q.Get("last_name")
	if len(password) > 1 || len(login) > 1 {
		if len(id) > 4 && len(lastName) > 2 {
			resp, err := irb.IrbBooksDetail(login, password, id, lastName)
			if err != nil {
				group := models.ErrorMessage{
					Error: 3,
				}
				b, _ := json.Marshal(group)
				fmt.Fprint(w, string(b))
			} else {
				fmt.Fprint(w, resp)
			}

		} else {
			group := models.ErrorMessage{
				Error: 2,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}
	} else {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	}
}

func GuidSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	guid := q.Get("guid")

	if len(login) < 1 && len(password) < 1 {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	} else {
		count := strings.Count(guid, "-")
		if count == 4 && len(guid) > 15 {
			result, err := irb.GuidSearchRecord(login, password, guid)
			if err != nil {
				println(err)
				group := models.ErrorMessage{
					Error: 3,
				}
				b, _ := json.Marshal(group)
				fmt.Fprint(w, string(b))
			} else {
				fmt.Fprint(w, result)
			}

		} else {
			group := models.ErrorMessage{
				Error: 2,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}

	}

}

func FormRecords(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	base := q.Get("base")
	pageStr := q.Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || len(base) == 0 {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	} else {
		totalRec, err := irb.SoloMfn(base, login, password)
		log.Println(totalRec)
		if err != nil {
			log.Println("Ошибка получения максимального количества записей в базе.")
		}
		recPerPage := 30
		totalPages := (totalRec + recPerPage - 1) / recPerPage
		if page < 1 {
			page = 1
		} else if page > totalPages {
			page = totalPages
		}

		start := (page - 1) * recPerPage
		end := start + recPerPage
		if end > totalRec {
			end = totalRec
		}
		start++
		end--

		records, err := irb.CollectRecords(base, login, password, start, end)
		if err == nil {
			fmt.Fprint(w, records)

		} else {
			group := models.ErrorMessage{
				Error: 2,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}

	}
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	base := q.Get("base")
	pageStr := q.Get("page")
	//Поиск в кэше по странице
	if elem, bol := cache.Get(base + "_" + pageStr); !bol {
		//Если пользователь не найден
	} else {
		switch v := elem.(type) {
		case string:
			fmt.Fprint(w, v)

		default:

		}

	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || len(base) == 0 {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	} else {
		totalRec, err := irb.SoloMfn("IKNBU", login, password)
		log.Println(totalRec)
		if err != nil {
			log.Println("Ошибка получения максимального количества записей в базе.")
		}
		recPerPage := 30
		totalPages := (totalRec + recPerPage - 1) / recPerPage
		if page < 1 {
			page = 1
		} else if page > totalPages {
			page = totalPages
		}

		start := (page - 1) * recPerPage
		end := start + recPerPage
		if end > totalRec {
			end = totalRec
		}

		records, err := irb.GenRecords(base, login, password, start, end)
		if err == nil {
			fmt.Fprint(w, records)

		} else {
			group := models.ErrorMessage{
				Error: 2,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}

	}
}

func MfnBlocks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	base := q.Get("base")
	if len(login) > 2 && len(password) > 2 && len(base) > 2 {
		res, err := irb.FindBlock(base, login, password)
		if err != nil {
			group := models.ErrorMessage{
				Error: 4,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		} else {
			fmt.Fprint(w, res)
		}

	} else {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	}
}

func UnblockRecs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	base := q.Get("base")
	mfn := q.Get("mfn")
	if len(login) > 2 && len(password) > 2 && len(mfn) != 0 {
		err := irb.UnlockMfns(base, mfn, login, password)
		if err != nil {
			group := models.ErrorMessage{
				Error: 4,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		} else {
			group := models.ErrorMessage{
				Error: 0,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))
		}
	} else {
		group := models.ErrorMessage{
			Error: 1,
		}
		b, _ := json.Marshal(group)
		fmt.Fprint(w, string(b))
	}
}

// TODO: Реализовать глобальную корректировку
func GlobalCorrect(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	gbl := models.GCor{
		Login:    q.Get("login"),
		Password: q.Get("password"),
		Base:     q.Get("base"),
		Records:  q.Get("records"),
		Gtype:    q.Get("type"),
		Field:    q.Get("field"),
		Value:    q.Get("value"),
		Actual:   q.Get("actual"),
	}
	if len(gbl.Value) == 0 {
		gbl.Value = "empty"
	}
	if len(gbl.Login) > 1 && len(gbl.Password) > 1 {
		if len(gbl.Gtype) > 2 && len(gbl.Field) != 0 && len(gbl.Value) != 0 {
			err := irb.Gbl(&gbl)
			if err != nil {
				//Обработка ошибки ирбиса
			}
		} else {
			//Ошибка - передаваемые значения некорректны
		}
	} else {
		//Ошибка - пустые логин и пароль
	}
}
