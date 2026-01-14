package data

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// XMLMappings provides helper functions for XML format conversions
type XMLMappings struct {
	schoolMap   map[string]string
	damageMap   map[string]string
	propertyMap map[string]string
	itemTypeMap map[string]string
	sizeMap     map[string]string
}

// NewXMLMappings creates a new XMLMappings instance
func NewXMLMappings() *XMLMappings {
	return &XMLMappings{
		schoolMap: map[string]string{
			"A": "Abjuration",
			"V": "Evocation",
			"E": "Evocation",
			"I": "Illusion",
			"D": "Divination",
			"N": "Necromancy",
			"T": "Transmutation",
			"C": "Conjuration",
			"P": "Enchantment",
		},
		damageMap: map[string]string{
			"B":  "Bludgeoning",
			"P":  "Piercing",
			"S":  "Slashing",
			"N":  "Necrotic",
			"R":  "Radiant",
			"F":  "Fire",
			"C":  "Cold",
			"L":  "Lightning",
			"T":  "Thunder",
			"A":  "Acid",
			"O":  "Poison",
			"PS": "Psychic",
			"FO": "Force",
		},
		propertyMap: map[string]string{
			"V":   "Versatile",
			"T":   "Thrown",
			"R":   "Reach",
			"A":   "Ammunition",
			"F":   "Finesse",
			"L":   "Light",
			"H":   "Heavy",
			"2":   "Two-handed",
			"LD":  "Loading",
			"AD":  "Ammunition",
			"FIN": "Finesse",
			"LGT": "Light",
			"HVY": "Heavy",
			"THN": "Two-handed",
		},
		itemTypeMap: map[string]string{
			"M":  "Martial Weapon",
			"S":  "Simple Weapon",
			"R":  "Ranged Weapon",
			"A":  "Armor",
			"G":  "Adventuring Gear",
			"W":  "Wondrous Item",
			"P":  "Potion",
			"SC": "Scroll",
			"RG": "Ring",
			"RD": "Rod",
			"ST": "Staff",
			"WD": "Wand",
		},
		sizeMap: map[string]string{
			"T": "Tiny",
			"S": "Small",
			"M": "Medium",
			"L": "Large",
			"H": "Huge",
			"G": "Gargantuan",
		},
	}
}

// GetSchoolName converts school abbreviation to full name
func (xm *XMLMappings) GetSchoolName(code string) string {
	if name, exists := xm.schoolMap[code]; exists {
		return name
	}
	return code
}

// GetDamageTypeName converts damage type abbreviation to full name
func (xm *XMLMappings) GetDamageTypeName(code string) string {
	if name, exists := xm.damageMap[code]; exists {
		return name
	}
	return code
}

// GetPropertyDescription converts property abbreviation to description
func (xm *XMLMappings) GetPropertyDescription(code string) string {
	if desc, exists := xm.propertyMap[code]; exists {
		return desc
	}
	return code
}

// GetItemTypeName converts item type abbreviation to full name
func (xm *XMLMappings) GetItemTypeName(code string) string {
	if name, exists := xm.itemTypeMap[code]; exists {
		return name
	}
	return code
}

// GetSizeName converts size abbreviation to full name
func (xm *XMLMappings) GetSizeName(code string) string {
	if name, exists := xm.sizeMap[code]; exists {
		return name
	}
	return code
}

// ExpandProperties expands comma-separated property codes to descriptions
func (xm *XMLMappings) ExpandProperties(properties string) []string {
	if properties == "" {
		return []string{}
	}

	result := []string{}
	for _, prop := range strings.Split(properties, ",") {
		prop = strings.TrimSpace(prop)
		if prop != "" {
			result = append(result, xm.GetPropertyDescription(prop))
		}
	}
	return result
}

// ParseAbilityScores parses ability score values from XML strings
func ParseAbilityScores(str, dex, con, intStr, wis, cha string) map[string]int {
	stats := make(map[string]int)

	stats["STR"] = parseIntOrZero(str)
	stats["DEX"] = parseIntOrZero(dex)
	stats["CON"] = parseIntOrZero(con)
	stats["INT"] = parseIntOrZero(intStr)
	stats["WIS"] = parseIntOrZero(wis)
	stats["CHA"] = parseIntOrZero(cha)

	return stats
}

// parseIntOrZero safely parses an integer string, returning 0 on error
func parseIntOrZero(s string) int {
	s = strings.TrimSpace(s)
	if val, err := parseInt(s); err == nil {
		return val
	}
	return 10 // Default ability score
}

// parseInt safely parses an integer
func parseInt(s string) (int, error) {
	// Handle strings like "+3" or "-2"
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	if val, err := parseInt(s); err == nil {
		return val, nil
	}
	return 0, fmt.Errorf("invalid integer: %s", s)
}

// ParseHitDice extracts hit dice information from HP string
func ParseHitDice(hpStr string) string {
	// Extract dice formula from HP string like "135 (18d10+36)"
	if strings.Contains(hpStr, "(") && strings.Contains(hpStr, ")") {
		start := strings.Index(hpStr, "(") + 1
		end := strings.Index(hpStr, ")")
		if start > 0 && end > start {
			return hpStr[start:end]
		}
	}
	return ""
}

