package cfg

import (
	"log"
	"os"
	"strconv"
)

type Cfg struct {
	DBName     string
	DBPort     int
	DBHost     string
	DBUsername string
	DBPassword string
	DBParams   string
	JWTSecret  string
	BCryptSalt int
}

func Load() *Cfg {
	var err error
	cfg := &Cfg{}

	cfg.DBName = os.Getenv("DB_NAME")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBUsername = os.Getenv("DB_USERNAME")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBParams = os.Getenv("DB_PARAMS")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	cfg.BCryptSalt, err = strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		log.Fatal("fail convert bcrypt salt to int:", err)
	}
	cfg.DBPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("fail convert db port to int:", err)
	}

	return cfg
}
