package common

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func ErrorHandler(err error, fatal bool) {
	if err != nil {
		log.Errorf("error: %v", err)

		if fatal {
			panic(fmt.Errorf("error: %s", err))
		}
	}
}

func OpenFile(filename string) *os.File {
	file, fileErr := os.Open(filename)
	if fileErr != nil {
		err := fmt.Errorf("failed to open file %v; [error_msg: %v]", filename, fileErr)
		ErrorHandler(err, true)
	}
	return file
}
