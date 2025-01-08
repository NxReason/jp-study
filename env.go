package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func Env() (map[string]string, error) {
	file, err := os.Open(".env")
	if err != nil { return nil, err }
	defer file.Close()

	result := make(map[string]string, 10)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { continue }

		split := strings.Split(line, "=")
		if len(split) < 2 {
			return nil, errors.New("Incorrect format at line: " + line)
		}
		
		result[split[0]] = split[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}