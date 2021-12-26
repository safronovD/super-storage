package main

import (
	"fmt"
	"log"

	"github.com/safronovD/super-storage/pkg"
)

func main() {
	fmt.Println("Implement me pls")

	manager := &pkg.FSManager{}

	writer, err := pkg.NewWriter("/tmp/data")
	if err != nil {
		log.Fatal(err)
	}

	//data := []byte("HellHello worlddddffff")
	size := uint8(4)

	data, err := manager.ReadFile("/home/user/test/cosi.mp4")
	if err != nil {
		log.Fatal(err)
	}

	dataFile := pkg.CreateDataFile(data, size, "SHA1")

	err = writer.Write(dataFile, "test")
	log.Print(err)

	result, err := writer.Read(dataFile, "test")
	log.Print(err)

	//log.Printf("in: %s\n", string(data))
	//log.Printf("out: %s\n", string(result))
	log.Printf("%d = %d", len(data), len(result))

	err = manager.WriteFile("/home/user/test/cosi2.mp4", result)
	if err != nil {
		log.Fatal(err)
	}
}
