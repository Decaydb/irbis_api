package web

import (
	"log"
	"net/http"
	"os"
)

func Routing() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/get.user", GetUser)
	mux.HandleFunc("/api/v1/worker.login", WorkerLogin)
	mux.HandleFunc("/api/v1/virtual.user", CreateVirtual)
	mux.HandleFunc("/api/v1/server.status", ServerStatus)
	mux.HandleFunc("/api/v1/reload.irbis", ReloadIrbis)
	mux.HandleFunc("/api/v1/on.hands", OnHands)
	mux.HandleFunc("/api/v1/on.hands.detail", OnHandsDetail)
	mux.HandleFunc("/api/v1/guid.search", GuidSearch)
	mux.HandleFunc("/api/v1/records", GetRecords)
	mux.HandleFunc("/api/v1/blocks", MfnBlocks)
	mux.HandleFunc("/api/v1/records.unlock", UnblockRecs)
	mux.HandleFunc("/api/v1/gbl", GlobalCorrect)

	f, err := os.OpenFile("logs/web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	infoLog.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	serr := http.ListenAndServe(":8080", mux)
	errorLog.Fatal(serr)
}
