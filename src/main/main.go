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

	"IPutOatsInGoats/gotableroller/src/rollabletable"
)

// TODO
// rename project to markdownTableRoller
// Monsters return nothing
var (
	rowRangePattern  = regexp.MustCompile(`(\d+).(\d+)`)      // matches roll ranges like '5-12' and captures the numbers as groups
	markdownListItem = regexp.MustCompile(`^(\d+\. |\* )`)    // identifies a line as a markdown list item, ie. '1. ' or '* '
	linkMatcher      = regexp.MustCompile(`\[.+?\]\((.+?)\)`) // matches markdown links like '[Link Label](path/to/table)' with a group for 'path/to/table'
)

func main() {
	args := os.Args

	query, err := parseArgs(args)
	checkError(err, "Bad command argument")

	rollTable := createRollableTable(query)
	rand.Seed(time.Now().UnixNano())
	result := rollTable.Roll()
	for len(linkMatcher.FindStringSubmatch(result)) != 0 {
		link := getLinkFromResult(result)
		subTable := createRollableTable(link[1])
		subResult := subTable.Roll()
		result = strings.Replace(result, link[0], subResult, 1)
	}
	fmt.Println(result)

}

func getLinkFromResult(result string) []string {
	query := linkMatcher.FindStringSubmatch(result)
	return query
}

func createRollableTable(query string) rollabletable.RollableTable {
	query = standardizeSearch(query)

	path, err := findTable(query, ".")
	checkError(err, "Error finding file")
	file, err := os.Open(path)
	checkError(err, "Error reading file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rollTable, err := rollabletable.ParseRollableTable(*scanner)
	checkError(err, "Error parsing table")

	return rollTable
}

func parseArgs(args []string) (query string, err error) {
	if len(args) < 2 {
		return "", fmt.Errorf("Please provide a table name")
	}

	if len(args) > 2 {
		return "", fmt.Errorf("Unknown options: %v", args[2:])
	}

	if contains([]string{"-h", "--h", "-help", "--help", "\\h", "\\help"}, args[1]) {
		printUsageAndExit()
	}

	return args[1], nil
}

func standardizeSearch(search string) string {
	search = strings.TrimPrefix(search, "./")
	search = strings.TrimPrefix(search, ".\\")
	search = filepath.FromSlash(search)
	search = strings.ToLower(search)
	if !strings.HasSuffix(search, ".md") {
		search = search + ".md"
	}
	return search
}

func findTable(search string, dir string) (tablePath string, err error) {
	if search == "" {
		return "", fmt.Errorf("Please provide a table name. Search: %s, Directory: %s", search, dir)
	}

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		checkError(err, "hereIAM")
		if tablePath != "" {
			return nil
		}
		fmt.Println("Path: " + path + ", Search: " + search)
		fmt.Println(strings.EqualFold(path, search))
		if strings.EqualFold(path, search) || strings.Contains(strings.ToLower(path), fmt.Sprintf("%c%s", os.PathSeparator, search)) {
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

func checkForSubQueries(tableRow string) string {
	subQueries := linkMatcher.FindAllStringSubmatch(tableRow, -1)
	for _, subQueryLink := range subQueries {
		subQuery := standardizeSearch(subQueryLink[1])
		subQueryFile, err := findTable(subQuery, ".")
		if err != nil {
			fmt.Printf("Couldn't find table for sub query: %s, Error: %v\n", subQuery, err)
			return tableRow
		}
		subQueryFileReader, err := os.Open(subQueryFile)
		if err != nil {
			fmt.Printf("Couldn't read file for sub query: %s, Error: %v\n", subQuery, err)
			return tableRow
		}
		subQueryValues, err := rollabletable.ParseRollableTable(*bufio.NewScanner(subQueryFileReader))
		if err != nil {
			fmt.Printf("Couldn't parse values for sub query: %s, Error: %s\n", subQuery, err)
			return tableRow
		}
		subQueryResult := subQueryValues.Roll()
		tableRow = strings.Replace(tableRow, subQueryLink[0], subQueryResult, 1)
	}
	return tableRow
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

func contains[T comparable](ts []T, t T) bool {
	for _, v := range ts {
		if v == t {
			return true
		}
	}
	return false
}
