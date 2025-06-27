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
func getLookahead(state *ParseState) string {
	if state.pos >= len(state.input) {
		log.Fatalf("Unexpected end of input")
	}
	return state.input[state.pos : state.pos+1]
}

func consume(state *ParseState, c string) {
	if getLookahead(state) == c {
		state.pos = state.pos + len(c)
	} else {
		log.Fatalf("Expected '%s' but found '%s'", c, getLookahead(state))
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

	if getLookahead(state) == "\"" {
		parseString(state)
		return nil
	}
	if getLookahead(state) == "[" {
		paraseArray(state)
		return nil
	}
	if getLookahead(state) == "{" {
		parseObject(state)
		return nil
	}
	if isValidNumberChar(rune(getLookahead(state)[0])) {
		paraseNumber(state)
		return nil
	}

	return errors.New("Invalid JSON")
}

func parseObject(state *ParseState) {
	consume(state, "{")
	parseWhitespace(state)
	lookahead := getLookahead(state)
	if lookahead == "}" {
		consume(state, "}")
		return
	}
	parseMembers(state)
	consume(state, "}")

}

func parseMembers(state *ParseState) {
	parseMember(state)
	lookahead := getLookahead(state)
	if lookahead == "," {
		consume(state, ",")
		parseMembers(state)
	} else if lookahead == "}" {
		consume(state, "}")
		return
	} else {
		log.Fatalf("Expected ',' or '}' but found '%s'", lookahead)
	}
}

func parseMember(state *ParseState) {
	parseString(state)
	parseString(state)
	parseWhitespace(state)
	consume(state, ":")
	parseElement(state)
}

func paraseArray(state *ParseState) {
	consume(state, "[")
	lookahead := getLookahead(state)
	if isWhitespace(rune(lookahead[0])) {
		parseWhitespace(state)
	}
	if getLookahead(state) == "]" {
		consume(state, "]")
		return
	}

	parseElements(state)

	if getLookahead(state) != "]" {
		log.Fatalf("Expected ']' but found '%s'", getLookahead(state))
	}
	consume(state, "]")

}

func parseElements(state *ParseState) {
	parseElement(state)
	lookahead := getLookahead(state)
	if lookahead == "," {
		consume(state, ",")
		parseElements(state)
	}
}

func parseElement(state *ParseState) {
	parseWhitespace(state)
	parseValue(state)
	parseWhitespace(state)
}

func paraseNumber(state *ParseState) {
	parseInteger(state)
	parseFraction(state)
	parseExponet(state)
}

func parseInteger(state *ParseState) {
	maybeParseMinus(state)
	lookahead := getLookahead(state)
	// cast lookahead to rune
	r := rune(lookahead[0])
	if isValidDigit(r) {
		parseDigit(state)
	} else {
		parseOneNine(state)
		parseDigits(state)
	}

}

func parseFraction(state *ParseState) {
	lookahead := getLookahead(state)
	if lookahead == "" {
		return
	}
	if lookahead == "." {
		consume(state, ".")
		parseDigits(state)

	}
}

func parseExponet(state *ParseState) {
	lookahead := getLookahead(state)
	// cast lookahead to rune
	if lookahead == "e" || lookahead == "E" {
		consume(state, lookahead)
		parseSign(state)
		parseDigits(state)
	} else {
		return
	}
}

func parseSign(state *ParseState) {
	lookahead := getLookahead(state)
	// cast lookahead to rune
	r := rune(lookahead[0])
	if isValidMinusChar(r) {
		consume(state, string(r))
	} else if r == '+' {
		consume(state, "+")
	} else if lookahead == "" {
		return
	} else {
		log.Fatalf("Expected '+' or '-' but found '%s'", lookahead)
	}
}

func maybeParseMinus(state *ParseState) {
	lookahead := getLookahead(state)
	// cast lookahead to rune
	r := rune(lookahead[0])
	if isValidMinusChar(r) {
		consume(state, string(r))
	}

}

func parseString(state *ParseState) {
	consume(state, "\"")
	parseCharacters(state)
	consume(state, "\"")

}

func parseCharacters(state *ParseState) {
	if getLookahead(state) == "\"" {
		return
	} else {
		parseCharacter(state)
		parseCharacters(state)
	}

}

func parseCharacter(state *ParseState) {
	lookahead := getLookahead(state)
	// cast lookahead to rune
	r := rune(lookahead[0])
	if isAllowedChar(r) {
		consume(state, string(r))
	} else {
		consume(state, "\\")
		paraseEscape(state)
	}
}

func paraseEscape(state *ParseState) {
	lookahead := getLookahead(state)
	switch lookahead {
	case "\"":
		consume(state, "\"")
	case "\\":
		consume(state, "\\")
	case "/":
		consume(state, "/")
	case "b":
		consume(state, "b")
	case "f":
		consume(state, "f")
	case "n":
		consume(state, "n")
	case "r":
		consume(state, "r")
	case "t":
		consume(state, "t")
	case "u":
		consume(state, "u")
		parseHex(state)
		parseHex(state)
		parseHex(state)
		parseHex(state)

	}

}

func parseHex(state *ParseState) {
	lookahead := getLookahead(state)
	r := rune(lookahead[0])
	if isValidHexChar(r) {
		consume(state, lookahead)
	} else {
		parseDigit(state)
	}

}

func parseDigits(state *ParseState) {
	parseDigit(state)
	lookahead := getLookahead(state)
	r := rune(lookahead[0])
	if isValidDigit(r) {
		parseDigits(state)
	} else {
		return
	}
}

func parseDigit(state *ParseState) {
	lookahead := getLookahead(state)
	if lookahead == "0" {
		consume(state, "0")
	} else {
		parseOneNine(state)
	}

}

func parseOneNine(state *ParseState) {
	lookahead := getLookahead(state)
	r := rune(lookahead[0])
	if isValidOneNine(r) {
		consume(state, lookahead)
	} else {
		log.Fatalf("Expected a digit between 1 and 9 but found '%s'", lookahead)
	}

}

func isValidDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isValidOneNine(r rune) bool {
	return r >= '1' && r <= '9'
}

func isValidHexChar(r rune) bool {
	return r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}
func isValidNumberChar(r rune) bool {
	return isValidDigit(r) || r == '-'
}

func isValidMinusChar(r rune) bool {
	return r == '-'
}

func isAllowedChar(r rune) bool {
	return r >= 0x20 && r <= 0x10FFFF && r != '"' && r != '\\'
}
func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
