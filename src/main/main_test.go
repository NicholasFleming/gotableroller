package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
