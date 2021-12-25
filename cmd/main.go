package main

import (
	"fmt"
	"log"

	"github.com/safronovD/super-storage/pkg"
)

func main() {
	fmt.Println("Implement me pls")

	writer, err := pkg.NewWriter("/tmp/data")
	if err != nil {
		log.Fatal(err)
	}

	data := []byte("HelloHello world")
	size := uint8(4)

	dataFile := pkg.CreateDataFile(data, size)

	err = writer.Write(dataFile, "test")
	log.Print(err)

	_, err = writer.Read(dataFile, "test")
	log.Print(err)
}
