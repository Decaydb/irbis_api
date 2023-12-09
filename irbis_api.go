package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/get.user", GetUser)
	http.HandleFunc("/api/v1/worker.login", WorkerLogin)
	http.HandleFunc("/api/v2/virtual.user", CreateVirtual)
	//mux := http.NewServeMux()
	//mux.HandleFunc("/api/v1/get.user", GetUser)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
