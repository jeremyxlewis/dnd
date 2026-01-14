package data

import (
	"fmt"
	"os"
	"path/filepath"
)

// Spell represents the structure of a spell
type Spell struct {
	Name         string   `json:"name"`
	Level        int      `json:"level"`
	School       string   `json:"school"`
	CastingTime  string   `json:"casting_time"`
	Range        string   `json:"range"`
	Components   []string `json:"components"`
	Material     string   `json:"material,omitempty"`
	Duration     string   `json:"duration"`
	Description  string   `json:"description"`
	HigherLevels *string  `json:"higher_levels,omitempty"`
	Classes      []string `json:"classes,omitempty"`
}

// Monster represents the structure of a monster
type Monster struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Size             string                 `json:"size"`
	Type             string                 `json:"type"`
	Alignment        string                 `json:"alignment"`
	ArmorClass       int                    `json:"armor_class"`
	ArmorClassStr    string                 `json:"ac,omitempty"`
	ArmorType        string                 `json:"armor_type"`
	HitPoints        int                    `json:"hit_points"`
	HitPointsStr     string                 `json:"hp,omitempty"`
	HitDice          string                 `json:"hit_dice"`
	Speed            string                 `json:"speed"`
	Stats            map[string]int         `json:"stats"`
	Str              string                 `json:"str,omitempty"`
	Dex              string                 `json:"dex,omitempty"`
	Con              string                 `json:"con,omitempty"`
	Int              string                 `json:"int,omitempty"`
	Wis              string                 `json:"wis,omitempty"`
	Cha              string                 `json:"cha,omitempty"`
	Skills           map[string]int         `json:"skills"`
	SkillsArray      []string               `json:"skill,omitempty"`
	Senses           map[string]interface{} `json:"senses"`
	Passive          string                 `json:"passive,omitempty"`
	Languages        []string               `json:"languages"`
	LanguagesStr     string                 `json:"languages,omitempty"`
	Challenge        string                 `json:"cr"`
	ChallengeXP      int                    `json:"challenge_xp"`
	SpecialAbilities []MonsterAbility       `json:"special_abilities"`
	Actions          []MonsterAbility       `json:"actions"`
	Traits           []interface{}          `json:"trait,omitempty"`
	ActionsArray     []interface{}          `json:"action,omitempty"`
}

// MonsterAbility represents an ability or action for monsters
type MonsterAbility struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description"`
	Hit         string `json:"hit,omitempty"`
	Reach       string `json:"reach,omitempty"`
	Range       string `json:"range,omitempty"`
	Damage      string `json:"damage,omitempty"`
}

// Item represents the structure of an item
type Item struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Text        []string `json:"text,omitempty"`
	Type        string   `json:"type,omitempty"`
	Value       string   `json:"value,omitempty"`
	Weight      string   `json:"weight,omitempty"`
	Damage      string   `json:"dmg1,omitempty"`
	DamageType  string   `json:"dmgType,omitempty"`
	Property    string   `json:"property,omitempty"`
}

// Species represents the structure of a species
type Species struct {
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	AbilityBonuses map[string]int `json:"ability_bonuses"`
	Speed          string         `json:"speed"`
	Languages      []string       `json:"languages"`
	Traits         []string       `json:"traits"`
}

// Background represents the structure of a background
type Background struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	SkillProficiencies []string `json:"skill_proficiencies"`
	ToolProficiencies  []string `json:"tool_proficiencies"`
	Equipment          []string `json:"equipment"`
	Feature            string   `json:"feature"`
	PersonalityTraits  []string `json:"personality_traits"`
	Ideals             []string `json:"ideals"`
}

// Class represents the structure of a class
type Class struct {
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	HitDie              string   `json:"hit_die"`
	PrimaryAbility      string   `json:"primary_ability"`
	SavingThrows        []string `json:"saving_throws"`
	ArmorProficiencies  []string `json:"armor_proficiencies"`
	WeaponProficiencies []string `json:"weapon_proficiencies"`
	ToolsProficiencies  []string `json:"tools_proficiencies"`
	SkillsCount         int      `json:"skills_count"`
	SkillsChoices       []string `json:"skills_choices"`
	Equipment           []string `json:"equipment"`
}

// Global store for loaded data
var (
	// SRD Parser instance (legacy)
	srdParser *SRDParser
	// Compendium Parser instance
	compendiumParser *CompendiumParser
	// Currently loaded compendium file
	currentCompendiumFile string
)

// LoadData loads all necessary data from Compendium XML files
func LoadData(dataPath string) error {
	// Use Compendium directory
	compendiumDir := dataPath
	if filepath.Base(dataPath) != "Compendium" {
		compendiumDir = filepath.Join(dataPath, "Compendium")
	}

	// Check if Compendium directory exists
	if _, err := os.Stat(compendiumDir); os.IsNotExist(err) {
		return fmt.Errorf("Compendium directory not found: %s", compendiumDir)
	}

	// Get list of available compendium files
	files, err := GetCompendiumFileList(compendiumDir)
	if err != nil {
		return fmt.Errorf("failed to list compendium files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no compendium files found in %s", compendiumDir)
	}

	// Default to Complete_Compendium.xml if available, otherwise use first file
	compendiumFile := GetDefaultCompendiumFile()
	found := false
	for _, file := range files {
		if file == compendiumFile {
			found = true
			break
		}
	}
	if !found {
		compendiumFile = files[0]
	}

	return LoadCompendiumData(compendiumDir, compendiumFile)
}

