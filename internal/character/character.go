package character

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Character represents a D&D 5e character
type Character struct {
	Name        string `json:"name"`
	Species     string `json:"species"`
	Class       string `json:"class"`
	Level       int    `json:"level"`
	Background  string `json:"background"`

	// Ability Scores
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`

	// Derived Stats
	HitPoints    int `json:"hit_points"`
	ArmorClass   int `json:"armor_class"`
	ProficiencyBonus int `json:"proficiency_bonus"`

	// Skills (simplified for now)
	Skills map[string]bool `json:"skills"` // e.g., {"Acrobatics": true, "Stealth": false}

	// Inventory (simplified for now)
	Inventory []string `json:"inventory"`
}

// NewCharacter creates a new character with default values
func NewCharacter(name, species, class, background string) *Character {
	return &Character{
		Name:        name,
		Species:     species,
		Class:       class,
		Level:       1,
		Background:  background,
		Strength:    10,
		Dexterity:   10,
		Constitution: 10,
		Intelligence: 10,
		Wisdom:      10,
		Charisma:    10,
		HitPoints:   10, // Placeholder, should be class-dependent
		ArmorClass:  10, // Placeholder
		ProficiencyBonus: 2, // Level 1 proficiency bonus
		Skills:      make(map[string]bool),
		Inventory:   []string{},
	}
}

// SaveCharacter saves a character to a JSON file
func SaveCharacter(char *Character, filePath string) error {
	data, err := json.MarshalIndent(char, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write character file: %w", err)
	}
	return nil
}

// LoadCharacter loads a character from a JSON file
func LoadCharacter(filePath string) (*Character, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	var char Character
	err = json.Unmarshal(data, &char)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal character: %w", err)
	}
	return &char, nil
}

// GetCharacterFilePath returns the standard path for a character file
func GetCharacterFilePath(charName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	// Create a .dnd-cli directory in the user's home directory
	appDir := filepath.Join(homeDir, ".dnd-cli")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err = os.Mkdir(appDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create application directory: %w", err)
		}
	}

	return filepath.Join(appDir, fmt.Sprintf("%s.json", charName)), nil
}

// LevelUp increases character level and updates basic stats
func (c *Character) LevelUp() {
	c.Level++
	c.HitPoints += 5 // Simplified HP gain
	// Proficiency bonus increases at certain levels (e.g., 5, 9, 13, 17)
	if c.Level == 5 || c.Level == 9 || c.Level == 13 || c.Level == 17 {
		c.ProficiencyBonus++
	}
	// TODO: Implement more complex leveling logic (class features, spell slots, etc.)
}
