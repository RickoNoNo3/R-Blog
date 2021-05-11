// +build !windows

package errorhelper

import (
	"errors"
	"os"
)

func CheckErrorFileNotFound(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}
