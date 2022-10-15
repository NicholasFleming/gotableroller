package src

import "fmt"

type color string

const (
	Reset     color = "\033[0m"
	Bold      color = "\033[1m"
	Italic    color = "\033[3m"
	Underline color = "\033[4m"
	Red       color = "\033[31m"
	Green     color = "\033[32m"
	Yellow    color = "\033[33m"
	Blue      color = "\033[34m"
	Purple    color = "\033[35m"
	Cyan      color = "\033[36m"
	Gray      color = "\033[37m"
	White     color = "\033[97m"
)

func Colorize(c color, s string) string {
	return string(c) + s + string(Reset)
}

func Colorizef(c color, format string, a ...any) string {
	return Colorize(c, fmt.Sprintf(format, a...))
}
