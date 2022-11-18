package rollabletable

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDice_Roll(t *testing.T) {
	assert.True(t, Dice{count: 1, sides: 1, DiceInterpreter: AdditionInterpreter{}}.Roll() == 1)
	assert.True(t, Dice{count: 5, sides: 1, DiceInterpreter: AdditionInterpreter{}}.Roll() == 5)
	assert.Greater(t, Dice{count: 2, sides: 2, DiceInterpreter: AdditionInterpreter{}}.Roll(), 1)
	assert.Less(t, Dice{count: 2, sides: 2, DiceInterpreter: AdditionInterpreter{}}.Roll(), 5)
}

func Test_rollAllDice(t *testing.T) {
	dice := Dice{
		count: 2,
		sides: 6,
	}
	result := dice.rollAllDice()
	fmt.Printf("result: %v\n", result)
	assert.Equal(t, 2, len(result))
	assert.GreaterOrEqual(t, result[0], 1)
	assert.LessOrEqual(t, result[0], 6)
	assert.GreaterOrEqual(t, result[1], 1)
	assert.LessOrEqual(t, result[1], 6)
}

func Test_diceFromString(t *testing.T) {
	s := (" 2d6 ")
	die, ok := parseDiceFromString(s)
	assert.True(t, ok)
	assert.Equal(t, 2, die.count)
	assert.Equal(t, 6, die.sides)
	assert.IsType(t, AdditionInterpreter{}, die.DiceInterpreter)
}

func Test_diceFromString_MDTableHeader(t *testing.T) {
	s := ("| 2d6 | result |")
	die, ok := parseDiceFromString(s)
	assert.True(t, ok)
	assert.Equal(t, 2, die.count)
	assert.Equal(t, 6, die.sides)
	assert.IsType(t, AdditionInterpreter{}, die.DiceInterpreter)
}

func Test_diceFromString_digitDie(t *testing.T) {
	s := ("| 1d66 | result |")
	die, ok := parseDiceFromString(s)
	assert.True(t, ok)
	assert.Equal(t, 2, die.count)
	assert.Equal(t, 6, die.sides)
	assert.IsType(t, DigitsInterpreter{}, die.DiceInterpreter)
}

func Test_AdditionInterpreter(t *testing.T) {
	var ai AdditionInterpreter
	dieResult := DieResult{1, 2, 3}
	assert.Equal(t, 6, ai.interpret(dieResult))
}

func Test_DigitsInterpreter(t *testing.T) {
	var ai DigitsInterpreter
	dieResult := DieResult{1, 2, 3}
	assert.Equal(t, 123, ai.interpret(dieResult))
}
