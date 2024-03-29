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

	"IPutOatsInGoats/gotableroller/src"
	"IPutOatsInGoats/gotableroller/src/rollabletable"
)

var (
	rowRangePattern  = regexp.MustCompile(`(\d+).(\d+)`)   // matches roll ranges like '5-12' and captures the numbers as groups
	markdownListItem = regexp.MustCompile(`^(\d+\. |\* )`) // identifies a line as a markdown list item, ie. '1. ' or '* '

	// Matches markdown links like '[Link Label](path/to/table)' with a group for 'path/to/table' or
	// Internal links like '[[path/to/table]]' with a group for 'path/to/table'
	// Group 1: Either '[foo](' or '[['; Group 2: The path to the table; Group 3: Either ')' or '|foo]]' or ']]'
	linkMatcher = regexp.MustCompile(`(\[.+?\]\(|\[\[)(.+?)(\)|\|.+?\]\]|\]\])`)
	usageText   = "Usage: gotableroller {TableName}\nTableName: the name of the markdown file containing the table. " +
		"This file must exist the same directory or a subdirectory of gotableroller. TableName may/maynot contain" +
		"the '.md' extension. It may contain path components as while. Examples: 'Weapons', 'weapons', 'weapons.md', " +
		"'Items/Weapons.md'"
)

// TODO
// terminal coloring doesnt work on windows

func main() {
	args := os.Args

	query, err := parseArgs(args)
	checkError(err, "Bad command argument")

	rollTables := createRollableTables(query)

	rand.Seed(time.Now().UnixNano())

	var results []string
	for _, table := range rollTables {
		result := rollOnTable(table)
		results = append(results, src.Colorize(src.Green, table.Name+": ")+result)
	}

	for _, result := range results {
		fmt.Println(result)
	}
}

func rollOnTable(rollTable rollabletable.RollableTable) string {
	result := rollTable.Roll()
	for len(linkMatcher.FindStringSubmatch(result)) != 0 {
		link := getLinkFromResult(result)
		subTable := createRollableTables(link.pathToTable)
		subResult := subTable[0].Roll()
		result = strings.Replace(result, link.originalLink, subResult, 1)
	}
	return result
}

type TableLink struct {
	originalLink string
	pathToTable  string
}

func getLinkFromResult(result string) TableLink {
	query := linkMatcher.FindStringSubmatch(result)
	return TableLink{
		originalLink: query[0],
		pathToTable:  query[2],
	}
}

func createRollableTables(query string) (rollTables []rollabletable.RollableTable) {
	query = standardizeSearch(query)
	paths, err := findTable(query, ".")
	checkError(err, "Error finding file")

	for _, path := range paths {
		table, err := rollableTableFromPath(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		rollTables = append(rollTables, table)
	}
	return rollTables
}

func rollableTableFromPath(path string) (rollabletable.RollableTable, error) {

	file, err := os.Open(path)
	checkError(err, "Error reading file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rollTable, err := rollabletable.ParseRollableTable(*scanner, path)
	if err != nil {
		return rollabletable.RollableTable{}, fmt.Errorf("Error parsing table: %s, %v", path, err)
	}
	return rollTable, nil
}

func parseArgs(args []string) (query string, err error) {

	if len(args) < 2 {
		return "", fmt.Errorf("Please provide a table name")
	}

	// TODO use flags
	if contains([]string{"-h", "--h", "-help", "--help", "\\h", "\\help"}, args[1]) {
		fmt.Println(usageText)
		os.Exit(0)
	}

	if contains([]string{"-ls", "--ls", "-list", "--list", "\\ls", "\\list"}, args[1]) {
		if len(args) > 2 {
			query = args[2]
		}
		output := printDirectoryOutput(".", 0, query)
		fmt.Println(output)
		os.Exit(0)
	}

	return args[1], nil
}

func printDirectoryOutput(dir string, depth int, query string) string {
	files, err := os.ReadDir(dir)
	checkError(err, "Error reading directory")
	var directories []os.DirEntry
	buffer := strings.Builder{}
	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			directories = append(directories, file)
		} else if strings.HasSuffix(file.Name(), ".md") {
			if strings.Contains(strings.ToLower(file.Name()), query) {
				buffer.WriteString(src.Colorize(src.Yellow, strings.Repeat("-", depth)+file.Name()+"\n"))
			}
		}
	}

	for _, subDir := range directories {
		buffer.WriteString(src.Colorize(src.Purple, strings.Repeat("-", depth)+dir+string(filepath.Separator)+subDir.Name()+"\n"))
		buffer.WriteString(printDirectoryOutput(dir+string(filepath.Separator)+subDir.Name(), depth+1, query))
	}

	return buffer.String()
}

func standardizeSearch(search string) string {
	search = strings.TrimPrefix(search, "./")
	search = strings.TrimPrefix(search, ".\\")
	search = filepath.FromSlash(search)
	search = strings.ToLower(search)

	return search
}

func findTable(search string, dir string) (paths []string, err error) {
	if search == "" {
		return []string{}, fmt.Errorf("Please provide a table name. Search: %s, Directory: %s", search, dir)
	}

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		checkError(err, "Error while walking directory")
		if d.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(path), strings.ToLower(search)) {
			paths = append(paths, path)
		}
		return nil
	})
	if len(paths) == 0 {
		return []string{}, fmt.Errorf("Table not found: %s", search)
	}
	return paths, err
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}

func contains[T comparable](ts []T, t T) bool {
	for _, v := range ts {
		if v == t {
			return true
		}
	}
	return false
}
