package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/scoutred/blds"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s datafile\n", os.Args[0])
		os.Exit(1)
	}

	fname := os.Args[1]

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	var records []blds.Record

	switch path.Ext(fname) {
	case ".csv":
		records, err = blds.FromCSV(file)
		if err != nil {
			log.Fatal(err)
		}
	case ".json":
		records, err = blds.FromJSON(file)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown file type ", path.Ext(fname))
	}

	for _, v := range records {
		// fmt.Println(v.EstProjCost)
		fmt.Println(v)
	}
}
