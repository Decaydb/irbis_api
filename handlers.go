package main

import (
	"encoding/json"
	"fmt"
	irb "irbis_api/core/irbis_hand"
	"irbis_api/core/validation"
	"irbis_api/internal/models"
	"log"
	"net/http"
	"strings"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	//v := r.URL.Values()
	//v.Add("last_name=")
	w.Header().Set("Contetn-Type", "application/json")
	var (
		errType     = []int{0, 0, 0, 0} // Первый элемент - ошибка авторизации. Второй - неверные данные пользователя
		description string
		errd        error
		errSend     int
	)
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	id := q.Get("id")
	lastName := q.Get("last_name")
	//var status string
	if len(login) < 1 || len(password) < 1 {
		//status = "Error"
		errType[0] = 1
	}
	if len(id) == 0 || len(lastName) < 1 {
		errType[1] = 1
	}

	if errType[0] == 0 && errType[1] == 0 {
		description, errd = irb.UserProfile(login, password, id, lastName)
		if errd != nil {
			errType[2] = 1
			log.Println(errd)
		}
	}

	switch {
	case errType[0] == 1:
		if errType[1] == 1 {
			errSend = 3
		} else {
			errSend = 1
		}
	case errType[0] == 0 && errType[1] == 1:
		errSend = 3
	case errType[0] == 0 && errType[1] == 0 && errType[2] == 1:
		errSend = 4
	}

	if errSend == 0 {
		// localhost:8080/api/v1/get.user?login=amogus&password=test1488&id=201240&last_name=Шамшурина
		//Получение данных о пользователе.
		//description, errd := irbis_hand.UserProfile(login, password, id, lastName)

		fmt.Fprint(w, description)

	} else {
		group := models.ErrorMessage{
			Error: errSend,
		}
		b, err := json.Marshal(group)
		if err != nil {
			log.Println("error:", err)
		}
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

func CreateVirtual(w http.ResponseWriter, r *http.Request) {
	//
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

func ServerStatus(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	login := q.Get("login")
	password := q.Get("password")
	if len(login) > 1 && len(password) > 1 {

		conn, run, comm, err := irb.IrbStatus(login, password)
		if err == nil {
			group := models.ServStatus{
				RegClients: conn,
				RunNow:     run,
				TotalComm:  comm,
			}
			b, _ := json.Marshal(group)
			fmt.Fprint(w, string(b))

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
