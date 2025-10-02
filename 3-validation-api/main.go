package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/jordan-wright/email"
)

// простая память для хранения токенов подтверждения
var (
	verificationTokens = make(map[string]string) // hash -> email
	mu                 sync.Mutex
)

// отправка письма
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		http.Error(w, "missing 'to' parameter", http.StatusBadRequest)
		return
	}

	// генерируем простой hash (можно заменить на crypto/rand)
	hash := fmt.Sprintf("%d", time.Now().UnixNano())

	// сохраняем в память
	mu.Lock()
	verificationTokens[hash] = to
	mu.Unlock()

	// готовим письмо
	e := email.NewEmail()
	e.From = os.Getenv("SMTP_USER")
	e.To = []string{to}
	e.Subject = "Email confirmation"
	e.Text = []byte(fmt.Sprintf("Для подтверждения перейдите по ссылке: http://localhost:8080/verify/%s", hash))

	// подключение к SMTP
	smtpAddr := os.Getenv("SMTP_ADDRESS") // например: smtp.gmail.com:587
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	err := e.Send(smtpAddr, smtp.PlainAuth("", smtpUser, smtpPass, smtpHost(smtpAddr)))
	if err != nil {
		http.Error(w, "send error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("confirmation email sent"))
}

// проверка подтверждения
func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hash := r.URL.Path[len("/verify/"):]
	if hash == "" {
		http.Error(w, "missing hash", http.StatusBadRequest)
		return
	}

	mu.Lock()
	emailAddr, ok := verificationTokens[hash]
	mu.Unlock()

	if !ok {
		http.Error(w, "invalid or expired hash", http.StatusNotFound)
		return
	}

	// подтверждено, можно удалить из памяти
	mu.Lock()
	delete(verificationTokens, hash)
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Email %s подтверждён!", emailAddr)))
}

// утилита для извлечения хоста из "smtp.gmail.com:587"
func smtpHost(addr string) string {
	for i := 0; i < len(addr); i++ {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}

func main() {
	// проверка конфигурации
	if os.Getenv("SMTP_USER") == "" || os.Getenv("SMTP_PASS") == "" || os.Getenv("SMTP_ADDRESS") == "" {
		log.Fatal("Нужно задать переменные окружения: SMTP_USER, SMTP_PASS, SMTP_ADDRESS")
	}

	http.HandleFunc("/send", sendEmailHandler)
	http.HandleFunc("/verify/", verifyEmailHandler)

	addr := ":8080"
	log.Printf("Server started at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
