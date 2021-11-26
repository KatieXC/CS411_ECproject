package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

const NUMARGS int = 6

type join_method int

const (
	NESTED_LOOP join_method = iota
	HASH
)

type table_data struct {
	Headers []string
	Data    map[string][]string // key = col name, value = [col value at row 0, col value at row 1, ...]
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

func print_table_start(table *table_data) {
	for _, h := range table.Headers {
		fmt.Printf("Column: %s\n", h)
		for i, val := range table.Data[h] {
			if i >= 5 {
				break
			}
			fmt.Printf("%s\n", val)
		}
	}
}

func setup_table(table *table_data, filename string, joinColName string) {
	fmt.Printf("Args: %s with col %s\n", filename, joinColName)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("CSV file could not be opened.\n")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	table.Headers = strings.Split(scanner.Text(), "|")
	for i, h := range table.Headers {
		// fmt.Printf("Header Found: %s\n", h)
		table.Data[h] = []string{}
		if strings.EqualFold(h, joinColName) { // Case Insensitive column check
			table.JoinCol = i
		}
	}

	for scanner.Scan() {
		row_str := scanner.Text()
		// fmt.Printf("Row = %s\n", row_str)
		row := strings.Split(row_str, "|")
		for i, value := range row[:len(row)-1] {
			table.Data[table.Headers[i]] = append(table.Data[table.Headers[i]], value)
		}
	}
	print_table_start(table)
}

func run_nested_loop(output_filename string) {
	joinHeader1 := table1.Headers[table1.JoinCol]
	joinHeader2 := table2.Headers[table2.JoinCol]

	outputFile, err := os.OpenFile(output_filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	check(err)
	out_writer := csv.NewWriter(outputFile)
	defer func() {
		out_writer.Flush()
		err = out_writer.Error()
		check(err)
		outputFile.Close()
	}()

	fmt.Println("Starting Timer...")
	// Start timer

	for i, r := range table1.Data[joinHeader1] {
		for j, s := range table2.Data[joinHeader2] {
			if r == s {
				// fmt.Printf("Creating row on %s==%s\n", r, s)
				var out_row []string
				for _, h1 := range table1.Headers {
					out_row = append(out_row, table1.Data[h1][i])
				}
				for hidx, h2 := range table2.Headers {
					if hidx == table2.JoinCol {
						continue
					}
					out_row = append(out_row, table2.Data[h2][j])
				}
				// for _, row_val := range out_row {
				// 	fmt.Printf("%s ", row_val)
				// }
				out_writer.Write(out_row)
			}
		}
	}

	// End timer
	fmt.Println("...Timer Ended")
}

func run_hash(output_filename string) {

}

func main() {
	// Read arguments
	argv := os.Args[1:]
	if len(argv) != NUMARGS {
		log.Fatal("usage: go run joins.go <input1.csv> <join_col_name1> <input2.csv> <join_col_name2> <join method> <output.csv>\n")
	}

	table1.Data = make(map[string][]string)
	table2.Data = make(map[string][]string)
	setup_table(&table1, argv[0], argv[1])
	setup_table(&table2, argv[2], argv[3])

	if strings.EqualFold(argv[4], "HASH") {
		run_hash(argv[5])
	} else if strings.EqualFold(argv[4], "NESTED_LOOP") {
		run_nested_loop(argv[5])
	} else {
		log.Fatal("Invalid join method request: 'HASH' or 'NESTED_LOOP'.")
	}

	fmt.Printf("End of Main\n")
}
