package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadCSVFile(t *testing.T) {
	reader, err := ReadCSVFile("/home/seafooler/ethTx_data/0to999999_NormalTransaction/0to999999_NormalTransaction.csv")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; ; i++{
		line, err := reader.Read()
		fmt.Printf("The %dth record: %s\n", i, line)
		if err != nil {
			break
		}
	}
}

func TestMakeFullFileName(t *testing.T) {
	fullFileName := MakeFullFileName(0)
	if strings.Compare("/home/seafooler/ethTx_data/0to999999_NormalTransaction/0to999999_NormalTransaction.csv",
		fullFileName) != 0 {
		t.Fatal("MakeFullFileName function makes a wrong file name")
	}
}
