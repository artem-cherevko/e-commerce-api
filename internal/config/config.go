package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT                  int64
	JWT_SECRET            string
	DB_DSN                string
	STRIPE_KEY            string
	STRIPE_WEBHOOK_SECRET string
}

func LoadENV() (*Env, error) {
	err := godotenv.Load(".env.dev")
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 0, 64)
	if err != nil {
		return nil, err
	}

	return &Env{
		PORT:                  port,
		JWT_SECRET:            os.Getenv("JWT_SECRET"),
		DB_DSN:                os.Getenv("DB_DSN"),
		STRIPE_KEY:            os.Getenv("STRIPE_KEY"),
		STRIPE_WEBHOOK_SECRET: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}, nil
}
