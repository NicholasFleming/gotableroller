package rollabletable

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var foo string = "foo"

func TestMain(m *testing.M) {
	for i := 1; i < 100; i++ {
		foo = "Run number " + strconv.Itoa(i)
		m.Run()

	}
}

func Test_isRollableMDList(t *testing.T) {
	fmt.Println("Foo: " + foo)
	assert.True(t, isRollableMDList("* foo"))
	assert.True(t, isRollableMDList("1. foo"))
	assert.False(t, isRollableMDList("foo"))
	assert.False(t, isRollableMDList("1 foo"))
	assert.False(t, isRollableMDList(""))
	assert.False(t, isRollableMDList("| foo | bar |"))
}

func Test_parseRollableMDList_unordered(t *testing.T) {
	s := []string{"* foo", "* bar", "* baz"}
	list := parseMDList(s)
	assert.Equal(t, MDList{"foo", "bar", "baz"}, list)
}

func Test_parseRollableMDList_ordered(t *testing.T) {
	s := []string{"1. foo", "2. bar", "3. baz"}
	list := parseMDList(s)
	assert.Equal(t, MDList{"foo", "bar", "baz"}, list)
}

func Test_parseRollableMDList_badFormat(t *testing.T) {
	s := []string{"- foo", "2 bar", "baz"}
	list := parseMDList(s)
	var emptyList MDList
	assert.Equal(t, emptyList, list)
}

func Test_isRollableMDTable(t *testing.T) {
	assert.False(t, isRollableMDTable("* foo"))
	assert.False(t, isRollableMDTable("1. foo"))
	assert.False(t, isRollableMDTable("foo"))
	assert.False(t, isRollableMDTable("1 foo"))
	assert.False(t, isRollableMDTable(""))
	assert.True(t, isRollableMDTable("| 1-2 | bar |"))
	assert.True(t, isRollableMDTable("|3-5|bar|"))
	assert.True(t, isRollableMDTable("|---|---|"))
	assert.False(t, isRollableMDTable("| foo | bar | baz |"))
}

func Test_parseMDTable(t *testing.T) {
	s := []string{"| foo | bar |", "|---|---|", "| 1-3 | A |", "| 4-6 | B |"}
	list := parseMDTable(s)
	assert.Equal(t, MDTable{{" foo ", " bar "}, {"---", "---"}, {" 1-3 ", " A "}, {" 4-6 ", " B "}}, list)
}

func Test_parseMDTableRow(t *testing.T) {
	row := [][]string{{" 7-20 ", " foo "}, {" 10 ", " bar "}, {" badinput ", " baz "}}
	min, max, value, ok := parseMDTableRow(row[0])
	assert.True(t, ok)
	assert.Equal(t, 7, min)
	assert.Equal(t, 20, max)
	assert.Equal(t, " foo ", value)

	min, max, value, ok = parseMDTableRow(row[1])
	assert.True(t, ok)
	assert.Equal(t, 10, min)
	assert.Equal(t, 10, max)
	assert.Equal(t, " bar ", value)

	min, max, value, ok = parseMDTableRow(row[2])
	assert.False(t, ok)
}

func Test_fromMDTable(t *testing.T) {
	table := MDTable{{" 1-3 ", " A "}, {" 4-6 ", " B "}}
	rollableTable, err := fromMDTable(table)
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{1: " A ", 2: "1", 3: "1", 4: " B ", 5: "4", 6: "4"}, rollableTable.table)
	assert.Equal(t, 6, rollableTable.max)
}

func Test_fromMDTable_badFormat(t *testing.T) {
	table := MDTable{{" A ", " B "}, {" bar ", " baz "}}
	_, err := fromMDTable(table)
	assert.Error(t, err)
}

func Test_fromMDList(t *testing.T) {
	list := MDList{"foo", "bar", "baz"}
	rollableTable := fromMDList(list)
	assert.Equal(t, map[int]string{1: "foo", 2: "bar", 3: "baz"}, rollableTable.table)
	assert.Equal(t, 3, rollableTable.max)
}

func Test_ParseRollableTable_table_withRanges(t *testing.T) {
	table, err := ParseRollableTable(*bufio.NewScanner(strings.NewReader("| foo | bar |\n|---|---|\n| 1-3 | A |\n| 4-6 | B |")))
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{1: " A ", 2: "1", 3: "1", 4: " B ", 5: "4", 6: "4"}, table.table)
	assert.Equal(t, 6, table.max)
}

func Test_ParseRollableTable_table_withoutRanges(t *testing.T) {
	table, err := ParseRollableTable(*bufio.NewScanner(strings.NewReader("| foo | bar |\n|---|---|\n| 1 | A |\n| 2 | B |\n| 3 | C |")))
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{1: " A ", 2: " B ", 3: " C "}, table.table)
	assert.Equal(t, 3, table.max)
}

func Test_ParseRollableTable_complexTable(t *testing.T) {
	table, err := ParseRollableTable(*bufio.NewScanner(strings.NewReader("| 2d6 | options |\n|---|---|\n| 2 | foo |\n| 3-7 | bar |\n| 8 | baz |\n| 9-12 | bing |")))
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{2: " foo ", 3: " bar ", 4: "3", 5: "3", 6: "3", 7: "3", 8: " baz ", 9: " bing ", 10: "9", 11: "9", 12: "9"}, table.table)
	assert.Equal(t, 12, table.max)
	assert.Equal(t, 2, table.dice.count)
	assert.Equal(t, 6, table.dice.sides)
}

func Test_parseDiceFromMDTable(t *testing.T) {
	table := MDTable{{" 2d20 ", " result "}, {"---", "---"}, {" 1-3 ", " A "}, {" 4-6 ", " B "}, {" 7-20 ", " C "}}
	rollableTable, err := fromMDTable(table)
	assert.Nil(t, err)
	assert.NotNil(t, rollableTable.dice)
	assert.Equal(t, 2, rollableTable.dice.count)
	assert.Equal(t, 20, rollableTable.dice.sides)
	assert.Equal(t, 20, rollableTable.max)
}

func Test_Roll(t *testing.T) {
	table := RollableTable{map[int]string{1: "foo", 2: "bar", 3: "baz"}, 3, Dice{1, 3, AdditionInterpreter{}}}
	match, err := regexp.Match(`foo|bar|baz`, []byte(table.Roll()))
	assert.NoError(t, err)
	assert.True(t, match)
}

func Test_Roll_WithDice1d3(t *testing.T) {
	table := RollableTable{map[int]string{1: "foo", 2: "bar", 3: "baz", 4: "bing", 5: "bong"}, 5, Dice{
		count:           1,
		sides:           3,
		DiceInterpreter: AdditionInterpreter{},
	}}
	match, err := regexp.Match(`foo|bar|baz`, []byte(table.Roll()))
	assert.NoError(t, err)
	assert.True(t, match)
}

func Test_Roll_WithDice2d2(t *testing.T) {
	table := RollableTable{map[int]string{1: "foo", 2: "bar", 3: "baz", 4: "bing", 5: "bong"}, 3, Dice{
		count:           2,
		sides:           2,
		DiceInterpreter: AdditionInterpreter{},
	}}
	match, err := regexp.Match(`bar|baz|bing`, []byte(table.Roll()))
	assert.NoError(t, err)
	assert.True(t, match)
}
