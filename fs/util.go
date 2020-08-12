package fs

import "os"

// DirExists checks if a directory exists
func DirExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}
