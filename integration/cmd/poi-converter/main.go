package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	Pipeline()
}

// Pipeline extracts result from csv and convert to json format
//
// Optimization for future: Pipeline could be a generic framework
//
// Extraction
// Extraction part load data from sources which could has many kind of readers(json, csv, etc)
// Extraction put loaded data into a channel
//
// Transformation
// Transformation part read data from Extraction's result, doing things(user
// defined functions) and output result to another channel
// There might be trick logic in this step when input information is not separated
// by lines:
// Input string might be: last part for object 1 | first part for object 2
// It need to implement buffer to handling this
// Transformation part could have a pool of worker to execute
//
// Loader/Publisher
// Loader part read data from channel and writes result to target format
// Similar as reader, it could have many writers
func Pipeline() {
	// extract: load data from csv
	csvFile, err := os.Open(flags.inputPath)
	if err != nil {
		glog.Fatalf("While open file %s, met error %v", flags.inputPath, err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	content, _ := reader.ReadAll()
	glog.Infof("Finish loading file of %s, it contains %d line of data\n", flags.inputPath, len(content))

	if len(content) < 1 {
		glog.Fatalf("No content in given file %s\n", flags.inputPath)
	}

	// transformation: convert csv to json format
	header := make([]string, 0)
	for _, attr := range content[0] {
		header = append(header, attr)
	}
	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, line := range content {
		buffer.WriteString("{")
		for j, element := range line {
			buffer.WriteString(`"` + header[j] + `":`)

			_, isFloatErr := strconv.ParseFloat(element, 64)
			_, isBoolErr := strconv.ParseBool(element)
			if isFloatErr == nil {
				buffer.WriteString(element)
			} else if isBoolErr == nil {
				buffer.WriteString(strings.ToLower(element))
			} else {
				buffer.WriteString((`"` + element + `"`))
			}

			if j < len(line)-1 {
				buffer.WriteString(",")
			}
		}

		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(`]`)
	glog.Info("Finish converting from csv to internal json string\n")

	// Load: save content to target file
	rawMessage := json.RawMessage(buffer.String())
	if err := ioutil.WriteFile(flags.outputPath, rawMessage, os.FileMode(0644)); err != nil {
		glog.Fatalf("While writing result to file %s, met error %v", flags.outputPath, err)
	}
	glog.Infof("Finish generating target file %s\n", flags.outputPath)

}
