package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strings"
)

const dataLoc = "/Users/katjad/Documents/eBird/ebd_201701_202005_prv_relApr-2020.zip"
const fileName = "ebd_201701_202005_prv_relApr-2020.txt"

func observationParser(f func([]string)) {
	// open reader for zipped folder
	zipReader, err := zip.OpenReader(dataLoc)
	if err != nil {
		panic(err)
	}
	// close at the end
	defer zipReader.Close()

	// open file in the zip folder
	for _, file := range zipReader.File {
		if file.Name == fileName {
			fileReader, err := file.Open()
			if err != nil {
				panic(err)
			}
			// open document line by line
			lineReader := bufio.NewReader(fileReader)
			for i := 0; i < 1000000; i++ {
				line, _, err := lineReader.ReadLine()
				// test for EOF
				if err == io.EOF {
					break
				} else if err != nil {
					panic(err)
				}
				observation := strings.Split(string(line), "\t")
				f(observation)
			}
			fileReader.Close()
		}

	}

}

func testFunction(lst []string) {
	fmt.Println(lst[0])
}

func main() {
	observationParser(testFunction)
}
