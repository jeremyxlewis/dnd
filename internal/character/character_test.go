package character

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCharacter(t *testing.T) {
	char := NewCharacter("TestChar", "Human", "Fighter", "Soldier", 10, 10, 10, 10, 10, 10)

	if char.Name != "TestChar" {
		t.Errorf("Expected name TestChar, got %s", char.Name)
	}
	if char.Level != 1 {
		t.Errorf("Expected level 1, got %d", char.Level)
	}
	// Add more assertions for default values
}

func TestSaveLoadCharacter(t *testing.T) {
	// Create a temporary directory for test files
	testDir, err := os.MkdirTemp("", "char-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up after test

	originalChar := NewCharacter("SaveLoadChar", "Elf", "Rogue", "Criminal", 10, 10, 10, 10, 10, 10)
	charFilePath := filepath.Join(testDir, "SaveLoadChar.json")

	// Save character
	err = SaveCharacter(originalChar, charFilePath)
	if err != nil {
		t.Fatalf("SaveCharacter failed: %v", err)
	}

	// Load character
	loadedChar, err := LoadCharacter(charFilePath)
	if err != nil {
		t.Fatalf("LoadCharacter failed: %v", err)
	}

	// Compare original and loaded character
	if originalChar.Name != loadedChar.Name {
		t.Errorf("Loaded character name mismatch: expected %s, got %s", originalChar.Name, loadedChar.Name)
	}
	// Add more comparisons for other fields
}

func TestLevelUp(t *testing.T) {
	char := NewCharacter("LevelUpChar", "Dwarf", "Cleric", "Acolyte", 10, 10, 10, 10, 10, 10)

	initialLevel := char.Level
	initialHP := char.HitPoints
	initialProficiency := char.ProficiencyBonus

	char.LevelUp()

	if char.Level != initialLevel+1 {
		t.Errorf("Expected level %d, got %d", initialLevel+1, char.Level)
	}
	if char.HitPoints != initialHP+5 {
		t.Errorf("Expected HP %d, got %d", initialHP+5, char.HitPoints)
	}
	// Test proficiency bonus increase at specific levels
	char.Level = 4 // Set to level 4, next level up should increase proficiency
	char.ProficiencyBonus = 2
	char.LevelUp()
	if char.ProficiencyBonus != initialProficiency+1 {
		t.Errorf("Expected proficiency bonus %d, got %d after level 5", initialProficiency+1, char.ProficiencyBonus)
	}
}

func TestGetCharacterFilePath(t *testing.T) {
	charName := "TestPathChar"
	filePath, err := GetCharacterFilePath(charName)
	if err != nil {
		t.Fatalf("GetCharacterFilePath failed: %v", err)
	}

	// Check if the path ends with .dnd-cli/TestPathChar.json
	if !strings.HasSuffix(filePath, filepath.Join(".dnd-cli", charName+".json")) {
		t.Errorf("Expected path to end with .dnd-cli/%s.json, got %s", charName, filePath)
	}

	// Ensure the directory is created
	appDir := filepath.Dir(filePath)
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		t.Errorf("Application directory %s was not created", appDir)
	}
	defer os.RemoveAll(appDir) // Clean up the created directory
}
