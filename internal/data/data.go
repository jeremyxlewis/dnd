package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

// processStructuredMonsters processes structured monsters to convert string fields to proper types
func processStructuredMonsters(monsters []Monster) []Monster {
	for i := range monsters {
		monster := &monsters[i]

		// Convert AC string to int
		if monster.ArmorClassStr != "" {
			if ac, err := strconv.Atoi(monster.ArmorClassStr); err == nil {
				monster.ArmorClass = ac
			}
		}

		// Convert HP string to int (extract first number)
		if monster.HitPointsStr != "" {
			if hpMatch := regexp.MustCompile(`^(\d+)`).FindStringSubmatch(monster.HitPointsStr); len(hpMatch) > 1 {
				if hp, err := strconv.Atoi(hpMatch[1]); err == nil {
					monster.HitPoints = hp
				}
			}
		}

		// Parse individual stats
		monster.Stats = make(map[string]int)
		statFields := map[string]string{"str": monster.Str, "dex": monster.Dex, "con": monster.Con, "int": monster.Int, "wis": monster.Wis, "cha": monster.Cha}
		for stat, value := range statFields {
			if value != "" {
				if statInt, err := strconv.Atoi(value); err == nil {
					monster.Stats[stat] = statInt
				}
			}
		}

		// Parse skills from array
		if len(monster.SkillsArray) > 0 {
			monster.Skills = make(map[string]int)
			for _, skillStr := range monster.SkillsArray {
				// Handle format like "Arcana +8, History +8" (multiple skills in one string)
				skills := strings.Split(skillStr, ",")
				for _, skill := range skills {
					skill = strings.TrimSpace(skill)
					if skillMatch := regexp.MustCompile(`(.+)\s+([+-]\d+)`).FindStringSubmatch(skill); len(skillMatch) > 2 {
						skillName := strings.TrimSpace(skillMatch[1])
						if modifier, err := strconv.Atoi(skillMatch[2]); err == nil {
							monster.Skills[skillName] = modifier
						}
					}
				}
			}
		}

		// Parse senses
		if monster.Passive != "" {
			if monster.Senses == nil {
				monster.Senses = make(map[string]interface{})
			}
			monster.Senses["passive Perception"] = monster.Passive
		}

		// Parse languages
		if monster.LanguagesStr != "" && len(monster.Languages) == 0 {
			monster.Languages = strings.Split(monster.LanguagesStr, ", ")
		}

		// Parse traits and actions from raw arrays
		if len(monster.Traits) > 0 {
			monster.SpecialAbilities = parseMonsterAbilities(monster.Traits)
		}
		if len(monster.ActionsArray) > 0 {
			monster.Actions = parseMonsterAbilities(monster.ActionsArray)
		}
	}
	return monsters
}

