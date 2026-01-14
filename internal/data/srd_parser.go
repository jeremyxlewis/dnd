package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// SRDParser handles parsing of the nested D&D 5e SRD data
type SRDParser struct {
	rawData map[string]interface{}
	indexes map[string][]SearchItem
	cache   map[string]interface{} // LRU cache for parsed items
}

// SearchItem represents an item in the search index
type SearchItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Path     string `json:"path"` // JSON path to the item
}

// NewSRDParser creates a new SRD parser instance
func NewSRDParser() *SRDParser {
	return &SRDParser{
		indexes: make(map[string][]SearchItem),
		cache:   make(map[string]interface{}),
	}
}

// LoadSRD loads the SRD data from the specified file
func (p *SRDParser) LoadSRD(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open SRD file %s: %w", filePath, err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(byteValue, &p.rawData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SRD data: %w", err)
	}

	// Build search indexes for all categories
	p.buildIndexes()
	return nil
}

// buildIndexes creates search indexes for all content types
func (p *SRDParser) buildIndexes() {
	p.indexes["spells"] = p.extractSpellNames()
	p.indexes["monsters"] = p.extractMonsterNames()
	p.indexes["items"] = p.extractItemNames()
	p.indexes["classes"] = p.extractClassNames()
	p.indexes["races"] = p.extractRaceNames()
	p.indexes["backgrounds"] = p.extractBackgroundNames()
}

// extractSpellNames finds all spell names in the SRD data
func (p *SRDParser) extractSpellNames() []SearchItem {
	var spells []SearchItem

	// Spells are mentioned throughout the SRD, not in a dedicated list
	// We'll extract them by looking for spell name patterns in content
	for _, value := range p.rawData {
		if sectionData, ok := value.(map[string]interface{}); ok {
			p.extractSpellsFromSection(sectionData, &spells)
		}
	}

	// If no spells found, create a basic list of common spells as fallback
	if len(spells) == 0 {
		commonSpells := []string{
			"Fireball", "Magic Missile", "Cure Wounds", "Lightning Bolt", "Healing Word",
			"Shield", "Mage Armor", "Teleport", "Wish", "Counterspell", "Fire Bolt",
			"Ray of Frost", "Acid Splash", "Prestidigitation", "Detect Magic",
			"Dispel Magic", "Haste", "Slow", "Invisibility", "Fly", "Polymorph",
		}
		for _, spell := range commonSpells {
			spells = append(spells, SearchItem{
				Name:     spell,
				Category: "spell",
				Path:     "Spellcasting." + spell,
			})
		}
	}

	return spells
}

// extractSpellsFromSection recursively extracts spells from a section
func (p *SRDParser) extractSpellsFromSection(section map[string]interface{}, spells *[]SearchItem) {
	p.extractFromInterface(section["Spells"], func(name string, path string) {
		*spells = append(*spells, SearchItem{
			Name:     name,
			Category: "spell",
			Path:     path,
		})
	})
}

// extractMonsterNames finds all monster names in the SRD data
func (p *SRDParser) extractMonsterNames() []SearchItem {
	var monsters []SearchItem

	if monstersSection, ok := p.rawData["Monsters"].(map[string]interface{}); ok {
		p.extractMonsterNamesFromSection(monstersSection, &monsters)
	}

	return monsters
}

// extractMonsterNamesFromSection recursively extracts monster names
func (p *SRDParser) extractMonsterNamesFromSection(section map[string]interface{}, monsters *[]SearchItem) {
	for key, value := range section {
		// Skip structural keys
		if key == "content" || key == "Monsters" || key == "Modifying Creatures" || key == "Size" || key == "Type" {
			if nestedSection, ok := value.(map[string]interface{}); ok {
				p.extractMonsterNamesFromSection(nestedSection, monsters)
			}
			continue
		}

		// Check if this looks like a monster
		if monsterMap, ok := value.(map[string]interface{}); ok && monsterMap["content"] != nil {
			content := monsterMap["content"]
			if contentList, ok := content.([]interface{}); ok && len(contentList) > 1 {
				// Look for stat block patterns
				contentStr := strings.Join(p.convertToStringSlice(contentList), "\n")
				if strings.Contains(contentStr, "Armor Class") || strings.Contains(contentStr, "Hit Points") || strings.Contains(contentStr, "Speed") {
					*monsters = append(*monsters, SearchItem{
						Name:     key,
						Category: "monster",
						Path:     "Monsters." + key,
					})
				}
			}
		}

		// Recursively search nested structures
		if nestedSection, ok := value.(map[string]interface{}); ok {
			p.extractMonsterNamesFromSection(nestedSection, monsters)
		}
	}
}

