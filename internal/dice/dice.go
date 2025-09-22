package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// DiceRoll represents the parsed components of a dice roll notation.
type DiceRoll struct {
	NumDice    int
	DieType    int
	Modifier   int
	Notation   string
}

// ParseDiceNotation parses a string like "2d6+5" into a DiceRoll struct.
func ParseDiceNotation(notation string) (*DiceRoll, error) {
	re := regexp.MustCompile(`^(\d+)d(\d+)(?:([+-])(\d+))?$`)
	matches := re.FindStringSubmatch(notation)

	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid dice notation: %s", notation)
	}

	numDice, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid number of dice: %w", err)
	}
	if numDice <= 0 {
		return nil, fmt.Errorf("number of dice must be positive")
	}

	dieType, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid die type: %w", err)
	}
	if dieType <= 0 {
		return nil, fmt.Errorf("die type must be positive")
	}

	modifier := 0
	if len(matches) > 3 && matches[3] != "" {
		modValue, err := strconv.Atoi(matches[4])
		if err != nil {
			return nil, fmt.Errorf("invalid modifier value: %w", err)
		}
		if matches[3] == "-" {
			modifier = -modValue
		} else {
			modifier = modValue
		}
	}

	return &DiceRoll{
		NumDice:    numDice,
		DieType:    dieType,
		Modifier:   modifier,
		Notation:   notation,
	}, nil
}

// Roll performs the dice roll based on the DiceRoll struct.
func (dr *DiceRoll) Roll() (int, []int) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	rolls := make([]int, dr.NumDice)
	total := 0

	for i := 0; i < dr.NumDice; i++ {
		roll := rand.Intn(dr.DieType) + 1 // rand.Intn(n) returns [0, n-1], so add 1 for [1, n]
		rolls[i] = roll
		total += roll
	}

	total += dr.Modifier
	return total, rolls
}