// parseMonsterAbilities converts raw ability arrays to MonsterAbility structs
func parseMonsterAbilities(abilities []interface{}) []MonsterAbility {
	var result []MonsterAbility
	for _, abilityInterface := range abilities {
		if abilityMap, ok := abilityInterface.(map[string]interface{}); ok {
			ability := MonsterAbility{}
			if name, ok := abilityMap["name"].(string); ok {
				ability.Name = name
			}
			if textArray, ok := abilityMap["text"].([]interface{}); ok {
				var textParts []string
				for _, text := range textArray {
					if textStr, ok := text.(string); ok {
						textParts = append(textParts, textStr)
					}
				}
				ability.Description = strings.Join(textParts, " ")
			}
			result = append(result, ability)
		}
	}
	return result
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

					// Try multiple fields for description
					if desc, ok := monsterData["text"].(string); ok {
						description = desc
					} else if desc, ok := monsterData["description"].(string); ok {
						description = desc
					} else {
						// Extract from traits and actions for structured data
						var descriptionParts []string

						// Extract from traits
						if traits, ok := monsterData["trait"].([]interface{}); ok {
							for _, trait := range traits {
								if traitMap, ok := trait.(map[string]interface{}); ok {
									if traitName, ok := traitMap["name"].(string); ok {
										if traitText, ok := traitMap["text"].([]interface{}); ok {
											var textParts []string
											for _, text := range traitText {
												if textStr, ok := text.(string); ok && textStr != "" {
													textParts = append(textParts, textStr)
												}
											}
											if len(textParts) > 0 {
												descriptionParts = append(descriptionParts, fmt.Sprintf("%s: %s", traitName, strings.Join(textParts, " ")))
											}
										}
									}
								}
							}
						}

						// Extract from actions
						if actions, ok := monsterData["action"].([]interface{}); ok {
							for _, action := range actions {
								if actionMap, ok := action.(map[string]interface{}); ok {
									if actionName, ok := actionMap["name"].(string); ok {
										if actionText, ok := actionMap["text"].([]interface{}); ok {
											var textParts []string
											for _, text := range actionText {
												if textStr, ok := text.(string); ok && textStr != "" {
													textParts = append(textParts, textStr)
												}
											}
											if len(textParts) > 0 {
												descriptionParts = append(descriptionParts, fmt.Sprintf("%s: %s", actionName, strings.Join(textParts, " ")))
											}
										}
									}
								}
							}
						}

						if len(descriptionParts) > 0 {
							description = strings.Join(descriptionParts, "\n")
						}
					}

					// Extract additional fields
					monster := Monster{
						Name: name,
					}

					// Parse basic stats
					if size, ok := monsterData["size"].(string); ok {
						monster.Size = size
					}
					if monsterType, ok := monsterData["type"].(string); ok {
						monster.Type = monsterType
					}
					if alignment, ok := monsterData["alignment"].(string); ok {
						monster.Alignment = alignment
					}
					if ac, ok := monsterData["ac"].(string); ok {
						if acInt, err := strconv.Atoi(ac); err == nil {
							monster.ArmorClass = acInt
						}
					}
					if hp, ok := monsterData["hp"].(string); ok {
						// Extract first number from hp string like "210 (20d10+100)"
						if hpMatch := regexp.MustCompile(`^(\d+)`).FindStringSubmatch(hp); len(hpMatch) > 1 {
							if hpInt, err := strconv.Atoi(hpMatch[1]); err == nil {
								monster.HitPoints = hpInt
							}
						}
					}
					if speed, ok := monsterData["speed"].(string); ok {
						monster.Speed = speed
					}

					// Parse individual stats (str, dex, con, int, wis, cha)
					monster.Stats = make(map[string]int)
					statNames := []string{"str", "dex", "con", "int", "wis", "cha"}
					for _, statName := range statNames {
						if statValue, ok := monsterData[statName].(string); ok {
							if statInt, err := strconv.Atoi(statValue); err == nil {
								monster.Stats[statName] = statInt
							}
						}
					}

					// Parse skills (array format like ["Perception +5"])
					if skillsArray, ok := monsterData["skill"].([]interface{}); ok {
						monster.Skills = make(map[string]int)
						for _, skillInterface := range skillsArray {
							if skillStr, ok := skillInterface.(string); ok {
								// Parse skill format like "Perception +5"
								if skillMatch := regexp.MustCompile(`(.+)\s+([+-]\d+)`).FindStringSubmatch(skillStr); len(skillMatch) > 2 {
									skillName := strings.TrimSpace(skillMatch[1])
									if modifier, err := strconv.Atoi(skillMatch[2]); err == nil {
										monster.Skills[skillName] = modifier
									}
								}
							}
						}
					}

					// Parse senses
					if sensesMap, ok := monsterData["senses"].(map[string]interface{}); ok {
						monster.Senses = make(map[string]interface{})
						for sense, value := range sensesMap {
							monster.Senses[sense] = value
						}
					}

					// Parse languages (can be string or array)
					if langs, ok := monsterData["languages"].([]interface{}); ok {
						for _, lang := range langs {
							if langStr, ok := lang.(string); ok {
								monster.Languages = append(monster.Languages, langStr)
							}
						}
					} else if langStr, ok := monsterData["languages"].(string); ok && langStr != "" {
						monster.Languages = strings.Split(langStr, ", ")
					}

					if cr, ok := monsterData["cr"].(string); ok {
						monster.Challenge = cr
					}

					// Parse passive perception
					if passive, ok := monsterData["passive"].(string); ok {
						if monster.Senses == nil {
							monster.Senses = make(map[string]interface{})
						}
						monster.Senses["passive Perception"] = passive
					}

					// Parse abilities and actions
					if traits, ok := monsterData["trait"].([]interface{}); ok {
						for _, trait := range traits {
							if traitMap, ok := trait.(map[string]interface{}); ok {
								ability := MonsterAbility{
									Name: getStringFromMap(traitMap, "name"),
								}
								if textArray, ok := traitMap["text"].([]interface{}); ok {
									// Join text parts more intelligently
									textParts := convertToStringSlice(textArray)
									ability.Description = strings.Join(textParts, " ")
								}
								monster.SpecialAbilities = append(monster.SpecialAbilities, ability)
							}
						}
					}

					if actions, ok := monsterData["action"].([]interface{}); ok {
						for _, action := range actions {
							if actionMap, ok := action.(map[string]interface{}); ok {
								ability := MonsterAbility{
									Name: getStringFromMap(actionMap, "name"),
								}
								if textArray, ok := actionMap["text"].([]interface{}); ok {
									// Join text parts more intelligently
									textParts := convertToStringSlice(textArray)
									ability.Description = strings.Join(textParts, " ")
								}
								if attackArray, ok := actionMap["attack"].([]interface{}); ok {
									for _, attack := range attackArray {
										if attackStr, ok := attack.(string); ok {
											ability.Hit = attackStr
										}
									}
								}
								monster.Actions = append(monster.Actions, ability)
							}
						}
					}

					// Only set description if we actually extracted meaningful content
					// Only create fallback description if we really have nothing
					if description == "" && len(monster.SpecialAbilities) == 0 && len(monster.Actions) == 0 {
						// Create a basic description from available data
						var basicDesc []string
						if monster.Size != "" && monster.Type != "" {
							basicDesc = append(basicDesc, fmt.Sprintf("A %s %s.", monster.Size, monster.Type))
						}
						if monster.ArmorClass > 0 {
							basicDesc = append(basicDesc, fmt.Sprintf("Armor Class %d.", monster.ArmorClass))
						}
						if monster.HitPoints > 0 {
							basicDesc = append(basicDesc, fmt.Sprintf("Hit Points %d.", monster.HitPoints))
						}
						if monster.Speed != "" {
							basicDesc = append(basicDesc, fmt.Sprintf("Speed %s.", monster.Speed))
						}
						description = strings.Join(basicDesc, " ")
					}

					monster.Description = description
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

// getStringFromMap safely extracts a string from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].([]interface{}); ok && len(val) > 0 {
		if str, ok := val[0].(string); ok {
			return str
		}
	}
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
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
					} else if textArray, ok := itemMap["text"].([]interface{}); ok && len(textArray) > 0 {
						// Convert text array to string description
						var textParts []string
						for _, text := range textArray {
							if textStr, ok := text.(string); ok && textStr != "" {
								textParts = append(textParts, textStr)
							}
						}
						description = strings.Join(textParts, "\n")
					}

					item := Item{
						Name:        name,
						Description: description,
					}

					// Add additional fields if available
					if typ, ok := itemMap["type"].(string); ok {
						item.Type = typ
					}
					if val, ok := itemMap["value"].(string); ok {
						item.Value = val
					}
					if wt, ok := itemMap["weight"].(string); ok {
						item.Weight = wt
					}
					if dmg, ok := itemMap["dmg1"].(string); ok {
						item.Damage = dmg
					}
					if dmgType, ok := itemMap["dmgType"].(string); ok {
						item.DamageType = dmgType
					}
					if prop, ok := itemMap["property"].(string); ok {
						item.Property = prop
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
						} else if textArray, ok := itemMap["text"].([]interface{}); ok && len(textArray) > 0 {
							// Convert text array to string description
							var textParts []string
							for _, text := range textArray {
								if textStr, ok := text.(string); ok && textStr != "" {
									textParts = append(textParts, textStr)
								}
							}
							description = strings.Join(textParts, "\n")
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
	var backgrounds []Background

	// Handle both array and object formats
	switch v := data.(type) {
	case []interface{}:
		// Direct array format
		for _, item := range v {
			if bgData, ok := item.(map[string]interface{}); ok {
				if name, ok := bgData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := bgData["description"].(string); ok {
						description = desc
					}

					background := Background{
						Name:        name,
						Description: description,
					}
					backgrounds = append(backgrounds, background)
				}
			}
		}
	case map[string]interface{}:
		// Object format - iterate through the map
		for _, value := range v {
			if bgData, ok := value.(map[string]interface{}); ok {
				if name, ok := bgData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := bgData["description"].(string); ok {
						description = desc
					}

					background := Background{
						Name:        name,
						Description: description,
					}
					backgrounds = append(backgrounds, background)
				}
			}
		}
	}

	// Fallback to hardcoded list if no data found
	if len(backgrounds) == 0 {
		background_list := []string{
			"Acolyte", "Charlatan", "Criminal", "Entertainer", "Folk Hero",
			"Guild Artisan", "Hermit", "Noble", "Outlander", "Sage",
			"Sailor", "Soldier", "Urchin",
		}

		for _, name := range background_list {
			backgrounds = append(backgrounds, Background{
				Name:        name,
				Description: "A D&D 5e background",
			})
		}
	}

	return backgrounds
}

