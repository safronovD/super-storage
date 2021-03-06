package pkg

import (
	"fmt"
	"os"
	"path"
)

type FSInteraction interface {
	Write(file *DataFile, name string) error
	Read(name string, blockSizeInt int, positionArr []int) ([]byte, error)
}

var (
	NewWriter = New
)

type FSInteractionImpl struct {
	FileDir     string
	FilePattern string
}

func New(path string) (FSInteraction, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &FSInteractionImpl{
		FileDir:     path,
		FilePattern: "data",
	}, nil
}

func (w *FSInteractionImpl) Write(file *DataFile, name string) error {
	var (
		filePath = path.Join(w.FileDir, w.FilePattern+"-"+name)
	)

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return fmt.Errorf("file %s exists", filePath)
	}

	err := os.WriteFile(filePath, file.toByte(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (w *FSInteractionImpl) Read(name string, blockSizeInt int, positionArr []int) (result []byte, err error) {
	var (
		filePath  = path.Join(w.FileDir, w.FilePattern+"-"+name)
		blocks    []Block
		blockSize = blockSizeInt
	)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data); i += blockSize {
		if i+blockSize > len(data) {
			return nil, fmt.Errorf("data is corrupted")
		}

		block := Block{data[i : i+blockSize]}
		blocks = append(blocks, block)
	}

	fmt.Println(len(blocks))

	for _, pos := range positionArr {
		result = append(result, blocks[pos].Data...)
	}

	return
}
