package pkg

type Block struct {
	Data []byte
}

func (b *Block) toByte() []byte {
	return b.Data
}

type DataFile struct {
	blockSize uint8
	hashFun   HashFun
	blocks    []Block
	// Id - position in blocks array
	blockHashId map[interface{}]int
	posBlockId  []int
}

func CreateDataFile(data []byte, blockSize uint8, hashFun string) *DataFile {
	var (
		blockSizeInt = int(blockSize)
		idx          = blockSizeInt
		pos          = 0
		df           = &DataFile{
			blockSize:   blockSize,
			hashFun:     getHashFun(hashFun),
			blocks:      []Block{},
			blockHashId: map[interface{}]int{},
			posBlockId:  []int{},
		}
	)

	for {
		if idx <= len(data)-1 {
			df.add(data[idx-blockSizeInt : idx])
		} else {
			endSlice := make([]byte, blockSize)
			endSlice = append(data[idx-blockSizeInt:], endSlice...)
			df.add(endSlice[:blockSize])
			break
		}
		idx += blockSizeInt
		pos++
	}

	return df
}

func (df *DataFile) add(data []byte) {
	var id int

	hsha1 := df.hashFun.getHash(data)

	if val, ok := df.blockHashId[hsha1]; ok {
		id = val
	} else {
		id = len(df.blocks)
		df.blocks = append(df.blocks, Block{data})
		df.blockHashId[hsha1] = id
	}

	df.posBlockId = append(df.posBlockId, id)
}

func (df *DataFile) toByte() (bytes []byte) {
	bytes = make([]byte, 0, len(df.blocks)*int(df.blockSize))
	for _, block := range df.blocks {
		bytes = append(bytes, block.toByte()...)
	}
	return
}
