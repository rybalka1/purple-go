package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Структуры для запросов и ответов
type PhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type SessionResponse struct {
	SessionID string `json:"sessionId"`
}

type CodeRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Code      int    `json:"code" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// Product struct
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// BuyRequest struct
type BuyRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// In-memory product list
var products = []Product{
	{ID: 1, Name: "Go Programming Book", Price: 49.99},
	{ID: 2, Name: "Sticker Pack", Price: 5.99},
	{ID: 3, Name: "Go Gopher Plush Toy", Price: 19.99},
}

// Временное хранилище для сессий и кодов (в реальном приложении использовалась бы база данных или Redis)
var sessions = make(map[string]struct {
	Phone string
	Code  int
})

// Секретный ключ для JWT (в реальном приложении должен быть в переменных окружения)
var jwtSecret = []byte("my_super_secret_key_for_jwt")

// GenerateJWT создает JWT-токен для пользователя
func GenerateJWT(phone string) (string, error) {
	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Токен действителен 72 часа
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware проверяет JWT-токен в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing or invalid"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("phone", claims["phone"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}

// SendCodeHandler обрабатывает запрос на отправку кода
func SendCodeHandler(c *gin.Context) {
	var req PhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Генерируем сессию
	sessionID := fmt.Sprintf("%x", rand.Int63()) // Простой способ генерации ID
	// 2. Генерируем код (4-значный)
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000 // от 1000 до 9999

	// 3. Сохраняем сессию и код (срок жизни 5 минут)
	sessions[sessionID] = struct {
		Phone string
		Code  int
	}{
		Phone: req.Phone,
		Code:  code,
	}

	// 4. Отправляем SMS (имитация)
	fmt.Printf("--- SMS SENT to %s: Code is %d (Session: %s) ---\n", req.Phone, code, sessionID)

	// 5. Возвращаем sessionId
	c.JSON(http.StatusOK, SessionResponse{SessionID: sessionID})
}

// VerifyCodeHandler обрабатывает запрос на проверку кода и выдачу токена
func VerifyCodeHandler(c *gin.Context) {
	var req CodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionData, ok := sessions[req.SessionID]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session ID"})
		return
	}

	if sessionData.Code != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
		return
	}

	// Код правильный, удаляем сессию
	delete(sessions, req.SessionID)

	// Генерируем JWT-токен
	token, err := GenerateJWT(sessionData.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: token})
}

// ProtectedHandler - защищенный ресурс
func ProtectedHandler(c *gin.Context) {
	phone, _ := c.Get("phone")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome, user with phone: %s", phone)})
}

// BuyProductHandler handles the logic for purchasing a product
func BuyProductHandler(c *gin.Context) {
	var req BuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone, _ := c.Get("phone")

	// Find the product
	var product Product
	found := false
	for _, p := range products {
		if p.ID == req.ProductID {
			product = p
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// "Process" the purchase
	totalPrice := product.Price * float64(req.Quantity)
	fmt.Printf("--- PURCHASE by %s: %d x %s for $%.2f ---\n", phone, req.Quantity, product.Name, totalPrice)

	c.JSON(http.StatusOK, gin.H{
		"message":     fmt.Sprintf("Successfully purchased %d of %s", req.Quantity, product.Name),
		"total_price": totalPrice,
	})
}

func main() {
	// Устанавливаем режим Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Группа для маршрутов авторизации
	auth := router.Group("/auth")
	{
		auth.POST("/send-code", SendCodeHandler)
		auth.POST("/verify-code", VerifyCodeHandler)
	}

	// Группа для защищенных маршрутов
	protected := router.Group("/api")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/protected", ProtectedHandler)
		protected.POST("/buy", BuyProductHandler)
	}

	fmt.Println("Server is running on :8080")
	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}