// ParseArmorClass extracts AC value and description
func ParseArmorClass(acStr string) (int, string) {
	parts := strings.SplitN(acStr, " ", 2)
	ac := 0
	if len(parts) > 0 {
		ac = parseIntOrZero(parts[0])
	}

	description := ""
	if len(parts) > 1 {
		description = parts[1]
	}

	return ac, description
}

// ParseChallengeRating extracts CR and calculates XP
func ParseChallengeRating(crStr string) (string, int) {
	// CR to XP mapping
	crXpMap := map[string]int{
		"0":   0,
		"1/8": 25,
		"1/4": 50,
		"1/2": 100,
		"1":   200,
		"2":   450,
		"3":   700,
		"4":   1100,
		"5":   1800,
		"6":   2300,
		"7":   2900,
		"8":   3900,
		"9":   5000,
		"10":  5900,
		"11":  7200,
		"12":  8400,
		"13":  10000,
		"14":  11500,
		"15":  13000,
		"16":  15000,
		"17":  18000,
		"18":  20000,
		"19":  22000,
		"20":  25000,
		"21":  33000,
		"22":  41000,
		"23":  50000,
		"24":  62000,
		"25":  75000,
		"26":  90000,
		"27":  105000,
		"28":  120000,
		"29":  135000,
		"30":  155000,
	}

	cr := strings.TrimSpace(crStr)
	xp := 0
	if xpVal, exists := crXpMap[cr]; exists {
		xp = xpVal
	}

	return cr, xp
}

// ParseSavingThrows parses comma-separated saving throw list
func ParseSavingThrows(savingThrows string) []string {
	if savingThrows == "" {
		return []string{}
	}

	throws := strings.Split(savingThrows, ",")
	for i, st := range throws {
		throws[i] = strings.TrimSpace(st)
	}

	return throws
}

// ParseSkills parses skill bonus string into map
func ParseSkills(skillsStr string) map[string]int {
	skills := make(map[string]int)

	if skillsStr == "" {
		return skills
	}

	// Parse skills like "Perception +3, Stealth +6"
	for _, skill := range strings.Split(skillsStr, ",") {
		skill = strings.TrimSpace(skill)
		parts := strings.Split(skill, " ")
		if len(parts) >= 2 {
			name := strings.Join(parts[:len(parts)-1], " ")
			bonus := parseIntOrZero(parts[len(parts)-1])
			if name != "" {
				skills[name] = bonus
			}
		}
	}

	return skills
}

// ParseSenses parses senses string into map
func ParseSenses(sensesStr, passiveStr string) map[string]interface{} {
	senses := make(map[string]interface{})

	if sensesStr != "" {
		// Parse senses like "darkvision 60 ft., passive Perception 12"
		for _, sense := range strings.Split(sensesStr, ",") {
			sense = strings.TrimSpace(sense)
			if sense != "" {
				senses[sense] = true
			}
		}
	}

	if passiveStr != "" {
		senses["Passive"] = passiveStr
	}

	return senses
}

// ParseLanguages parses comma-separated languages
func ParseLanguages(languagesStr string) []string {
	if languagesStr == "" {
		return []string{}
	}

	languages := strings.Split(languagesStr, ",")
	for i, lang := range languages {
		languages[i] = strings.TrimSpace(lang)
	}

	return languages
}

// GetCompendiumFileList returns sorted list of available compendium files
func GetCompendiumFileList(compendiumDir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(compendiumDir, "*.xml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list compendium files: %w", err)
	}

	var names []string
	for _, file := range files {
		name := filepath.Base(file)
		if strings.HasSuffix(name, ".xml") {
			names = append(names, name)
		}
	}

	// Sort alphabetically
	sort.Strings(names)

	return names, nil
}

// ValidateCompendiumFile checks if a file exists and is readable
func ValidateCompendiumFile(compendiumDir, filename string) error {
	if !strings.HasSuffix(filename, ".xml") {
		return fmt.Errorf("file must be an XML file")
	}

	fullPath := filepath.Join(compendiumDir, filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("compendium file not found: %s", filename)
	}

	return nil
}

// GetDefaultCompendiumFile returns the recommended default file
func GetDefaultCompendiumFile() string {
	return "Complete_Compendium.xml"
}

// FilterCompendiumFiles filters files by category
func FilterCompendiumFiles(files []string, category string) []string {
	var filtered []string

	category = strings.ToLower(category)
	for _, file := range files {
		lowerFile := strings.ToLower(file)

		switch category {
		case "core":
			if strings.Contains(lowerFile, "core") || strings.Contains(lowerFile, "rulebook") {
				filtered = append(filtered, file)
			}
		case "wotc":
			if strings.Contains(lowerFile, "wotc") && !strings.Contains(lowerFile, "partner") {
				filtered = append(filtered, file)
			}
		case "complete":
			if strings.Contains(lowerFile, "complete") {
				filtered = append(filtered, file)
			}
		default:
			// Return all files if no specific filter
			filtered = files
		}
	}

	return filtered
}
