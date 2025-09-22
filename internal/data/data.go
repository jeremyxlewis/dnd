package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Spell represents the structure of a spell from spells.json
type Spell struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
	Publisher   string                 `json:"publisher"`
	Book        string                 `json:"book"`
}

// Monster represents a simplified structure of a monster from monsters.json
type Monster struct {
	Name        string `json:"name"`
	Description string `json:"description"` // This might need to be inferred or constructed from other fields
}

// Item represents a simplified structure of an item from items.json
type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Species represents the structure of a species from species.json
type Species struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Background represents the structure of a background from backgrounds.json
type Background struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Class represents the structure of a class from classes.json
type Class struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Global store for loaded data
var (
	AllSpells      []Spell
	AllMonsters    []Monster
	AllItems       []Item
	AllSpecies     []Species
	AllBackgrounds []Background
	AllClasses     []Class
)

// LoadData loads all necessary JSON data into memory
func LoadData(dataPath string) error {
	// Load Spells
	err := loadJSONFile(filepath.Join(dataPath, "spells.json"), &AllSpells)
	if err != nil {
		return fmt.Errorf("failed to load spells data: %w", err)
	}

	// Load Monsters (simplified)
	err = loadJSONFile(filepath.Join(dataPath, "monsters.json"), &AllMonsters)
	if err != nil {
		return fmt.Errorf("failed to load monsters data: %w", err)
	}

	// Load Items (simplified)
	err = loadJSONFile(filepath.Join(dataPath, "items.json"), &AllItems)
	if err != nil {
		return fmt.Errorf("failed to load items data: %w", err)
	}

	// Load Species
	err = loadJSONFile(filepath.Join(dataPath, "species.json"), &AllSpecies)
	if err != nil {
		return fmt.Errorf("failed to load species data: %w", err)
	}

	// Load Backgrounds
	err = loadJSONFile(filepath.Join(dataPath, "backgrounds.json"), &AllBackgrounds)
	if err != nil {
		return fmt.Errorf("failed to load backgrounds data: %w", err)
	}

	// Load Classes
	err = loadJSONFile(filepath.Join(dataPath, "classes.json"), &AllClasses)
	if err != nil {
		return fmt.Errorf("failed to load classes data: %w", err)
	}

	return nil
}

// Helper function to load and unmarshal JSON files
func loadJSONFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}
	return nil
}

// GetSpellByName searches for a spell by its name (case-insensitive)
func GetSpellByName(name string) (*Spell, error) {
	lowerName := strings.ToLower(name)
	for _, spell := range AllSpells {
		if strings.ToLower(spell.Name) == lowerName {
			return &spell, nil
		}
	}
	return nil, fmt.Errorf("spell '%s' not found", name)
}

// GetMonsterByName searches for a monster by its name (case-insensitive)
func GetMonsterByName(name string) (*Monster, error) {
	lowerName := strings.ToLower(name)
	for _, monster := range AllMonsters {
		if strings.ToLower(monster.Name) == lowerName {
			return &monster, nil
		}
	}
	return nil, fmt.Errorf("monster '%s' not found", name)
}

