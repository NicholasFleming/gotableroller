package main

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"IPutOatsInGoats/gotableroller/src/rollabletable"
)

func Test_standardizeSearch(t *testing.T) {
	assert.Equal(t, "testtable", standardizeSearch("TestTable"))
	assert.Equal(t, "testtable.md", standardizeSearch("TestTable.md"))
	assert.Equal(t, "testtable", standardizeSearch("testtable"))
	assert.Equal(t, "testtable.md", standardizeSearch("testtable.md"))
	assert.Equal(t, "testtable", standardizeSearch("./testtable"))
	assert.Equal(t, "testtable", standardizeSearch(".\\testtable"))
	assert.Equal(t, filepath.FromSlash("test/testtable.md"), standardizeSearch("Test/TestTable.md"))
	assert.Equal(t, filepath.FromSlash("test/testtable"), standardizeSearch("Test/TestTable"))
}

func Test_parseArgs(t *testing.T) {
	query, err := parseArgs([]string{"foo", "TestTable"})
	assert.NoError(t, err)
	assert.Equal(t, "TestTable", query)
}

func Test_parseArgs_noQuery(t *testing.T) {
	_, err := parseArgs([]string{"foo"})
	assert.Error(t, err)
}

func Test_findFiles(t *testing.T) {
	path, err := findTable("testtable.md", "Test")
	assert.NoError(t, err)
	assert.Equal(t, []string{filepath.FromSlash("Test/TestTable.md"), filepath.FromSlash("Test/testdir/SubTestTable.md")}, path)
}

func Test_findFiles_subDir(t *testing.T) {
	path, err := findTable("subtesttable.md", "Test")
	assert.NoError(t, err)
	assert.Equal(t, []string{filepath.FromSlash("Test/testdir/SubTestTable.md")}, path)
}

func Test_findFiles_specifyPath(t *testing.T) {
	path, err := findTable(filepath.FromSlash("testdir/subtesttable.md"), "Test")
	assert.NoError(t, err)
	assert.Equal(t, []string{filepath.FromSlash("Test/testdir/SubTestTable.md")}, path)
}

func Test_findFiles_BadName(t *testing.T) {
	path, err := findTable("I_Dont_Exist", "Test")
	assert.Error(t, err)
	assert.Empty(t, path)
}

func Test_findFiles_EmptyName(t *testing.T) {
	path, err := findTable("", "Test")
	assert.Error(t, err)
	assert.Empty(t, path)
}

func Test_getLinkFromResult(t *testing.T) {
	link := getLinkFromResult("foo [text label](path/to/file) bar")
	assert.Len(t, link, 2)
	assert.Equal(t, "[text label](path/to/file)", link[0])
	assert.Equal(t, "path/to/file", link[1])
}

func Test_rollOnTable(t *testing.T) {
	table, err := rollabletable.ParseRollableTable(*bufio.NewScanner(strings.NewReader("* foo\n* bar\n* baz\n")), "mdtable")
	assert.NoError(t, err)
	result := rollOnTable(table)
	match, err := regexp.Match("foo|bar|baz", []byte(result))
	assert.NoError(t, err)
	assert.True(t, match)
}

func Test_rollOnTable_with_link(t *testing.T) {
	tables := createRollableTables("TestTable")
	result := rollOnTable(tables[0])
	assert.NotEmpty(t, result)
}

func Test_printDirectoryOutput(t *testing.T) {
	output := printDirectoryOutput(".", 0, "testtable")
	assert.True(t, strings.Contains(output, "./Test"))
	assert.True(t, strings.Contains(output, "-TestTable.md"))
	assert.True(t, strings.Contains(output, "-TestTableTable.md"))
	assert.True(t, strings.Contains(output, "-./Test/testdir"))
	assert.True(t, strings.Contains(output, "--SubTestTable.md"))
}

func Test_rollableTableFromPath(t *testing.T) {
	table, err := rollableTableFromPath(filepath.FromSlash("Test/TestTable.md"))
	assert.NoError(t, err)
	assert.True(t, strings.Contains(table.Name, "TestTable"))
	assert.NotEmpty(t, table.Roll())
}

func Test_createRollableTables(t *testing.T) {
	tables := createRollableTables("TestTable")
	assert.NotEmpty(t, tables)
	assert.True(t, strings.Contains(tables[0].Name, "TestTable"))
	assert.NotEmpty(t, tables[0].Roll())
}

func Test_contains(t *testing.T) {
	assert.True(t, contains([]string{"foo", "bar", "baz"}, "foo"))
	assert.False(t, contains([]string{"foo", "bar", "baz"}, "qux"))
}
