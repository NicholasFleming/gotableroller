package rollabletable

import (
	"math/rand"
	"strconv"
)

type Dice struct {
	Count int
	Sides int
}

func (d Dice) Roll() int {
	sum := 0
	for i := 0; i < d.Count; i++ {
		sum += rand.Intn(d.Sides) + 1
	}
	return sum
}

func getDiceFromString(s string) (Dice, bool) {
	matches := diePattern.FindStringSubmatch(s)
	if len(matches) == 0 {
		return Dice{}, false
	}
	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return Dice{}, false
	}
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return Dice{}, false
	}
	return Dice{Count: count, Sides: sides}, true
}
