package vse

import (
	"log"
	"os"
	"strings"
	"unicode"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// TODO: Optimize and use pointers
func stringMinifier(in string) (out string) {
	white := false
	for _, c := range strings.TrimSpace(in) {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return out
}
