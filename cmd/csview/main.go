package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	var fname, col string

	switch len(os.Args) {
	default:
		fmt.Printf("usage: %s data.csv [col]\n", os.Args[0])
		os.Exit(1)
	case 3:
		col = os.Args[2]
		fallthrough
	case 2:
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rcsv := csv.NewReader(file)

	rec, err := rcsv.Read()
	if err != nil {
		log.Fatal(err)
	}

	cols := map[string]int{}
	for i, v := range rec {
		cols[v] = i
	}

	if _, ok := cols[col]; len(cols) > 0 && !ok {
		log.Fatalf("column %v not found", col)
	}

	for {
		rec, err := rcsv.Read()
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil {
			os.Exit(1)
		}

		if len(col) > 0 {
			i := cols[col]
			if len(rec) <= i {
				log.Println("record does not have column ", col)
				continue
			}

			fmt.Println(rec[i])
		} else {
			fmt.Println(rec)
		}
	}

}
