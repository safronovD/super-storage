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

	data := []byte("HellHello worlddddffff")
	size := uint8(4)

	dataFile := pkg.CreateDataFile(data, size)

	err = writer.Write(dataFile, "test")
	log.Print(err)

	result, err := writer.Read(dataFile, "test")
	log.Print(err)

	log.Printf("in: %s\n", string(data))
	log.Printf("out: %s\n", string(result))
}
