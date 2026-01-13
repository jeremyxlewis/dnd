package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	ArmorType        string                 `json:"armor_type"`
	HitPoints        int                    `json:"hit_points"`
	HitDice          string                 `json:"hit_dice"`
	Speed            string                 `json:"speed"`
	Stats            map[string]int         `json:"stats"`
	Skills           map[string]int         `json:"skills"`
	Senses           map[string]interface{} `json:"senses"`
	Languages        []string               `json:"languages"`
	Challenge        string                 `json:"challenge"`
	ChallengeXP      int                    `json:"challenge_xp"`
	SpecialAbilities []MonsterAbility       `json:"special_abilities"`
	Actions          []MonsterAbility       `json:"actions"`
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
	Name        string `json:"name"`
	Description string `json:"description"`
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
	// Structured data for detailed views
	AllSpells      []Spell
	AllMonsters    []Monster
	AllItems       []Item
	AllSpecies     []Species
	AllBackgrounds []Background
	AllClasses     []Class

	// Legacy data for comprehensive fuzzy search
	AllLegacySpells      []Spell
	AllLegacyMonsters    []Monster
	AllLegacyItems       []Item
	AllLegacySpecies     []Species
	AllLegacyBackgrounds []Background
	AllLegacyClasses     []Class

	// Structured data container
	StructuredData struct {
		Metadata struct {
			Version   string `json:"version"`
			Source    string `json:"source"`
			Extracted string `json:"extracted"`
		} `json:"metadata"`
		Spells      []Spell      `json:"spells"`
		Monsters    []Monster    `json:"monsters"`
		Items       []Item       `json:"items"`
		Species     []Species    `json:"species"`
		Backgrounds []Background `json:"backgrounds"`
		Classes     []Class      `json:"classes"`
	}
)

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

// parseSpellsFromJSON extracts spells from JSON structure
func parseSpellsFromJSON(data interface{}) []Spell {
	var spells []Spell

	// Handle both array and object formats
	switch v := data.(type) {
	case []interface{}:
		// Direct array format (nick-aschenbach structure)
		for _, item := range v {
			if spellData, ok := item.(map[string]interface{}); ok {
				if name, ok := spellData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := spellData["description"].(string); ok {
						description = desc
					}

					spellEntry := Spell{
						Name:        name,
						Description: description,
					}
					spells = append(spells, spellEntry)
				}
			}
		}
	case map[string]interface{}:
		// Object format (original structure)
		var findSpells func(obj interface{}) []Spell
		findSpells = func(obj interface{}) []Spell {
			var found []Spell

			switch v := obj.(type) {
			case map[string]interface{}:
				for key, value := range v {
					// Skip structural keys
					if key == "content" || key == "Spellcasting" {
						found = append(found, findSpells(value)...)
						continue
					}

					// Check if this looks like a spell
					if spellMap, ok := value.(map[string]interface{}); ok && spellMap["content"] != nil {
						content := spellMap["content"]
						if contentList, ok := content.([]interface{}); ok && len(contentList) > 0 {
							// Parse out structured spell data for fuzzy search
							spellEntry := Spell{
								Name:        key,
								Description: strings.Join(convertToStringSlice(contentList), "\n"),
							}
							found = append(found, spellEntry)
						}
					}
					// Recursively search
					found = append(found, findSpells(value)...)
				}
			case []interface{}:
				for _, item := range v {
					found = append(found, findSpells(item)...)
				}
			}

			return found
		}
		spells = findSpells(data)
	}

	return spells
}

