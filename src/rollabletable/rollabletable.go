package rollabletable

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	rowRangePattern  = regexp.MustCompile(`(\d+)[-|–](\d+)`)     // matches roll ranges like '5-12' and captures the numbers as groups
	markdownListItem = regexp.MustCompile(`^(\d+\. |\* |- |– )`) // identifies a line as a markdown list item, ie. '1. ' or '* '
)

type RollableTable struct {
	Name  string
	table map[int]string
	max   int
	dice  Dice
}

func (rt RollableTable) Roll() string {
	result := rt.dice.Roll()

	index, err := strconv.Atoi(rt.table[result])
	if err != nil {
		return rt.table[result]
	}
	return rt.table[index]
}

func (rt RollableTable) AsMDTable() string {
	var table bytes.Buffer
	for k, v := range rt.table {
		table.WriteString(fmt.Sprintf("| %d | %s |\n", k, v))
	}
	return table.String()
}

func ParseRollableTable(scanner bufio.Scanner, name string) (RollableTable, error) {
	var doc []string
	for i := 0; i < 5; i++ { // Only check first couple lines before moving on
		scanner.Scan()
		doc = append(doc, scanner.Text())
		switch {
		case isRollableMDList(scanner.Text()):
			for scanner.Scan() {
				doc = append(doc, scanner.Text())
			}
			return fromMDList(parseMDList(doc), name), nil
		case isRollableMDTable(scanner.Text()):
			for scanner.Scan() {
				doc = append(doc, scanner.Text())
			}
			return fromMDTable(parseMDTable(doc), name)
		}
	}
	return RollableTable{}, fmt.Errorf("Not a Rollable Table")
}

func fromMDList(list MDList, name string) RollableTable {
	var rollableTable RollableTable
	rollableTable.Name = name
	rollableTable.table = make(map[int]string)
	rollableTable.max = len(list)
	rollableTable.dice = Dice{
		count:           1,
		sides:           len(list),
		DiceInterpreter: AdditionInterpreter{},
	}
	for i, line := range list {
		rollableTable.table[i+1] = line
	}

	return rollableTable
}

func fromMDTable(table MDTable, name string) (RollableTable, error) {
	var rollableTable RollableTable
	rollableTable.Name = name
	rollableTable.table = make(map[int]string)
	for _, row := range table {
		minRange, maxRange, value, ok := parseMDTableRow(row)
		if ok {
			rollableTable.table[minRange] = value
			for i := minRange + 1; i <= maxRange; i++ {
				rollableTable.table[i] = strconv.Itoa(minRange)
			}
			if rollableTable.max < maxRange {
				rollableTable.max = maxRange
			}
		}
	}
	if rollableTable.max == 0 || len(rollableTable.table) == 0 {
		return rollableTable, fmt.Errorf("Table not parsable as Rollable Table, table max: %d, table length: %d", rollableTable.max, len(rollableTable.table))
	}
	die, dieDefined := parseDiceFromString(table[0][0])
	if !dieDefined {
		die = Dice{
			count:           1,
			sides:           rollableTable.max,
			DiceInterpreter: AdditionInterpreter{},
		}
	}
	rollableTable.dice = die
	return rollableTable, nil
}

func parseMDTableRow(row []string) (min int, max int, value string, ok bool) {
	rowRange := rowRangePattern.FindAllStringSubmatch(row[0], -1)

	if len(rowRange) == 1 && len(rowRange[0]) == 3 {
		min, err := strconv.Atoi(rowRange[0][1])
		if err != nil {
			fmt.Printf("Error parsing min range, %v, %s\n", row, err)
			return 0, 0, "", false
		}
		max, err = strconv.Atoi(rowRange[0][2])
		if err != nil {
			fmt.Printf("Error parsing max range, %v, %s\n", row, err)
			return 0, 0, "", false
		}
		return min, max, row[1], true
	}
	num, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err == nil {
		return num, num, row[1], true
	}
	return 0, 0, "", false
}

type MDTable [][]string

func isRollableMDTable(s string) bool {
	if strings.HasPrefix(s, "|") {
		columns := strings.Count(s, "|") - strings.Count(s, `\|`)
		if columns == 3 {
			return true
		}
	}
	return false
}

func parseMDTable(contents []string) MDTable {
	var mdTable MDTable
	for _, line := range contents {
		if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
			mdTable = append(mdTable, strings.Split(line, "|")[1:3])
		}
	}
	return mdTable
}

type MDList []string

func isRollableMDList(s string) bool {
	return markdownListItem.MatchString(s)
}

func parseMDList(contents []string) MDList {
	var mdList MDList
	entry := ""
	for _, line := range contents {
		if markdownListItem.Match([]byte(line)) {
			if entry != "" {
				mdList = append(mdList, entry)
				entry = ""
			}
			entry += markdownListItem.ReplaceAllString(line, "")
		} else {
			if line != "" {
				entry += "\n" + markdownListItem.ReplaceAllString(line, "")
			}
		}
	}
	if entry != "" {
		mdList = append(mdList, entry)
	}
	return mdList
}
