package rakeja

import (
	"fmt"
	"os"
)

func debugf(format string, a ...any) {
	if os.Getenv("GORAKEJA_DEBUG") != "" {
		fmt.Printf(format, a...)
	}
}
