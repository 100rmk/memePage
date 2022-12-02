package conf

import (
	"os"
)

var Mongo = mongoDB{
	Url:    getEnv("DB_URL", "mongodb://localhost:27017"),
	DBName: getEnv("DB_NAME", "myDB"),
}

var Tg = telegram{
	Token: getEnv("TELEGRAM_TOKEN", ""),
}

var AppConf = appConfig{
	Host: getEnv("SERVER_HOST", "localhost"),
	Port: getEnv("SERVER_PORT", "1337"),
}

type mongoDB struct {
	Url    string
	DBName string
}

type telegram struct {
	Token string
}

type appConfig struct {
	Host string
	Port string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
