package rollabletable

import (
	"bufio"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isRollableMDList(t *testing.T) {
	assert.True(t, isRollableMDList("* foo"))
	assert.True(t, isRollableMDList("1. foo"))
	assert.False(t, isRollableMDList("foo"))
	assert.False(t, isRollableMDList("1 foo"))
	assert.False(t, isRollableMDList(""))
	assert.False(t, isRollableMDList("| foo | bar |"))
}

func Test_parseRollableMDList_unordered(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("* foo\n* bar\n* baz"))
	list := parseMDList(*scanner)
	assert.Equal(t, MDList{"foo", "bar", "baz"}, list)
}

func Test_parseRollableMDList_ordered(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("1. foo\n2. bar\n3. baz"))
	list := parseMDList(*scanner)
	assert.Equal(t, MDList{"foo", "bar", "baz"}, list)
}

func Test_parseRollableMDList_badFormat(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("- foo\n2 bar\nbaz"))
	list := parseMDList(*scanner)
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
	assert.False(t, isRollableMDTable("|---|---|"))
	assert.False(t, isRollableMDTable("| foo | bar |"))
}

func Test_parseRollableMDTable(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader("| foo | bar |\n|---|---|\n| 1-3 | A |\n| 4-6 | B |"))
	list := parseMDTable(*scanner)
	assert.Equal(t, MDTable{{" 1-3 ", " A "}, {" 4-6 ", " B "}}, list)
}

func Test_parseMDTableRow(t *testing.T) {
	row := []string{" 7-20 ", " foo "}
	min, max, value, err := parseMDTableRow(row)
	assert.NoError(t, err)
	assert.Equal(t, 7, min)
	assert.Equal(t, 20, max)
	assert.Equal(t, " foo ", value)
}

func Test_fromMDTable(t *testing.T) {
	table := MDTable{{" 1-3 ", " A "}, {" 4-6 ", " B "}}
	rollableTable, err := fromMDTable(table)
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{1: " A ", 2: "1", 3: "1", 4: " B ", 5: "4", 6: "4"}, rollableTable.table)
	assert.Equal(t, 6, rollableTable.max)
}

func Test_fromMDTable_badFormat(t *testing.T) {
	table := MDTable{{" 1-3 ", " A "}, {" bar ", " baz "}}
	_, err := fromMDTable(table)
	assert.Error(t, err)
}

func Test_fromMDList(t *testing.T) {
	list := MDList{"foo", "bar", "baz"}
	rollableTable := fromMDList(list)
	assert.Equal(t, map[int]string{1: "foo", 2: "bar", 3: "baz"}, rollableTable.table)
	assert.Equal(t, 3, rollableTable.max)
}

func Test_ParseRollableTable_table(t *testing.T) {
	table, err := ParseRollableTable(*bufio.NewScanner(strings.NewReader("| foo | bar |\n|---|---|\n| 1-3 | A |\n| 4-6 | B |")))
	assert.NoError(t, err)
	assert.Equal(t, map[int]string{1: " A ", 2: "1", 3: "1", 4: " B ", 5: "4", 6: "4"}, table.table)
	assert.Equal(t, 6, table.max)
}

func Test_Roll(t *testing.T) {
	table := RollableTable{map[int]string{1: "foo", 2: "bar", 3: "baz"}, 3}
	match, err := regexp.Match(`foo|bar|baz`, []byte(table.Roll()))
	assert.NoError(t, err)
	assert.True(t, match)

}
