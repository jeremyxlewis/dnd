package data

import (
	"fmt"
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
	// SRD Parser instance
	srdParser *SRDParser
)

// LoadData loads all necessary JSON data into memory
func LoadData(dataPath string) error {
	// Initialize SRD parser
	srdParser = NewSRDParser()

	// Try to load SRD data from dnd-5e-srd directory
	srdPath := filepath.Join(dataPath, "..", "dnd-5e-srd", "5esrd.json")
	err := srdParser.LoadSRD(srdPath)
	if err != nil {
		return fmt.Errorf("failed to load SRD data: %w", err)
	}

	return nil
}

// GetSpellByName searches for a spell by its name (case-insensitive)
// Searches SRD data only
func GetSpellByName(name string) (*Spell, error) {
	// Search SRD data only
	if srdParser != nil {
		return srdParser.GetSpellByName(name)
	}

	return nil, fmt.Errorf("SRD data not loaded - spell '%s' not found", name)
}

// GetMonsterByName searches for a monster by its name (case-insensitive)
// Searches SRD data only
func GetMonsterByName(name string) (*Monster, error) {
	// Search SRD data only
	if srdParser != nil {
		return srdParser.GetMonsterByName(name)
	}

	return nil, fmt.Errorf("SRD data not loaded - monster '%s' not found", name)
}

// GetItemByName searches for an item by its name (case-insensitive)
// Searches SRD data only
func GetItemByName(name string) (*Item, error) {
	// Search SRD data only
	if srdParser != nil {
		return srdParser.GetItemByName(name)
	}

	return nil, fmt.Errorf("SRD data not loaded - item '%s' not found", name)
}

// GetSpeciesByName searches for a species by its name (case-insensitive)
// Searches SRD data only
func GetSpeciesByName(name string) (*Species, error) {
	// Search SRD data only
	if srdParser != nil {
		return srdParser.GetRaceByName(name)
	}

	return nil, fmt.Errorf("SRD data not loaded - species '%s' not found", name)
}

// GetBackgroundByName searches for a background by its name (case-insensitive)
// Searches SRD data only
func GetBackgroundByName(name string) (*Background, error) {
	// Search SRD data only
	if srdParser != nil {
		return srdParser.GetBackgroundByName(name)
	}

	return nil, fmt.Errorf("SRD data not loaded - background '%s' not found", name)
}

// GetClassByName searches for a class by its name (case-insensitive)
// Searches SRD data only
func GetClassByName(name string) (*Class, error) {
	// Search SRD data only
	if srdParser != nil {
		return srdParser.GetClassByName(name)
	}

	return nil, fmt.Errorf("SRD data not loaded - class '%s' not found", name)
}

// GetSRDSpellNames returns all spell names from SRD data
func GetSRDSpellNames() []string {
	if srdParser != nil {
		return srdParser.GetSpellNames()
	}
	return []string{}
}

// GetSRDMonsterNames returns all monster names from SRD data
func GetSRDMonsterNames() []string {
	if srdParser != nil {
		return srdParser.GetMonsterNames()
	}
	return []string{}
}

// GetSRDItemNames returns all item names from SRD data
func GetSRDItemNames() []string {
	if srdParser != nil {
		return srdParser.GetItemNames()
	}
	return []string{}
}

// GetSRDClassNames returns all class names from SRD data
func GetSRDClassNames() []string {
	if srdParser != nil {
		return srdParser.GetClassNames()
	}
	return []string{}
}

// GetSRDRaceNames returns all race names from SRD data
func GetSRDRaceNames() []string {
	if srdParser != nil {
		return srdParser.GetRaceNames()
	}
	return []string{}
}

// GetSRDBackgroundNames returns all background names from SRD data
func GetSRDBackgroundNames() []string {
	if srdParser != nil {
		return srdParser.GetBackgroundNames()
	}
	return []string{}
}
