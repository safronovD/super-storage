package pkg

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

type Writer interface {
	Write(file []byte, name string) error
}

var (
	NewWriter = New
)

type WriterImpl struct {
	FileDir     string
	FilePattern string
	BlockSize   uint8
}

func New(path string) (Writer, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &WriterImpl{
		FileDir:     path,
		FilePattern: "data",
		BlockSize:   4,
	}, nil
}

type Block struct {
	Hash [20]byte
	Data []byte
	Id   uint32
}

type Blocks []Block

func (b *Block) toByte() (bytes []byte) {
	bytes = make([]byte, 0, len(b.Hash)+len(b.Data)+4)
	bytes = append(bytes, b.Hash[:20]...)
	bytes = append(bytes, b.Data...)
	bytes = append(bytes, []byte(strconv.Itoa(int(b.Id)))...)
	return
}

func (bs *Blocks) toByte() (bytes []byte) {
	bytes = make([]byte, 0, len(*bs)*24)
	for _, block := range *bs {
		bytes = append(bytes, block.toByte()...)
	}
	return
}

var (
	blocks      Blocks
	blockHashId = make(map[string]uint32, 0)
)

func (w *WriterImpl) Write(file []byte, name string) error {
	var (
		blockSize = int(w.BlockSize)
		idx       = blockSize
		filePath  = path.Join(w.FileDir, w.FilePattern+"-"+name)
	)

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return fmt.Errorf("file %s exists", filePath)
	}

	for {
		if idx <= len(file)-1 {
			if err := w.storeBlock(file[idx-blockSize : idx]); err != nil {
				return fmt.Errorf("store block failed: %+v", err)
			}
		} else {
			endSlice := make([]byte, blockSize)
			endSlice = append(file[idx-blockSize:], endSlice...)
			if err := w.storeBlock(endSlice[:blockSize]); err != nil {
				return fmt.Errorf("store block failed: %+v", err)
			}
			break
		}
		idx += blockSize
	}

	err := os.WriteFile(filePath, blocks.toByte(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (w *WriterImpl) storeBlock(block []byte) error {
	log.Printf(string(block))

	hsha1 := sha1.Sum(block)
	id := uint32(len(blocks))

	blocks = append(blocks, Block{
		Hash: hsha1,
		Data: block,
		Id:   id,
	})

	return nil
}
