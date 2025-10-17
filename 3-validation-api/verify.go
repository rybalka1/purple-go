package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func HandleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var data map[string]string
	json.Unmarshal(body, &data)

	emailAddr, ok := data["email"]
	if !ok {
		http.Error(w, "missing email", http.StatusBadRequest)
		return
	}

	hash := randomHash()
	saveVerification(Verification{Email: emailAddr, Hash: hash})

	if err := SendEmail(emailAddr, hash); err != nil {
		log.Println("Ошибка отправки письма:", err)
		http.Error(w, "send error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verification email sent"))
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/verify/"):]
	ok, _ := findAndRemove(hash)

	if ok {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func randomHash() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
