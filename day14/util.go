package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
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

func mustParseInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("ParseInt(%q) err=%v", s, err)
	}
	return v
}
