package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := LoadConfig()

	http.HandleFunc("/send", HandleSend)
	http.HandleFunc("/verify/", HandleVerify)

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	log.Printf("Сервер запущен на http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
