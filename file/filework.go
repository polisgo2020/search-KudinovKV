package file

import (
	"io/ioutil"
	"os"
)

// ReadFile read filename and return string data
func ReadFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile write data in filename
func WriteFile(data string, filename string) error {
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	_, err = outfile.WriteString(data)
	if err := outfile.Sync(); err != nil {
		return err
	}

	return nil
}
