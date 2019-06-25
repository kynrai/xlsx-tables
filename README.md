# XLSX Tables

## Motivation

XLSX files are often used to store and manage large tables. A nice free GUI for some, but a nightmare for golang devs asked to process several gigabytes of data in a spreadsheet. Existiing tools for go often load sheets into memory which make it very slow to process. The goal of this project is to provide a encoding/csv like interface over xslx files that store a table. We focus on streaming data instead of loading into RAM.

## Limitations

This library assumes that a given sheet in an xlsx format spreadsheet has nothing but a table in it. The equivelant of a CSV but in a spreadsheet. This library can can only read line by line for processing or loading into databases.

## Usage

```go
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"strings"

	xlsx "github.com/kynrai/xlsx-tables"
)

func main() {
	f, err := zip.OpenReader("xlsx_sample.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := xlsx.NewReader(f)
	r.Worksheet = "sheet2"

	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(strings.Join(row, ","))
	}
}
```
