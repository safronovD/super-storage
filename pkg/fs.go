package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type FSManager struct{}

func (fs *FSManager) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func (fs *FSManager) WriteFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, os.ModePerm)
}

func (fs *FSManager) ReadDir(dirPath string) ([][]byte, error) {
	var (
		resultErr  []string
		resultData [][]byte
	)

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			filePath := path.Join(dirPath, f.Name())
			data, err := fs.ReadFile(filePath)
			if err != nil {
				resultErr = append(resultErr, err.Error())
			} else {
				resultData = append(resultData, data)
			}
		}
	}

	if len(resultErr) != 0 {
		return nil, fmt.Errorf(strings.Join(resultErr, "\n"))
	}

	return resultData, nil
}