// parseClassesFromJSON creates classes from JSON data
func parseClassesFromJSON(data interface{}) []Class {
	var classes []Class

	// Handle both array and object formats
	switch v := data.(type) {
	case []interface{}:
		// Direct array format
		for _, item := range v {
			if classData, ok := item.(map[string]interface{}); ok {
				if name, ok := classData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := classData["description"].(string); ok {
						description = desc
					}

					class := Class{
						Name:        name,
						Description: description,
					}
					classes = append(classes, class)
				}
			}
		}
	case map[string]interface{}:
		// Object format - iterate through the map
		for _, value := range v {
			if classData, ok := value.(map[string]interface{}); ok {
				if name, ok := classData["name"].(string); ok && name != "" {
					description := ""
					if desc, ok := classData["description"].(string); ok {
						description = desc
					}

					class := Class{
						Name:        name,
						Description: description,
					}
					classes = append(classes, class)
				}
			}
		}
	}

	// Fallback to hardcoded list if no data found
	if len(classes) == 0 {
		class_list := []string{
			"Artificer", "Barbarian", "Bard", "Cleric", "Druid",
			"Fighter", "Monk", "Paladin", "Ranger", "Rogue",
			"Sorcerer", "Warlock", "Wizard",
		}

		for _, name := range class_list {
			classes = append(classes, Class{
				Name:        name,
				Description: "A D&D 5e class",
			})
		}
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
		AllMonsters = processStructuredMonsters(StructuredData.Monsters)
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
// Searches both structured and legacy data
func GetItemByName(name string) (*Item, error) {
	lowerName := strings.ToLower(name)

	// First search structured data
	for _, item := range AllItems {
		if strings.ToLower(item.Name) == lowerName {
			return &item, nil
		}
	}

	// Then search legacy data
	for _, item := range AllLegacyItems {
		if strings.ToLower(item.Name) == lowerName {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("item '%s' not found", name)
}

// GetSpeciesByName searches for a species by its name (case-insensitive)
// Searches both structured and legacy data
func GetSpeciesByName(name string) (*Species, error) {
	lowerName := strings.ToLower(name)

	// First search structured data
	for _, species := range AllSpecies {
		if strings.ToLower(species.Name) == lowerName {
			return &species, nil
		}
	}

	// Then search legacy data
	for _, species := range AllLegacySpecies {
		if strings.ToLower(species.Name) == lowerName {
			return &species, nil
		}
	}

	return nil, fmt.Errorf("species '%s' not found", name)
}

// GetBackgroundByName searches for a background by its name (case-insensitive)
// Searches both structured and legacy data
func GetBackgroundByName(name string) (*Background, error) {
	lowerName := strings.ToLower(name)

	// First search structured data
	for _, background := range AllBackgrounds {
		if strings.ToLower(background.Name) == lowerName {
			return &background, nil
		}
	}

	// Then search legacy data
	for _, background := range AllLegacyBackgrounds {
		if strings.ToLower(background.Name) == lowerName {
			return &background, nil
		}
	}

	return nil, fmt.Errorf("background '%s' not found", name)
}

// GetClassByName searches for a class by its name (case-insensitive)
// Searches both structured and legacy data
func GetClassByName(name string) (*Class, error) {
	lowerName := strings.ToLower(name)

	// First search structured data
	for _, class := range AllClasses {
		if strings.ToLower(class.Name) == lowerName {
			return &class, nil
		}
	}

	// Then search legacy data
	for _, class := range AllLegacyClasses {
		if strings.ToLower(class.Name) == lowerName {
			return &class, nil
		}
	}

	return nil, fmt.Errorf("class '%s' not found", name)
}
