package rollabletable

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var (
	standardDiePattern = regexp.MustCompile(`(\d+)d(\d+)`)
	digitsDiePattern   = regexp.MustCompile(`66|88`)
)

type DiceInterpreter interface {
	interpret(result DieResult) int
}

type AdditionInterpreter struct{}

func (ai AdditionInterpreter) interpret(result DieResult) int {
	sum := 0
	for _, v := range result {
		sum += v
	}
	return sum
}

type DigitsInterpreter struct{}

func (di DigitsInterpreter) interpret(result DieResult) int {
	sum := 0
	for _, v := range result {
		sum *= 10
		sum += v
	}
	return sum
}

type Dice struct {
	count int
	sides int
	DiceInterpreter
}

func (d Dice) rollAllDice() DieResult {
	var result []int
	for i := 0; i < d.count; i++ {
		result = append(result, rand.Intn(d.sides)+1)
	}
	return result
}

func (d Dice) Roll() int {
	fmt.Printf("Rolling Dice: %d, %d, %T\n", d.count, d.sides, d.DiceInterpreter)

	return d.DiceInterpreter.interpret(d.rollAllDice())
}

type DieResult []int

func (d DieResult) StandardInterpretation() int {
	sum := 0
	for _, v := range d {
		sum += v
	}
	return sum
}

func (d DieResult) DigitsInterpretation() int {
	sum := 0
	for i, v := range d {
		sum += v * (10 ^ i)
	}
	return sum
}

func parseDiceFromString(s string) (Dice, bool) {
	if digitsDiePattern.MatchString(s) {
		sides, err := strconv.Atoi(strings.Split(strings.Split(s, "d")[1], "")[0])
		if err != nil {
			return Dice{}, false
		}
		return Dice{
			count:           2,
			sides:           sides,
			DiceInterpreter: DigitsInterpreter{},
		}, true
	}

	matches := standardDiePattern.FindStringSubmatch(s)
	if len(matches) == 0 {
		return Dice{}, false
	}
	diceCount, err := strconv.Atoi(matches[1])
	if err != nil {
		return Dice{}, false
	}
	diceSides, err := strconv.Atoi(matches[2])
	if err != nil {
		return Dice{}, false
	}
	return Dice{
		count:           diceCount,
		sides:           diceSides,
		DiceInterpreter: AdditionInterpreter{},
	}, true
}
