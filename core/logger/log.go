package logger

import (
	"log"
	"os"
)

func Info(content ...any) {

	f, err := os.OpenFile("logs/web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	infoLog.Println(content...)
}

func Error(content ...any) {
	f, err := os.OpenFile("logs/web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog.Println(content...)
}

func Fatal(content ...any) {
	f, err := os.OpenFile("logs/web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	errorLog := log.New(f, "FATAL ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog.Fatal(content...)
}
