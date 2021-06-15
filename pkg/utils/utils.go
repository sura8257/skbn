package utils

import (
	"log"
	"os"
	"strings"
)

// SplitInTwo splits a string to two parts by a delimeter
func SplitInTwo(s, sep string) (string, string) {
	if !strings.Contains(s, sep) {
		log.Fatal(s, "does not contain", sep)
	}
	split := strings.Split(s, sep)
	return split[0], split[1]
}

// Check file stat
func CheckFileStat(file string) error {
	stat, err := os.Stat(file)
	if err != nil {
		log.Fatal(file, "failed to check its stat")
	}
	log.Printf("File: %s  Size: %d bytes", stat.Name(), stat.Size())

	return nil
}
