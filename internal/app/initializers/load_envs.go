package initializers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func load_envs() {

	fmt.Println("Loading environment variables...")
	if os.Getenv("GO_MODE") == "production" {
		godotenv.Load(".env.production")
		return
	} else {
		godotenv.Load(".env.development")
	}

}
