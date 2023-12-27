package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func Write(errType string, errBody error) {
	file, err := os.OpenFile("logs/"+errType+"Errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Запись логов невозможна!")
	}
	defer file.Close()
	time := time.Now()
	file.WriteString(fmt.Sprintf("%v", errBody) + "\n")
	logger := log.New(file, fmt.Sprintf("ERROR|%s;%s", strings.ToUpper(errType), time.Format("2006-01-02 15:04:05")), log.LstdFlags)
	logger.Println(err)
}
