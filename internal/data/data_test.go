package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDataAndLookup(t *testing.T) {
	// Create a temporary directory for test data
	testDir, err := os.MkdirTemp("", "dnd-data-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up after test

	// Create dummy spells.json
	dummySpells := `[
		{"name":"Test Spell 1","description":"Desc 1","properties":{},"publisher":"Pub 1","book":"Book 1"},
		{"name":"Test Spell 2","description":"Desc 2","properties":{},"publisher":"Pub 2","book":"Book 2"}
	]`
	if err := os.WriteFile(filepath.Join(testDir, "spells.json"), []byte(dummySpells), 0644); err != nil {
		t.Fatalf("Failed to write dummy spells.json: %v", err)
	}

	// Create dummy monsters.json
	dummyMonsters := `[
		{"name":"Test Monster 1","description":"Monster Desc 1"},
		{"name":"Test Monster 2","description":"Monster Desc 2"}
	]`
	if err := os.WriteFile(filepath.Join(testDir, "monsters.json"), []byte(dummyMonsters), 0644); err != nil {
		t.Fatalf("Failed to write dummy monsters.json: %v", err)
	}

	// Create dummy items.json
	dummyItems := `[
		{"name":"Test Item 1","description":"Item Desc 1"},
		{"name":"Test Item 2","description":"Item Desc 2"}
	]`
	if err := os.WriteFile(filepath.Join(testDir, "items.json"), []byte(dummyItems), 0644); err != nil {
		t.Fatalf("Failed to write dummy items.json: %v", err)
	}

	// Create dummy species.json
	dummySpecies := `[
		{"name":"Human","description":""},
		{"name":"Elf","description":""}
	]`
	if err := os.WriteFile(filepath.Join(testDir, "species.json"), []byte(dummySpecies), 0644); err != nil {
		t.Fatalf("Failed to write dummy species.json: %v", err)
	}

	// Create dummy backgrounds.json
	dummyBackgrounds := `[
		{"name":"Acolyte","description":""},
		{"name":"Soldier","description":""}
	]`
	if err := os.WriteFile(filepath.Join(testDir, "backgrounds.json"), []byte(dummyBackgrounds), 0644); err != nil {
		t.Fatalf("Failed to write dummy backgrounds.json: %v", err)
	}

	// Create dummy classes.json
	dummyClasses := `[
		{"name":"Fighter","description":""},
		{"name":"Wizard","description":""}
	]`
	if err := os.WriteFile(filepath.Join(testDir, "classes.json"), []byte(dummyClasses), 0644); err != nil {
		t.Fatalf("Failed to write dummy classes.json: %v", err)
	}

	// Load data from the temporary directory
	err = LoadData(testDir)
	if err != nil {
		t.Fatalf("LoadData failed: %v", err)
	}

	// Test spell lookup
	spell, err := GetSpellByName("Test Spell 1")
	if err != nil || spell.Name != "Test Spell 1" {
		t.Errorf("GetSpellByName failed for 'Test Spell 1': %v, got %v", err, spell)
	}
	_, err = GetSpellByName("Non Existent Spell")
	if err == nil {
		t.Errorf("GetSpellByName expected error for non-existent spell, got nil")
	}

	// Test monster lookup
	monster, err := GetMonsterByName("Test Monster 1")
	if err != nil || monster.Name != "Test Monster 1" {
		t.Errorf("GetMonsterByName failed for 'Test Monster 1': %v, got %v", err, monster)
	}
	_, err = GetMonsterByName("Non Existent Monster")
	if err == nil {
		t.Errorf("GetMonsterByName expected error for non-existent monster, got nil")
	}

	// Test item lookup
	item, err := GetItemByName("Test Item 1")
	if err != nil || item.Name != "Test Item 1" {
		t.Errorf("GetItemByName failed for 'Test Item 1': %v, got %v", err, item)
	}
	_, err = GetItemByName("Non Existent Item")
	if err == nil {
		t.Errorf("GetItemByName expected error for non-existent item, got nil")
	}

	// Test NPC generation
	npc := GenerateNPC()
	if npc.Name == "" || npc.Species == "" || npc.Background == "" {
		t.Errorf("GenerateNPC returned empty string: Name='%s', Species='%s', Background='%s'", npc.Name, npc.Species, npc.Background)
	}
}
