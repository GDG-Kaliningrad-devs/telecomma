package flag

import "os"

func Debug() bool {
	return os.Getenv("ENV") == "debug"
}

func AdminPass() string {
	pass := os.Getenv("ADMIN_PASS")
	if pass == "" {
		panic("setup password (ADMIN_PASS)")
	}

	return pass
}
