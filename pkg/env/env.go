package env

import (
	"github.com/joho/godotenv"
)

// NewEnvConfig new instance of configuration
func NewEnvConfig(configPath string) map[string]string {
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		panic(err)
	}

	return myEnv
}
