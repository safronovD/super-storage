package main

import (
	"flag"
	"fmt"
	"github.com/safronovD/super-storage/pkg"
	"github.com/safronovD/super-storage/pkg/manager"
	"github.com/safronovD/super-storage/pkg/redis"
	"log"
)

const (
	fsPath = "/tmp/data"
)

var (
	config = &redis.Config{
		Hostname: "0.0.0.0:6379",
		Username: "",
		Password: "",
		HashDB:   0,
		FileDB:   1,
	}
)

// flags
var (
	modeR = flag.Bool("read", false, "read")
	modeW = flag.Bool("write", false, "write")
	modeD = flag.Bool("delete", false, "delete")
	modeL = flag.Bool("list", false, "list")
	modeC = flag.Bool("compare", false, "compare")

	path = flag.String("path", "", "path")
	name = flag.String("name", "", "name")
	size = flag.Uint("size", 256, "size")
	hash = flag.String("hash", "SHA256", "hash")

	in  = flag.String("in", "", "in")
	out = flag.String("out", "", "out")
)

func performWrite(name, path, hash string, blockSize uint) error {
	redisS, err := redis.NewRedisHashImpl(config)
	if err != nil {
		return err
	}

	fsManager := &pkg.FSManager{}

	writer, err := pkg.NewWriter(fsPath)
	if err != nil {
		return err
	}

	data, err := fsManager.ReadFile(path)
	if err != nil {
		return err
	}

	dfManager := manager.NewDataFileManager(redisS, name)
	df, err := dfManager.Write(name, data, blockSize, hash)
	if err != nil {
		return err
	}

	return writer.Write(df, name)
	if err != nil {
		return err
	}

	return nil
}

func performRead(name, path string, size uint) error {
	redisS, err := redis.NewRedisHashImpl(config)
	if err != nil {
		return err
	}

	fsManager := &pkg.FSManager{}
	writer, err := pkg.NewWriter(fsPath)
	if err != nil {
		return err
	}

	dfManager := manager.NewDataFileManager(redisS, name)
	positions, err := dfManager.Read(name)
	if err != nil {
		return err
	}

	result, err := writer.Read(name, int(size), positions)
	if err != nil {
		return err
	}

	err = fsManager.WriteFile(path, result)
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

func performCompare(in, out string) error {
	manager := &pkg.FSManager{}

	dataIn, err := manager.ReadFile(in)
	if err != nil {
		return err
	}

	dataOut, err := manager.ReadFile(out)
	if err != nil {
		return err
	}

	result := pkg.CompareBytes(dataIn, dataOut)
	fmt.Printf("Percent - %f%%\n", result)

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
		err = performRead(*name, *path, *size)
	case *modeD:
		err = performDelete(*name)
	case *modeL:
		err = performList()
	case *modeC:
		err = performCompare(*in, *out)

	}

	if err != nil {
		log.Fatalf("err: %+v", err)
	}
}
