package dice

import (
	"testing"
)

func TestParseDiceNotation(t *testing.T) {
	tests := []struct {
		name     string
		notation string
		expectNumDice  int
		expectDieType  int
		expectModifier int
		expectErr      bool
	}{
		{name: "2d6", notation: "2d6", expectNumDice: 2, expectDieType: 6, expectModifier: 0, expectErr: false},
		{name: "1d20+5", notation: "1d20+5", expectNumDice: 1, expectDieType: 20, expectModifier: 5, expectErr: false},
		{name: "3d8-2", notation: "3d8-2", expectNumDice: 3, expectDieType: 8, expectModifier: -2, expectErr: false},
		{name: "1d4", notation: "1d4", expectNumDice: 1, expectDieType: 4, expectModifier: 0, expectErr: false},
		{name: "invalid-format", notation: "2d", expectErr: true},
		{name: "invalid-num-dice", notation: "0d6", expectErr: true},
		{name: "invalid-die-type", notation: "2d0", expectErr: true},
		{name: "invalid-modifier", notation: "1d20+", expectErr: true},
		{name: "no-d", notation: "26", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr, err := ParseDiceNotation(tt.notation)

			if (err != nil) != tt.expectErr {
				t.Errorf("ParseDiceNotation() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr {
				if dr.NumDice != tt.expectNumDice {
					t.Errorf("ParseDiceNotation() NumDice = %v, want %v", dr.NumDice, tt.expectNumDice)
				}
				if dr.DieType != tt.expectDieType {
					t.Errorf("ParseDiceNotation() DieType = %v, want %v", dr.DieType, tt.expectDieType)
				}
				if dr.Modifier != tt.expectModifier {
					t.Errorf("ParseDiceNotation() Modifier = %v, want %v", dr.Modifier, tt.expectModifier)
				}
			}
		})
	}
}

func TestRoll(t *testing.T) {
	tests := []struct {
		name     string
		notation string
		expectMin int
		expectMax int
	}{
		{name: "1d6", notation: "1d6", expectMin: 1, expectMax: 6},
		{name: "2d4+2", notation: "2d4+2", expectMin: 4, expectMax: 10},
		{name: "1d1-1", notation: "1d1-1", expectMin: 0, expectMax: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr, err := ParseDiceNotation(tt.notation)
			if err != nil {
				t.Fatalf("ParseDiceNotation() failed: %v", err)
			}

			// Roll multiple times to ensure randomness and range
			for i := 0; i < 100; i++ {
				total, _ := dr.Roll()
				if total < tt.expectMin || total > tt.expectMax {
					t.Errorf("Roll() for %s returned %d, expected between %d and %d", tt.notation, total, tt.expectMin, tt.expectMax)
				}
			}
		})
	}
}
