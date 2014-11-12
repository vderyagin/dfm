package host

import (
	"log"
	"os"
	"regexp"
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

// DotFileLocalSuffix returns a suffix to be added to host-specific dotfiles' paths.
func DotFileLocalSuffix() string {
	return ".host-" + Name()
}

// RemoveSuffix returns string without host-specific suffix, returns original
// string if it has no such suffix.
func RemoveSuffix(input string) string {
	return PathRegexp.ReplaceAllLiteralString(input, "")
}

// PathRegexp returns regexp matching path of any host-specific dotfile.
var PathRegexp = regexp.MustCompile(`\.host-[[:alpha:]]+\z`)
