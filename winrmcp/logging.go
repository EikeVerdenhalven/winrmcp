package winrmcp

import (
	"log"
	"os"
)

func debugLog(message string) {
	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Print(message)
	}
}
