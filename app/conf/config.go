package conf

import (
	"os"
)

var Postgres = postgresDB{
	Url:    getEnv("DB_URL", "postgresql://postgres:postgres@localhost/postgres"),
	DBName: getEnv("DB_NAME", "myDB"),
}

var Tg = telegram{
	Token: getEnv("TELEGRAM_TOKEN", ""),
}

var AppConf = appConfig{
	Host:        getEnv("SERVER_HOST", "localhost"),
	Port:        getEnv("SERVER_PORT", "1337"),
	ContentPath: getEnv("CONTENT_PATH", "./"),
	ServerUrl:   getEnv("SERVER_URL", "localhost"),
}

type postgresDB struct {
	Url    string
	DBName string
}

type telegram struct {
	Token string
}

type appConfig struct {
	Host        string
	Port        string
	ContentPath string
	ServerUrl   string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
