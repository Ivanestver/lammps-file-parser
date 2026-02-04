package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"lammps-file-parser/parser"
	"os"
)

func main() {
	infilePtr := flag.String("infile", "", "input lammps file with data")
	outfilePtr := flag.String("outfile", "", "output lammps file with data")
	flag.Parse()
	if len(*infilePtr) == 0 {
		fmt.Println("Wrong infile flag usage")
		return
	}
	if len(*outfilePtr) == 0 {
		fmt.Println("Wrong outfile flag usage")
		return
	}
	content, err := getFileContent(*infilePtr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lammpsStruct, err := parser.Parse(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	parsedJson, err := json.Marshal(lammpsStruct)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = writeJson(parsedJson, *outfilePtr); err == nil {
		fmt.Println("Done!")
	} else {
		fmt.Println(err.Error())
	}
}

func writeJson(parsedJson []byte, outfile string) error {
	return os.WriteFile(outfile, parsedJson, os.ModePerm)
}

func getFileContent(filename string) (string, error) {
	if content, err := os.ReadFile(filename); err == nil {
		return string(content), nil
	} else {
		return "", err
	}
}
