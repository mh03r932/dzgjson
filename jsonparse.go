package main

//https://build-your-own.org/blog/20231108_1week/

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

type ParseState struct {
	input string
	pos   int
}

func main() {
	file := flag.String("file", "", "file to parse")
	flag.Parse()
	fmt.Println(*file)

	if *file != "" {
		stringToParse := readFileData(*file)
		if !parseJSON(stringToParse) {
			log.Fatal("Failed to parse JSON")
		}
	}
}

func readFileData(filepath string) string {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File contents: %s\n", data)
	// contert to utf8 string and return
	return string(data)
}

func parseJSON(jsontext string) bool {
	state := &ParseState{
		input: jsontext,
		pos:   0,
	}

	// skip any whitespace
	parseWhitespace(state)
	err := parseValue(state)
	if err != nil {
		return false
	}
	parseWhitespace(state)

	if state.pos < len(state.input) {
		return false
	}

	return true
}

func parseWhitespace(state *ParseState) {
	for state.pos < len(state.input) && isWhitespace(rune(state.input[state.pos])) {
		state.pos++
	}
}

func parseValue(state *ParseState) error {
	if len(state.input[state.pos:]) >= 4 && state.input[state.pos:state.pos+4] == "true" {
		state.pos += 4
		return nil
	}
	if len(state.input[state.pos:]) >= 5 && state.input[state.pos:state.pos+5] == "false" {
		state.pos += 5
		return nil
	}
	if len(state.input[state.pos:]) >= 4 && state.input[state.pos:state.pos+4] == "null" {
		state.pos += 4
		return nil
	}

	return errors.New("Invalid JSON")
}

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
