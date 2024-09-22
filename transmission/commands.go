package transmission

import (
	"fmt"
	"os"
)

var OUTPUT_EMPTY_FILE = "sketch/transmission-list-1.txt"
var OUTPUT_REMOTE_FILE = "sketch/transmission-list-remote.txt"

func TransmissionList() ([]byte, error) {
	output, err := os.ReadFile(OUTPUT_REMOTE_FILE)
	if err != nil {
		return []byte(""), fmt.Errorf("can't read file: %w")
	}

	return output, nil
}
