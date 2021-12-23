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

	err = writer.Write([]byte("Hello world"), "test")
	log.Print(err)
}