// LoadCompendiumData loads a specific compendium file
func LoadCompendiumData(compendiumDir, filename string) error {
	// Validate file exists
	if err := ValidateCompendiumFile(compendiumDir, filename); err != nil {
		return err
	}

	// Initialize Compendium parser
	compendiumParser = NewCompendiumParser()

	// Load the compendium file
	compendiumPath := filepath.Join(compendiumDir, filename)
	err := compendiumParser.LoadCompendium(compendiumPath)
	if err != nil {
		return fmt.Errorf("failed to load compendium file '%s': %w", filename, err)
	}

	currentCompendiumFile = filename
	return nil
}

// GetAvailableCompendiumFiles returns list of available compendium files
func GetAvailableCompendiumFiles() ([]string, error) {
	// Use Compendium directory
	compendiumDir := "./Compendium"
	files, err := GetCompendiumFileList(compendiumDir)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetCurrentCompendiumFile returns the currently loaded compendium file
func GetCurrentCompendiumFile() string {
	return currentCompendiumFile
}

// GetSpellByName searches for a spell by its name (case-insensitive)
// Searches Compendium data
func GetSpellByName(name string) (*Spell, error) {
	// Search Compendium data first
	if compendiumParser != nil {
		return compendiumParser.GetSpellByName(name)
	}

	// Fallback to SRD data if available
	if srdParser != nil {
		return srdParser.GetSpellByName(name)
	}

	return nil, fmt.Errorf("No data loaded - spell '%s' not found", name)
}

// GetMonsterByName searches for a monster by its name (case-insensitive)
// Searches Compendium data
func GetMonsterByName(name string) (*Monster, error) {
	// Search Compendium data first
	if compendiumParser != nil {
		return compendiumParser.GetMonsterByName(name)
	}

	// Fallback to SRD data if available
	if srdParser != nil {
		return srdParser.GetMonsterByName(name)
	}

	return nil, fmt.Errorf("No data loaded - monster '%s' not found", name)
}

// GetItemByName searches for an item by its name (case-insensitive)
// Searches Compendium data
func GetItemByName(name string) (*Item, error) {
	// Search Compendium data first
	if compendiumParser != nil {
		return compendiumParser.GetItemByName(name)
	}

	// Fallback to SRD data if available
	if srdParser != nil {
		return srdParser.GetItemByName(name)
	}

	return nil, fmt.Errorf("No data loaded - item '%s' not found", name)
}

// GetSpeciesByName searches for a species by its name (case-insensitive)
// Searches Compendium data
func GetSpeciesByName(name string) (*Species, error) {
	// Search Compendium data first
	if compendiumParser != nil {
		return compendiumParser.GetRaceByName(name)
	}

	// Fallback to SRD data if available
	if srdParser != nil {
		return srdParser.GetRaceByName(name)
	}

	return nil, fmt.Errorf("No data loaded - species '%s' not found", name)
}

// GetBackgroundByName searches for a background by its name (case-insensitive)
// Searches Compendium data
func GetBackgroundByName(name string) (*Background, error) {
	// Search Compendium data first
	if compendiumParser != nil {
		return compendiumParser.GetBackgroundByName(name)
	}

	// Fallback to SRD data if available
	if srdParser != nil {
		return srdParser.GetBackgroundByName(name)
	}

	return nil, fmt.Errorf("No data loaded - background '%s' not found", name)
}

// GetClassByName searches for a class by its name (case-insensitive)
// Searches Compendium data
func GetClassByName(name string) (*Class, error) {
	// Search Compendium data first
	if compendiumParser != nil {
		return compendiumParser.GetClassByName(name)
	}

	// Fallback to SRD data if available
	if srdParser != nil {
		return srdParser.GetClassByName(name)
	}

	return nil, fmt.Errorf("No data loaded - class '%s' not found", name)
}

// GetSRDSpellNames returns all spell names from loaded data
func GetSRDSpellNames() []string {
	if compendiumParser != nil {
		return compendiumParser.GetSpellNames()
	}
	if srdParser != nil {
		return srdParser.GetSpellNames()
	}
	return []string{}
}

// GetSRDMonsterNames returns all monster names from loaded data
func GetSRDMonsterNames() []string {
	if compendiumParser != nil {
		return compendiumParser.GetMonsterNames()
	}
	if srdParser != nil {
		return srdParser.GetMonsterNames()
	}
	return []string{}
}

// GetSRDItemNames returns all item names from loaded data
func GetSRDItemNames() []string {
	if compendiumParser != nil {
		return compendiumParser.GetItemNames()
	}
	if srdParser != nil {
		return srdParser.GetItemNames()
	}
	return []string{}
}

// GetSRDClassNames returns all class names from loaded data
func GetSRDClassNames() []string {
	if compendiumParser != nil {
		return compendiumParser.GetClassNames()
	}
	if srdParser != nil {
		return srdParser.GetClassNames()
	}
	return []string{}
}

// GetSRDRaceNames returns all race names from loaded data
func GetSRDRaceNames() []string {
	if compendiumParser != nil {
		return compendiumParser.GetRaceNames()
	}
	if srdParser != nil {
		return srdParser.GetRaceNames()
	}
	return []string{}
}

// GetSRDBackgroundNames returns all background names from loaded data
func GetSRDBackgroundNames() []string {
	if compendiumParser != nil {
		return compendiumParser.GetBackgroundNames()
	}
	if srdParser != nil {
		return srdParser.GetBackgroundNames()
	}
	return []string{}
}
