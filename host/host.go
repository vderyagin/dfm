package host

import (
	"log"
	"os"
)

// Name returns a hostname of a host machine. Can be overridden by setting
// HOST enviroment variable.
func Name() string {
	if envHost := os.Getenv("HOST"); len(envHost) != 0 {
		return envHost
	}

	name, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to retrieve hostname.")
	}

	return name
}
