package addtorrent

import (
	"strings"

	"github.com/atotto/clipboard"
)

func isMagnet(s string) bool {
	return strings.HasPrefix(s, "magnet:?xt=")
}

func getMagnetFromClipboard() string {
	placeholder, err := clipboard.ReadAll()
	// Clipboard empty
	if err != nil {
		placeholder = ""
	}
    if isMagnet(placeholder) {
		return placeholder
	}
	return ""
}