// parseMonstersFromJSON extracts monsters from the complex nested JSON structure
func parseMonstersFromJSON(data interface{}) []Monster {
	var monsters []Monster

	// Handle both array and object formats
	switch v := data.(type) {
	case []interface{}:
		// Direct array format (from DND-5E-Data)
		for _, item := range v {
			if monsterData, ok := item.(map[string]interface{}); ok {
				if name, ok := monsterData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := monsterData["text"].(string); ok {
						description = desc
					} else if desc, ok := monsterData["description"].(string); ok {
						description = desc
					}

					monster := Monster{
						Name:        name,
						Description: description,
					}
					monsters = append(monsters, monster)
				}
			}
		}
	case map[string]interface{}:
		// Object format (from SRD)
		// Helper function to recursively find monsters
		var findMonsters func(obj interface{}) []Monster
		findMonsters = func(obj interface{}) []Monster {
			var found []Monster

			switch v := obj.(type) {
			case map[string]interface{}:
				for key, value := range v {
					// Skip structural keys
					if key == "content" || key == "Monsters" || key == "Modifying Creatures" || key == "Size" || key == "Type" {
						found = append(found, findMonsters(value)...)
						continue
					}

					// Check if this looks like a monster
					if monsterMap, ok := value.(map[string]interface{}); ok && monsterMap["content"] != nil {
						content := monsterMap["content"]
						if contentList, ok := content.([]interface{}); ok && len(contentList) > 1 {
							// Look for stat block patterns
							contentStr := strings.Join(convertToStringSlice(contentList), "\n")
							if strings.Contains(contentStr, "Armor Class") || strings.Contains(contentStr, "Hit Points") || strings.Contains(contentStr, "Speed") {
								// Parse out structured monster data for fuzzy search
								monster := Monster{
									Name:        key,
									Description: strings.TrimSpace(contentStr),
								}
								found = append(found, monster)
							}
						}
					}
					// Recursively search
					found = append(found, findMonsters(value)...)
				}
			case []interface{}:
				for _, item := range v {
					found = append(found, findMonsters(item)...)
				}
			}

			return found
		}
		monsters = findMonsters(data)
	}

	return monsters
}

// convertToStringSlice converts []interface{} to []string
func convertToStringSlice(slice []interface{}) []string {
	var result []string
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

// parseItemsFromJSON extracts items from JSON structure
func parseItemsFromJSON(data interface{}) []Item {
	var items []Item

	// Handle both array and object formats
	switch v := data.(type) {
	case []interface{}:
		// Direct array format (nick-aschenbach structure)
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if name, ok := itemMap["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := itemMap["description"].(string); ok {
						description = desc
					}

					item := Item{
						Name:        name,
						Description: description,
					}
					items = append(items, item)
				}
			}
		}
	case map[string]interface{}:
		// Object format (original structure)
		var findItems func(obj interface{}) []Item
		findItems = func(obj interface{}) []Item {
			var found []Item

			switch v := obj.(type) {
			case map[string]interface{}:
				for key, value := range v {
					// Skip structural keys
					if key == "content" {
						found = append(found, findItems(value)...)
						continue
					}

					// Check if this looks like an item
					if itemMap, ok := value.(map[string]interface{}); ok && itemMap["name"] != nil {
						name := ""
						description := ""

						if nameVal, ok := itemMap["name"].(string); ok {
							name = nameVal
						}
						if descVal, ok := itemMap["description"].(string); ok {
							description = descVal
						}

						if name != "" {
							item := Item{
								Name:        name,
								Description: description,
							}
							found = append(found, item)
						}
					}
					// Recursively search
					found = append(found, findItems(value)...)
				}
			case []interface{}:
				for _, item := range v {
					found = append(found, findItems(item)...)
				}
			}

			return found
		}
		items = findItems(data)
	}

	return items
}

