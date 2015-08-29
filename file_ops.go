package confection2

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

const (
	dirPerm  = 0755
	filePerm = 0633
)

func readFile(file string) ([]byte, error) {
	if ok := fileExist(file); ok {
		return ioutil.ReadFile(file)
	}

	return nil, errors.New("Config file does not exist")
}

func writeFile(file string, body []byte) error {
	var fd *os.File
	var err error
	if ok := fileExist(file); !ok {
		if err = mkdirp(file); err != nil {
			return err
		}
		if fd, err = os.Create(file); err != nil {
			return err
		}
	} else if fd, err = os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, filePerm); err != nil {
		return err
	}
	defer fd.Close()

	if _, err := fd.Write(body); err != nil {
		return err
	}

	return nil
}

func fileExist(file string) bool {
	_, err := os.Stat(file)
	return (err == nil)
}

func mkdirp(file string) error {
	dir := path.Dir(file)
	return os.MkdirAll(dir, dirPerm)
}