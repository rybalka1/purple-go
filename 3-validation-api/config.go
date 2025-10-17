package main

type Config struct {
	Email    string
	Password string
	Address  string
	Port     int
}

func LoadConfig() Config {
	return Config{
		Email:    "youremail@example.com",
		Password: "yourpassword",
		Address:  "localhost",
		Port:     8081,
	}
}
