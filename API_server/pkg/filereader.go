package pkg

import (
	"os"
	"strings"
)

func ReadFromFile(name string) ([]string, error) {
	buf, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(buf), "\n"), nil
}
