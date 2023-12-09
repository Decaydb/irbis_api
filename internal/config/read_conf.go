package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Dbase struct {
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbName string `json:"db_name"`
}

type Vk struct {
	Vkey    string `json:"vkey"`
	GroupId int    `json:"group_id"`
}

func RdbConfig() string {

	file, err := os.ReadFile("/config/db.json")
	if err != nil {
		log.Println(err)
	}
	var dbcon Dbase
	err = json.Unmarshal(file, &dbcon)
	if err != nil {
		log.Fatal(err)
	}
	// Строка подключения к базе данных
	connection := fmt.Sprintf("user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbcon.DbUser, dbcon.DbPass, dbcon.DbName)

	return connection
}

func RmainConf() (string, int) {

	file, err := os.ReadFile("config/vk.json")
	if err != nil {
		log.Println(err)
	}
	var vkcon Vk
	err = json.Unmarshal(file, &vkcon)
	if err != nil {
		log.Fatal(err)
	}
	return vkcon.Vkey, vkcon.GroupId
}
