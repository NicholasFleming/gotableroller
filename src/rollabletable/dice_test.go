package rollabletable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDice_Roll(t *testing.T) {
	assert.True(t, Dice{Count: 1, Sides: 1}.Roll() == 1)
	assert.True(t, Dice{Count: 5, Sides: 1}.Roll() == 5)
	assert.Greater(t, Dice{Count: 2, Sides: 2}.Roll(), 1)
	assert.Less(t, Dice{Count: 2, Sides: 2}.Roll(), 5)
}

func Test_diceFromString(t *testing.T) {
	s := (" 2d6 ")
	die, ok := getDiceFromString(s)
	assert.True(t, ok)
	assert.Equal(t, 2, die.Count)
	assert.Equal(t, 6, die.Sides)
}

func Test_diceFromString_MDTableHeader(t *testing.T) {
	s := ("| 2d6 | result |")
	die, ok := getDiceFromString(s)
	assert.True(t, ok)
	assert.Equal(t, 2, die.Count)
	assert.Equal(t, 6, die.Sides)
}
