package util

import (
	"fmt"
	"log"
)

func CheckError(err error) {
	if err != nil {
        log.Print(fmt.Sprintf("%w", err))
		panic(err)
	}
}
