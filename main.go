package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

//
// TODO
// rename project to markdownTableRoller
// handle number ranges eg 1-4, 5-6 for varying probabilities, probably with table syntax
// handle tables with same name in multiple directories. print paths and one random selection?
// work on windows
// if [subtable] is not found just print [subtable]
// add more maze rats tables
// update README
// add usage string
// create a "markdown" interface incase it comes time to add a 3rd party md package
//

var dotSlash string = "." + string(os.PathSeparator)

func main() {
	args := os.Args

	query, err := parseArgs(args)
	checkError(err, "Bad command argument")

	path, err := findTable(query, dotSlash)
	checkError(err, "Error finding file")

	file, err := os.Open(path)
	checkError(err, "Error reading file")
	defer file.Close()

	tableValues, err := parseTableValues(file)
	checkError(err, "Error parsing table values")

	result := rollTheDice(tableValues)
	fmt.Println(result)
}

func parseArgs(args []string) (query string, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("Please provide a table name")
	}

	if slices.Contains([]string{"-h", "--h", "-help", "--help", "\\h", "\\help"}, args[1]) {
		printUsageAndExit()
	}

	query = strings.TrimPrefix(args[1], dotSlash)
	query = strings.TrimSuffix(query, ".md")

	return query, nil
}

func findTable(search string, dir string) (tablePath string, err error) {
	if search == "" {
		return "", fmt.Errorf("Please provide a table name. Search: %s, Directory: %s", search, dir)
	}

	search = standardizeSearch(search)
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if tablePath != "" {
			return nil
		}
		if strings.Contains(strings.ToLower(path), search) {
			tablePath = path
			return nil
		}
		return nil
	})
	if tablePath == "" {
		return "", fmt.Errorf("Table not found: %s", search)
	}
	return tablePath, err
}

func standardizeSearch(search string) string {
	search = filepath.FromSlash(search)
	search = strings.ToLower(search)
	if !strings.HasSuffix(search, ".md") {
		search = search + ".md"
	}
	return search
}

func parseTableValues(tableFile *os.File) ([]string, error) {
	scanner := bufio.NewScanner(tableFile)
	var tableValues []string
	for scanner.Scan() {
		matcher, err := regexp.Compile(`^(\d+\. |\* )`) // Ex. starts with '* ' or '1. '
		checkError(err, "Error parsing table values")
		isValue := matcher.Match([]byte(scanner.Text()))
		if isValue {
			noPrefix := matcher.ReplaceAllString(scanner.Text(), "")
			// TODO Do this after the result is picked to avoid extra work
			str := checkForSubQueries(noPrefix)

			tableValues = append(tableValues, str)

		}
	}
	if len(tableValues) == 0 {
		return nil, fmt.Errorf("No table values found")
	}
	return tableValues, nil
}

func checkForSubQueries(tableRow string) string {
	linkMatcher := regexp.MustCompile(`\[.+?\]\(.+?\)`)
	prefixMatcher := regexp.MustCompile(`\[.+?\]\(`)

	hasSubQueries := linkMatcher.Match([]byte(tableRow))

	if hasSubQueries {

		subQueries := linkMatcher.FindAllString(tableRow, -1)
		for _, subQueryLink := range subQueries {
			subQuery := prefixMatcher.ReplaceAllString(subQueryLink, "")
			subQuery = strings.TrimSuffix(subQuery, ")")
			subQueryFile, err := findTable(subQuery, dotSlash)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Couldn't find table for sub query: %s\n", subQuery)
				return tableRow
			}
			subQueryFileReader, err := os.Open(subQueryFile)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Couldn't read file for sub query: %s\n", subQuery)
				return tableRow
			}
			subQueryValues, err := parseTableValues(subQueryFileReader)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Couldn't parse values for sub query: %s\n", subQuery)
				return tableRow
			}
			subQueryResult := rollTheDice(subQueryValues)
			tableRow = strings.Replace(tableRow, subQueryLink, subQueryResult, 1)
		}
	}
	return tableRow
}

func rollTheDice(tableValues []string) string {
	rand.Seed(time.Now().UnixNano())
	return tableValues[rand.Intn(len(tableValues))]
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsageAndExit() {
	fmt.Println("Usage: gotableroller {TableName}\nTableName: the name of the markdown file containing the table. " +
		"This file must exist the same directory or a subdirectory of gotableroller. TableName may/maynot contain" +
		"the '.md' extension. It may contain path components as while. Examples: 'Weapons', 'weapons', 'weapons.md', " +
		"'Items/Weapons.md'")
	os.Exit(0)
}
