package main

import (
	"bufio"
	"os"
)

func readLines(f string) ([]string, error) {
	rd, err := os.Open(f)
        if err != nil {
		return nil, err
        }
	defer rd.Close()

	var lines []string
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		t := scanner.Text()
		lines = append(lines, t)
	}
	return lines, scanner.Err()
}
