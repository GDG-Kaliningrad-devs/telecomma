package flag

import "os"

func Debug() bool {
	return os.Getenv("ENV") == "debug"
}
