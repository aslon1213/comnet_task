package initializers

import "github.com/joho/godotenv"

// laod envs from .env and .env.example files
func load_envs() error {
	if err := godotenv.Load(".env", ".env.example"); err != nil {
		return err
	}
	return nil
}
