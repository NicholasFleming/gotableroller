package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseArgs(t *testing.T) {
	args := []string{"foo", "TestTable"}
	query, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, "TestTable", query)
}

func Test_parseArgs_withExtension(t *testing.T) {
	args := []string{"foo", "TestTable.md"}
	query, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, "TestTable", query)
}

func Test_parseArgs_dotSlash(t *testing.T) {
	args := []string{"foo", dotSlash + "TestTable"}
	query, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, "TestTable", query)
}

func Test_parseArgs_dotSlash_withExtension(t *testing.T) {
	args := []string{"foo", dotSlash + "TestTable.md"}
	query, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, "TestTable", query)
}

func Test_parseArgs_pathToFile(t *testing.T) {
	args := []string{"foo", filepath.FromSlash("Test/TestTable.md")}
	query, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/TestTable"), query)
}

func Test_parseArgs_pathToFileInSubDir(t *testing.T) {
	args := []string{"foo", filepath.FromSlash("Test/testdir/SubTestTable.md")}
	query, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/testdir/SubTestTable"), query)
}

func Test_findFiles(t *testing.T) {
	path, err := findTable("TestTable", dotSlash+"Test")
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/TestTable.md"), path)
}

func Test_findFiles_subDir(t *testing.T) {
	path, err := findTable("SubTestTable", dotSlash+"Test")
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/testdir/SubTestTable.md"), path)
}

func Test_findFiles_specifyPath(t *testing.T) {
	path, err := findTable(filepath.FromSlash("testdir/SubTestTable"), dotSlash+"Test")
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/testdir/SubTestTable.md"), path)

	path, err = findTable(filepath.FromSlash("Test/testdir/SubTestTable"), dotSlash+"Test")
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/testdir/SubTestTable.md"), path)
}

func Test_findFiles_caseInsensitive(t *testing.T) {
	path, err := findTable("testtable", dotSlash+"Test")
	assert.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("Test/TestTable.md"), path)
}

func Test_findFiles_BadName(t *testing.T) {
	path, err := findTable("I_Dont_Exist", dotSlash+"Test")
	assert.Error(t, err)
	assert.Empty(t, path)
}

func Test_findFiles_EmptyName(t *testing.T) {
	path, err := findTable("", dotSlash+"Test")
	assert.Error(t, err)
	assert.Empty(t, path)
}

func Test_parseTableValues_bulleted(t *testing.T) {
	file, _ := os.Open(filepath.FromSlash("Test/TestTable.md"))
	defer file.Close()
	actual, err := parseTableValues(file)
	assert.NoError(t, err)
	assert.Contains(t, actual, "Option1")
	assert.Contains(t, actual, "Option2")
	assert.Contains(t, actual, "Option3")
}

func Test_parseTableValues_numbered(t *testing.T) {
	file, _ := os.Open(filepath.FromSlash("Test/testdir/SubTestTable.md"))
	defer file.Close()
	actual, err := parseTableValues(file)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Sub Option1", "Sub Option2", "Sub Option3"}, actual)
}

func Test_linkedTables(t *testing.T) {
	file, _ := os.Open(filepath.FromSlash("Test/TestTable.md"))
	file2, _ := os.Open(filepath.FromSlash("Test/testdir/SubTestTable.md"))
	defer file.Close()
	defer file2.Close()
	tableValues, err := parseTableValues(file)
	assert.NoError(t, err)
	subTableValues, err := parseTableValues(file2)
	assert.NoError(t, err)
	assert.NotContains(t, tableValues, "Option with [SubTestTable](testdir/SubTestTable)")
	found := false
	for _, value := range tableValues {
		for _, subValue := range subTableValues {
			if strings.Contains(value, subValue) {
				found = true
			}
		}
	}
	assert.True(t, found)
}