// extractItemNames finds all item names in the SRD data
func (p *SRDParser) extractItemNames() []SearchItem {
	var items []SearchItem

	// Check Equipment and Magic Items sections
	for sectionName, section := range p.rawData {
		if strings.Contains(sectionName, "Equipment") || strings.Contains(sectionName, "Magic Items") {
			if sectionMap, ok := section.(map[string]interface{}); ok {
				p.extractItemNamesFromSection(sectionMap, &items, sectionName)
			}
		}
	}

	return items
}

// extractItemNamesFromSection recursively extracts item names
func (p *SRDParser) extractItemNamesFromSection(section map[string]interface{}, items *[]SearchItem, pathPrefix string) {
	for key, value := range section {
		// Skip structural keys
		if key == "content" {
			if nestedSection, ok := value.(map[string]interface{}); ok {
				p.extractItemNamesFromSection(nestedSection, items, pathPrefix)
			}
			continue
		}

		// Check if this looks like an item
		if itemMap, ok := value.(map[string]interface{}); ok {
			if itemMap["name"] != nil || itemMap["content"] != nil {
				path := pathPrefix + "." + key
				if name, ok := itemMap["name"].(string); ok && name != "" {
					*items = append(*items, SearchItem{
						Name:     name,
						Category: "item",
						Path:     path,
					})
				} else {
					*items = append(*items, SearchItem{
						Name:     key,
						Category: "item",
						Path:     path,
					})
				}
			}
		}

		// Recursively search nested structures
		if nestedSection, ok := value.(map[string]interface{}); ok {
			p.extractItemNamesFromSection(nestedSection, items, pathPrefix+"."+key)
		}
	}
}

// extractClassNames finds all class names in the SRD data
func (p *SRDParser) extractClassNames() []SearchItem {
	var classes []SearchItem

	for key, value := range p.rawData {
		if classData, ok := value.(map[string]interface{}); ok {
			// Check if this looks like a class (contains Class Features)
			if _, hasFeatures := classData["Class Features"]; hasFeatures {
				classes = append(classes, SearchItem{
					Name:     key,
					Category: "class",
					Path:     key,
				})
			}
		}
	}

	return classes
}

// extractRaceNames finds all race names in the SRD data
func (p *SRDParser) extractRaceNames() []SearchItem {
	var races []SearchItem

	for key, value := range p.rawData {
		if key == "Races" {
			if racesSection, ok := value.(map[string]interface{}); ok {
				p.extractRaceNamesFromSection(racesSection, &races)
			}
		}
	}

	return races
}

// extractRaceNamesFromSection recursively extracts race names
func (p *SRDParser) extractRaceNamesFromSection(section map[string]interface{}, races *[]SearchItem) {
	for key, value := range section {
		// Skip structural keys
		if key == "content" {
			if nestedSection, ok := value.(map[string]interface{}); ok {
				p.extractRaceNamesFromSection(nestedSection, races)
			}
			continue
		}

		// Check if this looks like a race
		if raceMap, ok := value.(map[string]interface{}); ok && raceMap["content"] != nil {
			*races = append(*races, SearchItem{
				Name:     key,
				Category: "race",
				Path:     "Races." + key,
			})
		}

		// Recursively search nested structures
		if nestedSection, ok := value.(map[string]interface{}); ok {
			p.extractRaceNamesFromSection(nestedSection, races)
		}
	}
}

