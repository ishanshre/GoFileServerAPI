package helpers

import "os"

func SetFilePermission(realtiveFilepath string) error {
	return os.Chmod(realtiveFilepath, 0644)
}
