package main

import (
	"log"
	"time"

	"github.com/sura8257/skbn/pkg/skbn"
)

func main() {
	src := "path/to/copy/from"
	dst := "s3://bucket/path/to/copy/to"
	parallel := 0     
	bufferSize := 0 

	start := time.Now()
	if err := skbn.Copy(src, dst, parallel, bufferSize); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	log.Printf("Copy execution time: %s", elapsed)
}
