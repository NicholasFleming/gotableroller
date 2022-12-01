package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	caser = cases.Title(language.English)
)

func main() {
	file, err := os.Open("DreadTable.md")

	if err != nil {
		fmt.Println("Error Opening DreadTable: " + err.Error())
		os.Exit(-1)
	}

	scanner := bufio.NewScanner(file)

	h1 := ""
	h2 := ""
	h3 := ""
	var output []string

	for scanner.Scan() {

		line := scanner.Text()

		switch {
		case strings.EqualFold(line, ""):
			continue
		case strings.HasPrefix(line, "###"):
			fmt.Println(line)
			if len(output) > 0 {
				writeTable(h1, h2, h3, output)
			}
			output = []string{}
			h3 = strings.TrimSpace(strings.TrimPrefix(line, "###"))
		case strings.HasPrefix(line, "##"):
			fmt.Println(line)
			if len(output) > 0 {
				writeTable(h1, h2, h3, output)
			}
			output = []string{}
			h3 = ""
			h2 = strings.TrimSpace(strings.TrimPrefix(line, "##"))
		case strings.HasPrefix(line, "#"):
			fmt.Println(line)
			if len(output) > 0 {
				writeTable(h1, h2, h3, output)
			}
			output = []string{}
			h2, h3 = "", ""
			h1 = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			err := os.Mkdir(sanitizeFileName("./"+h1), 0777)
			if err != nil {
				fmt.Println("Error making directory: " + err.Error())
				os.Exit(-1)
			}
		default:
			output = append(output, line)
		}
	}

	writeTable(h1, h2, h3, output)

}

func writeTable(h1 string, h2 string, h3 string, output []string) {
	if len(output) > 0 {
		name := h1 + "/" + h1
		if h2 != "" {
			name = name + " - " + h2
		}
		if h3 != "" {
			name = name + " - " + h3
		}
		err := os.WriteFile(sanitizeFileName(name), []byte(strings.Join(output, "\n")), 0777)
		if err != nil {
			fmt.Println("Error writing file: " + err.Error())
			os.Exit(-1)
		}
	}
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, "?", "")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "&", "")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, "  ", " ")

	return path.Clean(caser.String(name))
}
