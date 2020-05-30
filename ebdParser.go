package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const dataLoc = "/Users/katjad/Documents/eBird/ebd_201701_202005_prv_relApr-2020.zip"
const fileName = "ebd_201701_202005_prv_relApr-2020.txt"
const outputDirectory = "/Users/katjad/Documents/eBird/Output"

// function that can be passed a criterion by which observations are agregated
// The function opens each line of the data file and applies f to it
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
			//discard the headers
			lineReader.ReadLine()
			for i := 0; ; i++ {
				if i%1000000 == 0 {
					fmt.Println(i / 1000000)
				}
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

func main() {
	fmt.Println(time.Now())
	byCountry := make(map[string]map[int64]int)
	observationParser(func(observation []string) {
		if len(observation) < 30 {
			return
		}
		date, err := time.Parse("2006-01-02", observation[27])
		if err != nil {
			return
		}
		countryCode := observation[13]
		if _, ok := byCountry[countryCode]; !ok {
			byCountry[countryCode] = make(map[int64]int)
		}
		byCountry[countryCode][date.Unix()]++

	})
	jsonByCountryData, err := json.Marshal(byCountry)
	if err != nil {
		panic(err)
	}
	jsonByCountryFile, err := os.Create(strings.Join([]string{outputDirectory, "byCountry.json"}, "/"))
	if err != nil {
		panic(err)
	}
	defer jsonByCountryFile.Close()

	jsonByCountryFile.Write(jsonByCountryData)
	fmt.Println(time.Now())
}
