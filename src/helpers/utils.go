package helpers

import (
	"log"
)

func CheckAndPrintErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}