package rollabletable

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var (
	rowRangePattern  = regexp.MustCompile(`(\d+)-(\d+)`)   // matches roll ranges like '5-12' and captures the numbers as groups
	markdownListItem = regexp.MustCompile(`^(\d+\. |\* )`) // identifies a line as a markdown list item, ie. '1. ' or '* '
)

type RollableTable struct {
	table map[int]string
	max   int
}

func (rt RollableTable) Roll() string {
	roll := rand.Intn(rt.max) + 1
	i, err := strconv.Atoi(rt.table[roll])
	if err != nil {
		return rt.table[roll]
	}
	return rt.table[i]
}

func (rt RollableTable) AsMDTable() string {
	var table bytes.Buffer
	for k, v := range rt.table {
		table.WriteString(fmt.Sprintf("| %d | %s |\n", k, v))
	}
	return table.String()
}

func ParseRollableTable(scanner bufio.Scanner) (RollableTable, error) {
	for i := 0; i < 5; i++ { // Only check first couple lines before moving on
		scanner.Scan()
		switch {
		case isRollableMDList(scanner.Text()):
			return fromMDList(parseMDList(scanner)), nil
		case isRollableMDTable(scanner.Text()):
			return fromMDTable(parseMDTable(scanner))
		}
	}
	return RollableTable{}, fmt.Errorf("Not a Rollable Table")
}

func fromMDList(list MDList) RollableTable {
	var rollableTable RollableTable
	rollableTable.table = make(map[int]string)
	rollableTable.max = len(list)
	for i, line := range list {
		rollableTable.table[i+1] = line
	}

	return rollableTable
}

func fromMDTable(table MDTable) (RollableTable, error) {
	var rollableTable RollableTable
	rollableTable.table = make(map[int]string)
	for _, row := range table {
		minRange, maxRange, value, err := parseMDTableRow(row)
		if err != nil {
			return RollableTable{}, fmt.Errorf("Error in Table Format: %w", err)
		}
		rollableTable.table[minRange] = value
		for i := minRange + 1; i <= maxRange; i++ {
			rollableTable.table[i] = strconv.Itoa(minRange)
		}
		if rollableTable.max < maxRange {
			rollableTable.max = maxRange
		}
	}
	return rollableTable, nil
}

func parseMDTableRow(row []string) (min int, max int, value string, err error) {
	rowRange := rowRangePattern.FindAllStringSubmatch(row[0], -1)
	if len(rowRange) != 1 || len(rowRange[0]) != 3 {
		return 0, 0, "", fmt.Errorf("Bad row range value: %v", row[0])
	}
	min, err = strconv.Atoi(rowRange[0][1])
	if err != nil {
		return 0, 0, "", fmt.Errorf("Error parsing table value min roll range, row: %s, Error: %w", row[0], err)
	}
	max, err = strconv.Atoi(rowRange[0][2])
	if err != nil {
		return 0, 0, "", fmt.Errorf("Error parsing table value max roll range, row: %s, Error: %w", row[0], err)
	}
	return min, max, row[1], nil
}

type MDTable [][]string

func isRollableMDTable(s string) bool {
	if strings.HasPrefix(s, "|") {
		columns := strings.Count(s, "|") - strings.Count(s, `\|`)
		if columns == 3 && rowRangePattern.MatchString(strings.Split(s, "|")[1]) {
			return true
		}
	}
	return false
}

func parseMDTable(scanner bufio.Scanner) MDTable {
	var mdTable MDTable
	line := scanner.Text()
	if strings.HasPrefix(line, "|") && rowRangePattern.MatchString(line) {
		mdTable = append(mdTable, strings.Split(line, "|")[1:3])
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "|") && rowRangePattern.MatchString(line) {
			mdTable = append(mdTable, strings.Split(line, "|")[1:3])
		}
	}
	return mdTable
}

type MDList []string

func isRollableMDList(s string) bool {
	return markdownListItem.MatchString(s)
}

func parseMDList(scanner bufio.Scanner) MDList {
	var mdList MDList
	if markdownListItem.Match([]byte(scanner.Text())) {
		mdList = append(mdList, markdownListItem.ReplaceAllString(scanner.Text(), ""))
	}
	for scanner.Scan() {
		if markdownListItem.Match([]byte(scanner.Text())) {
			mdList = append(mdList, markdownListItem.ReplaceAllString(scanner.Text(), ""))
		}
	}
	return mdList
}