// GetItemByName searches for an item by its name (case-insensitive)
func GetItemByName(name string) (*Item, error) {
	lowerName := strings.ToLower(name)
	for _, item := range AllItems {
		if strings.ToLower(item.Name) == lowerName {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item '%s' not found", name)
}

// GetSpeciesByName searches for a species by its name (case-insensitive)
func GetSpeciesByName(name string) (*Species, error) {
	lowerName := strings.ToLower(name)
	for _, species := range AllSpecies {
		if strings.ToLower(species.Name) == lowerName {
			return &species, nil
		}
	}
	return nil, fmt.Errorf("species '%s' not found", name)
}

// GetBackgroundByName searches for a background by its name (case-insensitive)
func GetBackgroundByName(name string) (*Background, error) {
	lowerName := strings.ToLower(name)
	for _, background := range AllBackgrounds {
		if strings.ToLower(background.Name) == lowerName {
			return &background, nil
		}
	}
	return nil, fmt.Errorf("background '%s' not found", name)
}

// GetClassByName searches for a class by its name (case-insensitive)
func GetClassByName(name string) (*Class, error) {
	lowerName := strings.ToLower(name)
	for _, class := range AllClasses {
		if strings.ToLower(class.Name) == lowerName {
			return &class, nil
		}
	}
	return nil, fmt.Errorf("class '%s' not found", name)
}

// NPC represents a generated non-player character
type NPC struct {
	Name             string
	Species          string
	Background       string
	PersonalityTrait string
	Ideal            string
	Bond             string
	Flaw             string
	Backstory        string
}

// GenerateNPC generates a random NPC with expanded details
func GenerateNPC() NPC {
	rand.Seed(time.Now().UnixNano())

	// Random Name
	firstNames := []string{"Elara", "Borin", "Lyra", "Gareth", "Seraphina", "Kaelen", "Thrain", "Mira", "Dorian", "Lirael"}
	lastNames := []string{"Stonehand", "Brightwood", "Shadowbrook", "Ironhide", "Whisperwind", "Darkwater", "Stormforge", "Moonshadow", "Fireheart", "Oakenshield"}
	name := fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))])

	// Random Species
	species := "Human"
	if len(AllSpecies) > 0 {
		species = AllSpecies[rand.Intn(len(AllSpecies))].Name
	}

	// Random Background
	background := "Commoner"
	if len(AllBackgrounds) > 0 {
		background = AllBackgrounds[rand.Intn(len(AllBackgrounds))].Name
	}

	// Random Personality Traits
	personalityTraits := []string{
		"I idolize a particular hero of my faith, and constantly refer to their deeds and example.",
		"I see omens in death.",
		"I am suspicious of strangers and expect the worst of them.",
		"I am tolerant (or intolerant) of other faiths and respect (or condemn) the worship of other gods.",
		"I've enjoyed fine food, drink, and high society among my temple's elite.",
		"I've spent so long in the temple that I have little practical experience dealing with people in the outside world.",
	}
	personalityTrait := personalityTraits[rand.Intn(len(personalityTraits))]

	// Random Ideals
	ideals := []string{
		"Change: We must help bring about the changes the gods are constantly working in the world.",
		"Tradition: The ancient ways of our faith must be preserved.",
		"Charity: I always try to help those in need, no matter what the personal cost.",
		"Power: I hope to one day rise to the top of my faith's religious hierarchy.",
		"Faith: I trust that my deity will guide my actions. I have faith that if I work hard, things will go well.",
		"Aspiration: I seek to prove myself worthy of my god's favor by matching my actions against his or her teachings.",
	}
	ideal := ideals[rand.Intn(len(ideals))]

	// Random Bonds
	bonds := []string{
		"I owe my life to the priest who took me in when my parents died.",
		"Everything I do is for the common people.",
		"I will someday get revenge on the corrupt temple hierarchy who branded me a heretic.",
		"I owe a debt I can never repay to the person who took pity on me.",
		"I protect those who cannot protect themselves.",
		"My temple is my home, and I will defend it with my life.",
	}
	bond := bonds[rand.Intn(len(bonds))]

	// Random Flaws
	flaws := []string{
		"I am suspicious of strangers and expect the worst of them.",
		"Once I pick a goal, I become obsessed with it to the detriment of everything else in my life.",
		"I am suspicious of strangers and expect the worst of them.",
		"I am suspicious of strangers and expect the worst of them.",
		"I am suspicious of strangers and expect the worst of them.",
		"I am suspicious of strangers and expect the worst of them.",
	}
	// Wait, I repeated. Let me fix.
	flaws = []string{
		"I am suspicious of strangers and expect the worst of them.",
		"Once I pick a goal, I become obsessed with it to the detriment of everything else in my life.",
		"I am greedy and will do anything for money.",
		"I am prone to fits of rage when provoked.",
		"I have a secret that could ruin me if it got out.",
		"I am overly trusting of others.",
	}
	flaw := flaws[rand.Intn(len(flaws))]

	// Random Backstory Snippet
	backstorySnippets := []string{
		"Born in a small village, I was always drawn to the mysteries of the divine. After a near-death experience, I dedicated my life to serving the gods.",
		"Growing up in the shadow of a great temple, I learned the ways of faith from a young age. Now, I travel the lands spreading the word.",
		"Orphaned at a young age, I found solace in the teachings of my deity. The temple became my family, and I vowed to protect it.",
		"As a child, I witnessed a miracle that changed my life forever. Since then, I've sought to understand and replicate such wonders.",
		"Exiled from my homeland for heretical beliefs, I now wander as a pilgrim, seeking truth in distant lands.",
		"Raised by devout parents, I inherited their faith and now carry on their legacy as a humble servant of the divine.",
	}
	backstory := backstorySnippets[rand.Intn(len(backstorySnippets))]

	return NPC{
		Name:             name,
		Species:          species,
		Background:       background,
		PersonalityTrait: personalityTrait,
		Ideal:            ideal,
		Bond:             bond,
		Flaw:             flaw,
		Backstory:        backstory,
	}
}