// parseSpeciesFromJSON extracts species from JSON structure
func parseSpeciesFromJSON(data interface{}) []Species {
	var species []Species

	// Handle both array and object formats
	switch v := data.(type) {
	case []interface{}:
		// Direct array format
		for _, item := range v {
			if speciesData, ok := item.(map[string]interface{}); ok {
				if name, ok := speciesData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := speciesData["description"].(string); ok {
						description = desc
					}

					abilityBonuses := make(map[string]int)
					if bonuses, ok := speciesData["ability_bonuses"].(map[string]interface{}); ok {
						for stat, val := range bonuses {
							if valInt, ok := val.(float64); ok {
								abilityBonuses[stat] = int(valInt)
							}
						}
					}

					speed := ""
					if spd, ok := speciesData["speed"].(string); ok {
						speed = spd
					}

					var languages []string
					if langs, ok := speciesData["languages"].([]interface{}); ok {
						for _, lang := range langs {
							if langStr, ok := lang.(string); ok {
								languages = append(languages, langStr)
							}
						}
					}

					var traits []string
					if trts, ok := speciesData["traits"].([]interface{}); ok {
						for _, trait := range trts {
							if traitStr, ok := trait.(string); ok {
								traits = append(traits, traitStr)
							}
						}
					}

					speciesEntry := Species{
						Name:           name,
						Description:    description,
						AbilityBonuses: abilityBonuses,
						Speed:          speed,
						Languages:      languages,
						Traits:         traits,
					}
					species = append(species, speciesEntry)
				}
			}
		}
	case map[string]interface{}:
		// Object format
		for _, value := range v {
			if speciesData, ok := value.(map[string]interface{}); ok {
				if name, ok := speciesData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := speciesData["description"].(string); ok {
						description = desc
					}

					abilityBonuses := make(map[string]int)
					if bonuses, ok := speciesData["ability_bonuses"].(map[string]interface{}); ok {
						for stat, val := range bonuses {
							if valInt, ok := val.(float64); ok {
								abilityBonuses[stat] = int(valInt)
							}
						}
					}

					speed := ""
					if spd, ok := speciesData["speed"].(string); ok {
						speed = spd
					}

					var languages []string
					if langs, ok := speciesData["languages"].([]interface{}); ok {
						for _, lang := range langs {
							if langStr, ok := lang.(string); ok {
								languages = append(languages, langStr)
							}
						}
					}

					var traits []string
					if trts, ok := speciesData["traits"].([]interface{}); ok {
						for _, trait := range trts {
							if traitStr, ok := trait.(string); ok {
								traits = append(traits, traitStr)
							}
						}
					}

					speciesEntry := Species{
						Name:           name,
						Description:    description,
						AbilityBonuses: abilityBonuses,
						Speed:          speed,
						Languages:      languages,
						Traits:         traits,
					}
					species = append(species, speciesEntry)
				}
			}
		}
	}

	return species
}

func parseBackgroundsFromJSON(data interface{}) []Background {
	background_list := []string{
		"Acolyte", "Charlatan", "Criminal", "Entertainer", "Folk Hero",
		"Guild Artisan", "Hermit", "Noble", "Outlander", "Sage",
		"Sailor", "Soldier", "Urchin",
	}

	var backgrounds []Background
	for _, name := range background_list {
		backgrounds = append(backgrounds, Background{
			Name:        name,
			Description: "A D&D 5e background",
		})
	}

	return backgrounds
}

// parseClassesFromJSON creates basic classes
func parseClassesFromJSON(data interface{}) []Class {
	class_list := []string{
		"Artificer", "Barbarian", "Bard", "Cleric", "Druid",
		"Fighter", "Monk", "Paladin", "Ranger", "Rogue",
		"Sorcerer", "Warlock", "Wizard",
	}

	var classes []Class
	for _, name := range class_list {
		classes = append(classes, Class{
			Name:        name,
			Description: "A D&D 5e class",
		})
	}

	return classes
}

