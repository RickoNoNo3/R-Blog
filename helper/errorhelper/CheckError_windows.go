// +build windows

package errorhelper

import (
	"errors"
	"os"
	"syscall"
)

func CheckErrorFileNotFound(err error) bool {
	return errors.Is(err, os.ErrNotExist) || err == syscall.ERROR_FILE_NOT_FOUND || err == syscall.ERROR_PATH_NOT_FOUND
}
