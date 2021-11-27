package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const NUMARGS int = 6
const STEPLENGTH int = 20 // lineitem |><| orders -> set to 20
// (Sticking to rows in R * rows in S ~= 2x10^10)

func check(msg string, err error) {
	if err != nil {
		log.Fatal(msg)
	}
}

func setup_table(filename string, joinCol string) ([][]string, int) {
	var table [][]string
	row_count := 0
	joinIdx := -1

	csv_file, err := os.Open(filename)
	check("CSV file 1 could not be opened.\n", err)
	defer csv_file.Close()
	scanner := bufio.NewScanner(csv_file)

	for scanner.Scan() {
		if row_count == 0 || row_count%STEPLENGTH == 0 {
			// Choose delimiter
			table = append(table, strings.Split(scanner.Text(), ","))
			// table = append(table, strings.Split(scanner.Text(), "|"))
		}
		row_count++
	}
	for i, h := range table[0] {
		if h == joinCol {
			joinIdx = i
		}
	}

	fmt.Printf("Table: %s joining on column: %s (idx=%d)\n#Rows = %d\n", filename, joinCol, joinIdx, len(table)-1)
	return table, joinIdx
}

func main() {
	// Read arguments
	argv := os.Args[1:]
	if len(argv) != NUMARGS {
		log.Fatal("usage: go run joins.go <input1.csv> <join_col_name1> <input2.csv> <join_col_name2> <join method> <output.csv>\n")
	}
	filename1 := argv[0]
	filename2 := argv[2]
	joinCol1 := argv[1]
	joinCol2 := argv[3]
	joinMethod := argv[4]
	outputFilename := argv[5]

	/* Setup Tables */
	table1, joinIdx1 := setup_table(filename1, joinCol1)
	table2, joinIdx2 := setup_table(filename2, joinCol2)
	if joinIdx1 == -1 || joinIdx2 == -1 {
		log.Fatal("Join Column does not exist.\n")
	}

	/* Output File Setup */
	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	check("Output file couldn't be made.\n", err)
	out_writer := csv.NewWriter(outputFile)
	defer func() {
		out_writer.Flush()
		err = out_writer.Error()
		check("Output Writing Error.\n", err)
		outputFile.Close()
	}()

	/* Perform Join */
	if strings.EqualFold(joinMethod, "HASH") {
		fmt.Println("Timer Starting...")
		start := time.Now() // Time join operation

		// Hash Join
		hashmap := make(map[string][]string)
		for _, r := range table1 {
			hashmap[r[joinIdx1]] = r
		}
		for _, s := range table2 {
			r, check_id := hashmap[s[joinIdx2]]
			if check_id {
				out_writer.Write(append(r, s...))
			}
		}

		duration := time.Since(start)
		fmt.Println("...Timer Ended")
		fmt.Println(duration)
	} else if strings.EqualFold(joinMethod, "NESTED_LOOP") {
		fmt.Println("Timer Starting...")
		start := time.Now() // Time join operation

		// Nested Loop Join
		for _, r := range table1 {
			for _, s := range table2 {
				if r[joinIdx1] == s[joinIdx2] {
					out_writer.Write(append(r, s...))
				}
			}
		}

		duration := time.Since(start)
		fmt.Println("...Timer Ended")
		fmt.Println(duration)
	} else {
		log.Fatal("Invalid join method request: 'HASH' or 'NESTED_LOOP'.")
	}
	fmt.Printf("End of Join\n")
}