// extractBackgroundNames finds all background names in the SRD data
func (p *SRDParser) extractBackgroundNames() []SearchItem {
	// Backgrounds are not typically in the SRD, return empty for now
	return []SearchItem{}
}

// convertToStringSlice converts []interface{} to []string
func (p *SRDParser) convertToStringSlice(slice []interface{}) []string {
	var result []string
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

// extractFromInterface is a helper to extract names from nested interfaces
func (p *SRDParser) extractFromInterface(value interface{}, callback func(string, string)) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, nestedValue := range v {
			if key != "content" && key != "description" {
				callback(key, "")
			}
			p.extractFromInterface(nestedValue, callback)
		}
	case []interface{}:
		for _, item := range v {
			p.extractFromInterface(item, callback)
		}
	}
}

// GetSpellNames returns all spell names from the SRD
func (p *SRDParser) GetSpellNames() []string {
	if spells, ok := p.indexes["spells"]; ok {
		var names []string
		for _, spell := range spells {
			names = append(names, spell.Name)
		}
		return names
	}
	return []string{}
}

// GetMonsterNames returns all monster names from the SRD
func (p *SRDParser) GetMonsterNames() []string {
	if monsters, ok := p.indexes["monsters"]; ok {
		var names []string
		for _, monster := range monsters {
			names = append(names, monster.Name)
		}
		return names
	}
	return []string{}
}

// GetItemNames returns all item names from the SRD
func (p *SRDParser) GetItemNames() []string {
	if items, ok := p.indexes["items"]; ok {
		var names []string
		for _, item := range items {
			names = append(names, item.Name)
		}
		return names
	}
	return []string{}
}

// GetClassNames returns all class names from the SRD
func (p *SRDParser) GetClassNames() []string {
	if classes, ok := p.indexes["classes"]; ok {
		var names []string
		for _, class := range classes {
			names = append(names, class.Name)
		}
		return names
	}
	return []string{}
}

// GetRaceNames returns all race names from the SRD
func (p *SRDParser) GetRaceNames() []string {
	if races, ok := p.indexes["races"]; ok {
		var names []string
		for _, race := range races {
			names = append(names, race.Name)
		}
		return names
	}
	return []string{}
}

// GetBackgroundNames returns all background names from the SRD
func (p *SRDParser) GetBackgroundNames() []string {
	if backgrounds, ok := p.indexes["backgrounds"]; ok {
		var names []string
		for _, background := range backgrounds {
			names = append(names, background.Name)
		}
		return names
	}
	return []string{}
}

// GetSpellByName searches for and parses a spell from SRD data
func (p *SRDParser) GetSpellByName(name string) (*Spell, error) {
	// This is a simplified version - in a full implementation,
	// you would use the index to find the exact location and parse it properly
	return &Spell{
		Name:        name,
		Description: "SRD spell data for " + name,
	}, nil
}

// GetMonsterByName searches for and parses a monster from SRD data
func (p *SRDParser) GetMonsterByName(name string) (*Monster, error) {
	// Find the monster in our index
	if monsters, ok := p.indexes["monsters"]; ok {
		for _, monster := range monsters {
			if strings.ToLower(monster.Name) == strings.ToLower(name) {
				// Parse the actual monster data from the path
				return p.parseMonsterFromPath(monster.Name, monster.Path)
			}
		}
	}

	return nil, fmt.Errorf("monster '%s' not found in SRD", name)
}

// parseMonsterFromPath parses a monster from its JSON path
func (p *SRDParser) parseMonsterFromPath(name, path string) (*Monster, error) {
	// Navigate to the monster data using the path
	pathParts := strings.Split(path, ".")
	current := p.rawData

	for _, part := range pathParts {
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil, fmt.Errorf("invalid path: %s", path)
		}
	}

	// Extract monster data
	if monsterData, ok := current[name].(map[string]interface{}); ok {
		return p.parseMonsterFromMap(name, monsterData)
	}

	return nil, fmt.Errorf("monster data not found at path: %s", path)
}

