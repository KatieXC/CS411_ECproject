package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

const NUMARGS int = 6

type join_method int

const (
	NESTED_LOOP join_method = iota
	HASH
)

type table_data struct {
	Reader  *csv.Reader
	Headers []string
	JoinCol int
}

var err error
var method join_method
var table1 table_data
var table2 table_data
var outputFile *os.File

func check(err error) {
	if err != nil {
		log.Fatal("Error occurred!\n")
	}
}

func setup_tables(filename1 string, joinColName1 string, filename2 string, joinColName2 string) {
	fmt.Printf("Args: %s with col %s, %s with col %s\n", filename1, joinColName1, filename2, joinColName2)
	f1, err := os.Open(filename1)
	if err != nil {
		log.Fatal("CSV File 1 could not be opened.\n")
	}
	f2, err := os.Open(filename2)
	if err != nil {
		log.Fatal("CSV File 2 could not be opened.\n")
	}
	defer f1.Close()
	defer f2.Close()

	table1.Reader = csv.NewReader(f1)
	table2.Reader = csv.NewReader(f2)

	table1.Headers, err = table1.Reader.Read()
	check(err)
	for i, h := range table1.Headers {
		if h == joinColName1 {
			table1.JoinCol = i
		}
	}
	table2.Headers, err = table2.Reader.Read()
	check(err)
	for i, h := range table2.Headers {
		if h == joinColName2 {
			table2.JoinCol = i
		}
	}
}

func main() {
	// Read arguments
	argv := os.Args[1:]
	if len(argv) != NUMARGS {
		log.Fatal("usage: go run joins.go <input1.csv> <join_col_name1> <input2.csv> <join_col_name2> <join method> <output.csv>\n")
	}

	if argv[4] == "HASH" {
		method = HASH
	} else if argv[4] == "NESTED_LOOP" {
		method = NESTED_LOOP
	} else {
		log.Fatal("Invalid join type.\n")
	}

	setup_tables(argv[0], argv[1], argv[2], argv[3])
	outputFile, err := os.Create(argv[5])
	if err != nil {
		log.Fatal("Output file could not be created.\n")
	}
	defer outputFile.Close()

	fmt.Printf("End of Main\n")
}
