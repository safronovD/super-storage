package pkg

import (
	"crypto/sha1"
)

type Block struct {
	Data []byte
}

func (b *Block) toByte() []byte {
	return b.Data
}

type DataFile struct {
	blockSize uint8
	blocks    []Block
	// Id - position in blocks array
	blockHashId map[interface{}]int
	posBlockId  map[int]int
}

func CreateDataFile(data []byte, blockSize uint8) *DataFile {
	var (
		blockSizeInt = int(blockSize)
		idx          = blockSizeInt
		pos          = 0
		df           = &DataFile{
			blockSize:   blockSize,
			blocks:      []Block{},
			blockHashId: map[interface{}]int{},
			posBlockId:  map[int]int{},
		}
	)

	for {
		if idx <= len(data)-1 {
			df.add(data[idx-blockSizeInt:idx], pos)
		} else {
			endSlice := make([]byte, blockSize)
			endSlice = append(data[idx-blockSizeInt:], endSlice...)
			df.add(endSlice[:blockSize], pos)
			break
		}
		idx += blockSizeInt
		pos++
	}

	return df
}

func (df *DataFile) add(data []byte, pos int) {
	var id int

	hsha1 := sha1.Sum(data)

	if val, ok := df.blockHashId[hsha1]; ok {
		id = val
	} else {
		id = len(df.blocks)
		df.blocks = append(df.blocks, Block{data})
		df.blockHashId[hsha1] = id
	}

	df.posBlockId[pos] = id
}

func (df *DataFile) toByte() (bytes []byte) {
	bytes = make([]byte, 0, len(df.blocks)*int(df.blockSize))
	for _, block := range df.blocks {
		bytes = append(bytes, block.toByte()...)
	}
	return
}