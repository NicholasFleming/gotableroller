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
)

/* TODO
Add Maze rats md to use for examples and dev since its creative commons
rename project to markdownTableRoller
handle number ranges eg 1-4, 5-6 for varying probabilities, probably with table syntax
handle tables with same name in multiple directories
work on windows
if [subtable] is not found just print [subtable]
add more maze rats tables
*/

func main() {
	args := os.Args

	query, err := parseArgs(args)
	checkError(err, "Bad command argument")

	path, err := findTable(query, "./")
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

	query = strings.TrimPrefix(args[1], "."+string(os.PathSeparator))
	query = strings.TrimSuffix(query, ".md")

	return query, nil
}

func findTable(search string, dir string) (tablePath string, err error) {
	if search == "" {
		return "", fmt.Errorf("Please provide a table name. Search: %s, Directory: %s", search, dir)
	}
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if tablePath != "" {
			return nil
		}
		if strings.Contains(strings.ToLower(path), strings.ToLower(search+".md")) {
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
		matcher, err := regexp.Compile(`^(\d+\. |\* )`)
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
	matcher, err := regexp.Compile(`\[.+?\]`)
	if err != nil {
		fmt.Printf("Couldn't understand subwuery in: %s\n", tableRow)
		return tableRow
	}
	hasSubQueries := matcher.Match([]byte(tableRow))

	if hasSubQueries {
		subQueries := matcher.FindAllString(tableRow, -1)
		for _, subQuery := range subQueries {
			subQuery = strings.TrimPrefix(subQuery, "[")
			subQuery = strings.TrimSuffix(subQuery, "]")
			subQueryFile, err := findTable(subQuery, "./")
			if err != nil {
				fmt.Printf("Couldn't find table for sub query: %s\n", subQuery)
				return tableRow
			}
			subQueryFileReader, err := os.Open(subQueryFile)
			if err != nil {
				fmt.Printf("Couldn't read file for sub query: %s\n", subQuery)
				return tableRow
			}
			subQueryValues, err := parseTableValues(subQueryFileReader)
			if err != nil {
				fmt.Printf("Couldn't parse values for sub query: %s\n", subQuery)
				return tableRow
			}
			subQueryResult := rollTheDice(subQueryValues)
			// TODO handle multiple instances of the subquery
			foo := strings.Replace(tableRow, fmt.Sprintf("[%s]", subQuery), subQueryResult, 1)
			tableRow = foo
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
