package manager

import (
	"fmt"
	"github.com/safronovD/super-storage/pkg"
	storage "github.com/safronovD/super-storage/pkg/redis"
	"log"
	"strconv"
)

type DataFileManager interface {
	Write(fileID string, data []byte, blockSize uint, hashFun string) (*pkg.DataFile, error)
	Read(fileID string) ([]int, error)
}

type dataFileManagerImpl struct {
	storage           storage.HashStorage
	dataFile          *pkg.DataFile
	localStorage      string
	lastBlockPosition int
}

func NewDataFileManager(storage storage.HashStorage, fileLink string) DataFileManager {
	return &dataFileManagerImpl{
		storage:           storage,
		dataFile:          &pkg.DataFile{},
		localStorage:      fileLink,
		lastBlockPosition: 0,
	}
}

func (manager *dataFileManagerImpl) Write(fileID string, data []byte, blockSize uint, hashFun string) (*pkg.DataFile, error) {
	if exist, err := manager.storage.CheckFileID(fileID); exist == 1 {
		if err != nil {
			return nil, err
		}
		log.Printf("This file is alredy exist")
		return nil, fmt.Errorf("this file is already exist")
	}

	manager.dataFile = &pkg.DataFile{
		BlockSize:  blockSize,
		HashFun:    pkg.GetHashFun(hashFun),
		Blocks:     []pkg.Block{},
		PosBlockId: []int{},
	}

	var (
		blockSizeInt = int(blockSize)
		idx          = blockSizeInt
	)

	for {
		if idx <= len(data)-1 {
			if err := manager.add(data[idx-blockSizeInt:idx], blockSizeInt, fileID); err != nil {
				return nil, err
			}
		} else {
			endSlice := make([]byte, blockSize)
			endSlice = append(data[idx-blockSizeInt:], endSlice...)
			if err := manager.add(endSlice[:blockSize], blockSizeInt, fileID); err != nil {
				return nil, err
			}
		}
		idx += blockSizeInt
	}

	return manager.dataFile, nil
}

func (manager *dataFileManagerImpl) Read(fileID string) ([]int, error) {
	exist, err := manager.storage.CheckFileID(fileID)
	if err != nil {
		return nil, err
	}
	if exist != 1 {
		return nil, fmt.Errorf("this file doesn't exist %v", fileID)
	}

	fileMeta, err := manager.storage.Read(fileID)
	if err != nil {
		return nil, err
	}

	result := make([]int, 0)
	for _, elem := range fileMeta {
		position, err := strconv.Atoi(elem[0])
		if err != nil {
			position, err = strconv.Atoi(elem[1])
			if err != nil {
				return nil, err
			}
		}
		result = append(result, position)
	}
	return result, nil
}

func (manager *dataFileManagerImpl) add(data []byte, blockSizeInt int, fileID string) error {
	hash := manager.dataFile.HashFun.GetHash(data)
	exist, _, err := manager.storage.CheckHash(hash)
	if err != nil {
		return err
	}
	if exist != 1 {
		if err := manager.storage.WriteNewHash(hash, manager.localStorage, manager.lastBlockPosition); err != nil {
			return err
		}
		manager.dataFile.Add(data, manager.lastBlockPosition)
		manager.lastBlockPosition += blockSizeInt
	}
	if err := manager.storage.WriteBlockPosition(fileID, hash); err != nil {
		return err
	}
	return nil
}
