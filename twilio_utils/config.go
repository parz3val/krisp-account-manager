package twilio_utils

import (
	"os"
	"github.com/joho/godotenv"
)

type Settings struct {
	AccountSid string
	AuthToken  string
}

func GetSettings() (settings Settings) {
	// read .env file
	err := godotenv.Load(".env")
	if err != nil {
		return Settings{}
	}
	return Settings{
		AccountSid: os.Getenv("ACCOUNT_SID"),
		AuthToken:  os.Getenv("AUTH_TOKEN"),
	}
}
