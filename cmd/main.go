package main

import (
	"flag"
	"github.com/safronovD/super-storage/pkg"
	"log"
)

const (
	fsPath = "/tmp/data"
)

// flags
var (
	modeR = flag.Bool("read", false, "read")
	modeW = flag.Bool("write", false, "write")
	modeD = flag.Bool("delete", false, "delete")
	modeL = flag.Bool("list", false, "list")

	path = flag.String("path", "", "path")
	name = flag.String("name", "", "name")
	size = flag.Uint("size", 256, "size")
	hash = flag.String("hash", "SHA256", "hash")
)

func performWrite(name, path, hash string, blockSize uint) error {
	manager := &pkg.FSManager{}

	writer, err := pkg.NewWriter(fsPath)
	if err != nil {
		return err
	}

	data, err := manager.ReadFile(path)
	if err != nil {
		return err
	}

	dataFile := pkg.CreateDataFile(data, blockSize, hash)

	return writer.Write(dataFile, name)
	if err != nil {
		return err
	}

	return nil
}

func performRead(name, path string) error {
	manager := &pkg.FSManager{}
	writer, err := pkg.NewWriter(fsPath)
	if err != nil {
		return err
	}

	//get DataFile FromDB
	dataFile := &pkg.DataFile{}

	result, err := writer.Read(dataFile, name)
	if err != nil {
		return err
	}

	err = manager.WriteFile(path, result)
	if err != nil {
		return err
	}

	return nil
}

func performDelete(name string) error {
	// delete file from DB
	return nil
}

func performList() error {
	// print file list from DB
	return nil
}

func main() {
	var (
		err error
	)

	flag.Parse()

	switch {
	case *modeW:
		err = performWrite(*name, *path, *hash, *size)
	case *modeR:
		err = performRead(*name, *path)
	case *modeD:
		err = performDelete(*name)
	case *modeL:
		err = performList()
	}

	if err != nil {
		log.Fatalf("err: %+v", err)
	}
}
