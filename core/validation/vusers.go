package validation

import (
	"fmt"
	"regexp"
	"strings"
)

func NamesValidation(firstName, lastName, middleName string) (err error) {
	pattern := `^[а-яА-Яa-zA-z]+$`
	//f, _ := regexp.Match(`^[а-яА-Яa-zA-z]+$`)
	f, _ := regexp.Match(pattern, []byte(firstName))
	l, _ := regexp.Match(pattern, []byte(lastName))
	m, _ := regexp.Match(pattern, []byte(middleName))
	if f != true && l != true && m != true {
		err = fmt.Errorf("В ФИО не только буквы")
	}

	return
}

func PhoneValidation(phone string) (string, error) {
	pattern := `^[\+\-]?[0-9s\-]+$`
	valid, _ := regexp.Match(pattern, []byte(phone))
	if valid != true {
		return "", fmt.Errorf("Not a phone number")
	}
	p := strings.NewReplacer(" ", "", ".", "", "\n", "", "\r", "")
	phone = p.Replace(phone)

	return phone, nil
}

func TextEditor(target string) (editResult string) {
	e := strings.NewReplacer(" ", "", ".", "", ",", "", "\n", "")
	editResult = e.Replace(target)
	return
}
