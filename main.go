package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//
// TODO
// rename project to markdownTableRoller
// handle number ranges eg 1-4, 5-6 for varying probabilities, probably with table syntax
// handle tables with same name in multiple directories. print paths and one random selection?
// update README
//

func main() {
	args := os.Args

	query, err := parseArgs(args)
	checkError(err, "Bad command argument")

	query = standardizeSearch(query)

	path, err := findTable(query, ".")
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

func parseTableValues(tableFile *os.File) ([]string, error) {
	scanner := bufio.NewScanner(tableFile)
	var tableValues []string
	for scanner.Scan() {
		markdownListItem, err := regexp.Compile(`^(\d+\. |\* )`)
		checkError(err, "Error parsing table values")
		isValue := markdownListItem.Match([]byte(scanner.Text()))
		if isValue {
			noPrefix := markdownListItem.ReplaceAllString(scanner.Text(), "")
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
	// Ex. Matches markdown links like [Link Label](path/to/table) With subgroups: [1] = path/to/table
	linkMatcher := regexp.MustCompile(`\[.+?\]\((.+?)\)`)

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
		subQueryValues, err := parseTableValues(subQueryFileReader)
		if err != nil {
			fmt.Printf("Couldn't parse values for sub query: %s, Error: %s\n", subQuery, err)
			return tableRow
		}
		subQueryResult := rollTheDice(subQueryValues)
		tableRow = strings.Replace(tableRow, subQueryLink[0], subQueryResult, 1)
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

func contains[T comparable](ts []T, t T) bool {
	for _, v := range ts {
		if v == t {
			return true
		}
	}
	return false
}

func readMarkdownTable(scanner *bufio.Scanner) (table [][]string, err error) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "|") && !strings.HasPrefix(line, "|--") {
			table = append(table, strings.Split(line, "|")[1:3])
		}
	}
	return table, nil
}

func rollOnTable(table [][]string) (result string, err error) {
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(len(table)) // this is wrong, it should go from 1 to max of the last range
	rowRangePattern := regexp.MustCompile(`(\d+).(\d+)`)
	for _, row := range table {
		rowRange := rowRangePattern.FindAllStringSubmatch(row[0], -1)
		if len(rowRange) > 0 {
			min, err := strconv.Atoi(rowRange[0][1])
			if err != nil {
				return "", fmt.Errorf("Error parsing table value min roll range, row: %s, Error: %w", row[0], err)
			}
			max, err := strconv.Atoi(rowRange[0][2])
			if err != nil {
				return "", fmt.Errorf("Error parsing table value max roll range, row: %s, Error: %w", row[0], err)
			}
			fmt.Printf("Min: %d, Max: %d, Roll: %d", min, max, roll)
			if roll >= min && roll <= max {
				result = row[1]
				break
			}
		}
	}
	return result, nil
}
