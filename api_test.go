package main

import (
	"encoding/json"
	"irbis_api/core/irbis_hand"
	"log"
	"os"
	"testing"
)

type TestAcces struct {
	User string `json:"test_user"`
	Pass string `json:"test_pass"`
	Base string `json:"test_base"`
	UId  string `json:"test_userid"`
	UFam string `json:"test_family"`
}

var ts = configR("settings/test_config.json")

func configR(path string) *TestAcces {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("Не удалось прочитать конфигурационный файл: ", err)
	}
	var ta TestAcces
	err = json.Unmarshal(file, &ta)
	if err != nil {
		log.Println(err)
	}
	return &ta
}
func TestBooksOnHand(t *testing.T) {
	respond, err := irbis_hand.UserBooksOnHands(ts.User, ts.Pass, ts.UId, ts.UFam)
	if err != nil {
		t.Errorf("Expected error: nil, but got: %v", err)
	}
	check := json.Valid([]byte(respond))
	if !check {
		t.Errorf("Expected valid json, but... got this: %v", respond)
	}
}

func TestGenRecords(t *testing.T) {
	testCases := []struct {
		start int
		end   int
	}{
		{1, 20},
		{21, 40},
		{41, 60},
		{61, 80},
	}
	for _, e := range testCases {
		respond, err := irbis_hand.GenRecords("IKNBU", ts.User, ts.Pass, e.start, e.end)
		if err != nil {
			t.Errorf("Expected error: nil, but got: %v", err)
		}

		check := json.Valid([]byte(respond))
		if !check {
			t.Errorf("Expected valid json, but... got this: %v", respond)
		}

	}
}

func TestCollectRecords(t *testing.T) {
	tc := []struct {
		base  string
		start int
		end   int
	}{
		{"IKNBU", 1, 20},
		{"IKNBU", 21, 31},
		{"IKNBU", 32, 42},
		{"IKNBU", 43, 53},
	}
	for _, e := range tc {
		resp, err := irbis_hand.CollectRecords(e.base, ts.User, ts.Pass, e.start, e.end)
		if json.Valid([]byte(resp)) || err != nil {
			t.Errorf("Expected json array, but got: %s", resp)
		}

		log.Println(resp)
	}
}
