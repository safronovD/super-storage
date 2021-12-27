package pkg

type Block struct {
	Data []byte
}

func (b *Block) toByte() []byte {
	return b.Data
}

type DataFile struct {
	BlockSize  uint
	HashFun    HashFun
	Blocks     []Block
	PosBlockId []int
}

func (df *DataFile) Add(data []byte, posBlockID int) {
	df.Blocks = append(df.Blocks, Block{data})
	df.PosBlockId = append(df.PosBlockId, posBlockID)
}

func (df *DataFile) toByte() (bytes []byte) {
	bytes = make([]byte, 0, len(df.Blocks)*int(df.BlockSize))
	for _, block := range df.Blocks {
		bytes = append(bytes, block.toByte()...)
	}
	return
}
