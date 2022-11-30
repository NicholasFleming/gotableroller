package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("DreadTable.md")

	if err != nil {
		fmt.Println(err)
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
			writeTable(h1, h2, h3, output)
			output = []string{}
			h3 = strings.TrimPrefix(line, "###")
		case strings.HasPrefix(line, "##"):
			writeTable(h1, h2, h3, output)
			output = []string{}
			h2 = strings.TrimPrefix(line, "##")
		case strings.HasPrefix(line, "#"):
			err := os.Mkdir("/"+h1, 0777)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			writeTable(h1, h2, h3, output)
			output = []string{}
			h1 = strings.TrimPrefix(line, "#")
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
		os.WriteFile(name, []byte(strings.Join(output, "\n")), 0777)
	}
}
