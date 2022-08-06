package fs

import (
	"log"
	"os"
)

func AssertNotExists(path string) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Fatal("File or directory already exists")
	}
}

func AppendString(path, s string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(s); err != nil {
		log.Fatal(err)
	}
}