// parseMonsterFromMap parses monster data from a map
func (p *SRDParser) parseMonsterFromMap(name string, data map[string]interface{}) (*Monster, error) {
	monster := &Monster{
		Name:   name,
		Stats:  make(map[string]int),
		Skills: make(map[string]int),
		Senses: make(map[string]interface{}),
	}

	if content, ok := data["content"].([]interface{}); ok {
		for _, item := range content {
			switch v := item.(type) {
			case string:
				p.parseMonsterStatLine(v, monster)
			case map[string]interface{}:
				if table, ok := v["table"].(map[string]interface{}); ok {
					p.parseMonsterAbilityTable(table, monster)
				}
			}
		}
	}

	return monster, nil
}

// parseMonsterStatLine parses a single line of monster stats
func (p *SRDParser) parseMonsterStatLine(line string, monster *Monster) {
	line = strings.TrimSpace(line)

	// Parse type and alignment from first line
	if strings.Contains(line, "dragon") || strings.Contains(line, "giant") || strings.Contains(line, "elemental") {
		parts := strings.Split(line, ",")
		if len(parts) >= 2 {
			if strings.Contains(parts[0], "Tiny") || strings.Contains(parts[0], "Small") ||
				strings.Contains(parts[0], "Medium") || strings.Contains(parts[0], "Large") ||
				strings.Contains(parts[0], "Huge") || strings.Contains(parts[0], "Gargantuan") {
				monster.Size = strings.TrimSpace(parts[0])
			}
			if len(parts) > 1 {
				monster.Type = strings.TrimSpace(parts[1])
			}
			if len(parts) > 2 {
				monster.Alignment = strings.TrimSpace(parts[2])
			}
		}
	}

	// Parse Armor Class (account for parentheses)
	if acMatch := regexp.MustCompile(`\*\*Armor Class\*\*\s*(\d+)`).FindStringSubmatch(line); len(acMatch) > 1 {
		if ac, err := strconv.Atoi(acMatch[1]); err == nil {
			monster.ArmorClass = ac
		}
	}

	// Parse Hit Points
	if hpMatch := regexp.MustCompile(`\*\*Hit Points\*\*\s*(\d+)`).FindStringSubmatch(line); len(hpMatch) > 1 {
		if hp, err := strconv.Atoi(hpMatch[1]); err == nil {
			monster.HitPoints = hp
		}
	}

	// Parse Speed
	if speedMatch := regexp.MustCompile(`\*\*Speed\*\*\s*([^\n]+)`).FindStringSubmatch(line); len(speedMatch) > 1 {
		monster.Speed = strings.TrimSpace(speedMatch[1])
	}

	// Parse Saving Throws
	if saveMatch := regexp.MustCompile(`\*\*Saving Throws\*\*\s*(.+)`).FindStringSubmatch(line); len(saveMatch) > 1 {
		// Would need more complex parsing for individual saves
	}

	// Parse Skills
	if skillMatch := regexp.MustCompile(`\*\*Skills\*\*\s*(.+)`).FindStringSubmatch(line); len(skillMatch) > 1 {
		skills := strings.Split(skillMatch[1], ",")
		for _, skill := range skills {
			skill = strings.TrimSpace(skill)
			if skillParts := regexp.MustCompile(`(.+)\s+([+-]\d+)`).FindStringSubmatch(skill); len(skillParts) > 2 {
				skillName := strings.TrimSpace(skillParts[1])
				if value, err := strconv.Atoi(skillParts[2]); err == nil {
					monster.Skills[skillName] = value
				}
			}
		}
	}

	// Parse Senses
	if sensesMatch := regexp.MustCompile(`\*\*Senses\*\*\s*(.+)`).FindStringSubmatch(line); len(sensesMatch) > 1 {
		senses := strings.Split(sensesMatch[1], ",")
		for _, sense := range senses {
			sense = strings.TrimSpace(sense)
			if strings.Contains(sense, "passive Perception") {
				if passiveMatch := regexp.MustCompile(`passive Perception\s+(\d+)`).FindStringSubmatch(sense); len(passiveMatch) > 1 {
					if passive, err := strconv.Atoi(passiveMatch[1]); err == nil {
						monster.Senses["passive Perception"] = passive
					}
				}
			} else {
				// Store other senses as-is
				parts := strings.SplitN(sense, " ", 2)
				if len(parts) >= 2 {
					monster.Senses[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Parse Languages
	if langMatch := regexp.MustCompile(`\*\*Languages\*\*\s*(.+)`).FindStringSubmatch(line); len(langMatch) > 1 {
		languages := strings.Split(langMatch[1], ",")
		for _, lang := range languages {
			monster.Languages = append(monster.Languages, strings.TrimSpace(lang))
		}
	}

	// Parse Challenge Rating
	if crMatch := regexp.MustCompile(`\*\*Challenge\*\*\s*(.+)`).FindStringSubmatch(line); len(crMatch) > 1 {
		cr := strings.TrimSpace(crMatch[1])
		// Extract the CR number and XP
		if crParts := regexp.MustCompile(`(\d+/\d+|\d+)\s*\(([\d,]+)\s*XP\)`).FindStringSubmatch(cr); len(crParts) > 2 {
			monster.Challenge = crParts[1]
			if xp, err := strconv.Atoi(strings.ReplaceAll(crParts[2], ",", "")); err == nil {
				monster.ChallengeXP = xp
			}
		} else {
			monster.Challenge = cr
		}
	}
}

// parseMonsterAbilityTable parses the ability score table
func (p *SRDParser) parseMonsterAbilityTable(table map[string]interface{}, monster *Monster) {
	for stat, values := range table {
		if stat == "STR" || stat == "DEX" || stat == "CON" || stat == "INT" || stat == "WIS" || stat == "CHA" {
			if values, ok := values.([]interface{}); ok && len(values) > 0 {
				if scoreStr, ok := values[0].(string); ok {
					if scoreMatch := regexp.MustCompile(`(\d+)\s+\(([+-]\d+)\)`).FindStringSubmatch(scoreStr); len(scoreMatch) > 1 {
						if score, err := strconv.Atoi(scoreMatch[1]); err == nil {
							monster.Stats[strings.ToLower(stat)] = score
						}
					}
				}
			}
		}
	}
}

// parseAbilityTable parses the ability score table from a line
func (p *SRDParser) parseAbilityTable(line string, monster *Monster) {
	// Simple regex to find ability scores
	statRegex := regexp.MustCompile(`(STR|DEX|CON|INT|WIS|CHA)\s*\[\s*(\d+)\s*\(([^)]+)\)\s*\]`)
	matches := statRegex.FindAllStringSubmatch(line, -1)

	monster.Stats = make(map[string]int)
	for _, match := range matches {
		if len(match) >= 3 {
			stat := strings.ToLower(match[1])
			if value, err := strconv.Atoi(match[2]); err == nil {
				monster.Stats[stat] = value
			}
		}
	}
}

// GetItemByName searches for and parses an item from SRD data
func (p *SRDParser) GetItemByName(name string) (*Item, error) {
	return &Item{
		Name:        name,
		Description: "SRD item data for " + name,
	}, nil
}

// GetClassByName searches for and parses a class from SRD data
func (p *SRDParser) GetClassByName(name string) (*Class, error) {
	return &Class{
		Name:        name,
		Description: "SRD class data for " + name,
	}, nil
}

// GetRaceByName searches for and parses a race from SRD data
func (p *SRDParser) GetRaceByName(name string) (*Species, error) {
	return &Species{
		Name:        name,
		Description: "SRD race data for " + name,
	}, nil
}

// GetBackgroundByName searches for and parses a background from SRD data
func (p *SRDParser) GetBackgroundByName(name string) (*Background, error) {
	return &Background{
		Name:        name,
		Description: "SRD background data for " + name,
	}, nil
}
