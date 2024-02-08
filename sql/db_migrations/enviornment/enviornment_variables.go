package enviornment

import (
	"os"
)

type EnvironmentVariables struct {
	PostgresAddress  string
	PostgresPort     string
	PostgresDB       string
	PostgresUsername string
	PostgresPassword string
}

func ProcessEnvironmentVariables() EnvironmentVariables {
	// In all cases the default behavior should be for the docker compose setup
	env := EnvironmentVariables{
		PostgresAddress:  "localhost",
		PostgresPort:     "5432",
		PostgresDB:       "postgres",
		PostgresUsername: "postgres",
		PostgresPassword: "testpassword",
	}

	envPostgresAddress := os.Getenv("POSTGRES_ADDRESS")
	envPostgresPort := os.Getenv("POSTGRES_PORT")
	envPostgresDB := os.Getenv("POSTGRES_DB")
	envPostgresUsername := os.Getenv("POSTGRES_USERNAME")
	envPostgresPassword := os.Getenv("POSTGRES_PASSWORD")

	if len(envPostgresAddress) != 0 {
		env.PostgresAddress = envPostgresAddress
	}

	if len(envPostgresPort) != 0 {
		env.PostgresPort = envPostgresPort
	}

	if len(envPostgresDB) != 0 {
		env.PostgresDB = envPostgresDB
	}

	if len(envPostgresUsername) != 0 {
		env.PostgresUsername = envPostgresUsername
	}

	if len(envPostgresPassword) != 0 {
		env.PostgresPassword = envPostgresPassword
	}

	return env
}