// loadLegacyDataForSearch loads comprehensive legacy data for fuzzy search
func loadLegacyDataForSearch(dataPath string) error {
	// Load Spells with special parsing for comprehensive search
	var spellsData interface{}
	err := loadJSONFile(filepath.Join(dataPath, "spells.json"), &spellsData)
	if err != nil {
		return fmt.Errorf("failed to load spells data: %w", err)
	}
	AllLegacySpells = parseSpellsFromJSON(spellsData)

	// Load Monsters with special parsing for comprehensive search
	var monstersData interface{}
	err = loadJSONFile(filepath.Join(dataPath, "monsters.json"), &monstersData)
	if err != nil {
		return fmt.Errorf("failed to load monsters data: %w", err)
	}
	AllLegacyMonsters = parseMonstersFromJSON(monstersData)

	// Load Items from JSON
	var itemsData interface{}
	err = loadJSONFile(filepath.Join(dataPath, "items.json"), &itemsData)
	if err != nil {
		return fmt.Errorf("failed to load items data: %w", err)
	}
	AllLegacyItems = parseItemsFromJSON(itemsData)

	// Load Species from JSON
	var speciesData interface{}
	err = loadJSONFile(filepath.Join(dataPath, "species.json"), &speciesData)
	if err != nil {
		return fmt.Errorf("failed to load species data: %w", err)
	}
	AllLegacySpecies = parseSpeciesFromJSON(speciesData)

	// Load Backgrounds from JSON
	var backgroundsData interface{}
	err = loadJSONFile(filepath.Join(dataPath, "backgrounds.json"), &backgroundsData)
	if err != nil {
		return fmt.Errorf("failed to load backgrounds data: %w", err)
	}
	AllLegacyBackgrounds = parseBackgroundsFromJSON(backgroundsData)

	// Load Classes from JSON
	var classesData interface{}
	err = loadJSONFile(filepath.Join(dataPath, "classes.json"), &classesData)
	if err != nil {
		// Use hardcoded fallback if classes.json has parsing errors
		AllLegacyClasses = parseClassesFromJSON(nil)
	} else {
		AllLegacyClasses = parseClassesFromJSON(classesData)
	}

	return nil
}

// LoadData loads all necessary JSON data into memory
func LoadData(dataPath string) error {
	// Load structured data for detailed views
	structuredPath := filepath.Join(dataPath, "structured.json")
	err := loadJSONFile(structuredPath, &StructuredData)
	if err != nil {
		// Use basic fallback if structured.json doesn't exist
		AllSpells = []Spell{}
		AllMonsters = []Monster{}
		AllItems = []Item{}
		AllSpecies = []Species{}
		AllBackgrounds = []Background{}
		AllClasses = []Class{}
	} else {
		// Assign structured data to global variables
		AllSpells = StructuredData.Spells
		AllMonsters = StructuredData.Monsters
		AllItems = StructuredData.Items
		AllSpecies = StructuredData.Species
		AllBackgrounds = StructuredData.Backgrounds
		AllClasses = StructuredData.Classes
	}

	// Always load legacy data for comprehensive fuzzy search
	err = loadLegacyDataForSearch(dataPath)
	if err != nil {
		return fmt.Errorf("failed to load legacy data: %w", err)
	}

	return nil
}

// GetSpellByName searches for a spell by its name (case-insensitive)
// Searches both structured and legacy data
func GetSpellByName(name string) (*Spell, error) {
	lowerName := strings.ToLower(name)

	// First search structured data
	for _, spell := range AllSpells {
		if strings.ToLower(spell.Name) == lowerName {
			return &spell, nil
		}
	}

	// Then search legacy data
	for _, spell := range AllLegacySpells {
		if strings.ToLower(spell.Name) == lowerName {
			return &spell, nil
		}
	}

	return nil, fmt.Errorf("spell '%s' not found", name)
}

// GetMonsterByName searches for a monster by its name (case-insensitive)
// Searches both structured and legacy data
func GetMonsterByName(name string) (*Monster, error) {
	lowerName := strings.ToLower(name)

	// First search structured data
	for _, monster := range AllMonsters {
		if strings.ToLower(monster.Name) == lowerName {
			return &monster, nil
		}
	}

	// Then search legacy data
	for _, monster := range AllLegacyMonsters {
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
