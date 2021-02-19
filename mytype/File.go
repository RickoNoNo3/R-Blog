package mytype

import (
	"io/ioutil"
	"os"

	"rickonono3/r-blog/helper/typehelper"
)

type File struct {
	Entity Entity
}

func (file *File) getFilePath() string {
	return "file_" + typehelper.MustItoa(file.Entity.Id)
}

func (file *File) getData() (data []byte, err error) {
	var osFile *os.File
	if osFile, err = os.Open(file.getFilePath()); err == nil {
		return ioutil.ReadAll(osFile)
	}
	return
}
