package data

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CompendiumParser handles parsing of Compendium XML files
type CompendiumParser struct {
	spells      []Spell
	monsters    []Monster
	items       []Item
	species     []Species
	backgrounds []Background
	classes     []Class

	// Mapping tables for format conversion
	schoolMap   map[string]string
	damageMap   map[string]string
	propertyMap map[string]string
	itemTypeMap map[string]string
}

// CompendiumXML represents the root structure of a Compendium XML file
type CompendiumXML struct {
	XMLName     xml.Name        `xml:"compendium"`
	Version     string          `xml:"version,attr"`
	Spells      []XMLSpell      `xml:"spell"`
	Monsters    []XMLMonster    `xml:"monster"`
	Items       []XMLItem       `xml:"item"`
	Races       []XMLRace       `xml:"race"`
	Backgrounds []XMLBackground `xml:"background"`
	Classes     []XMLClass      `xml:"class"`
}

// XML structures for parsing
type XMLSpell struct {
	XMLName    xml.Name `xml:"spell"`
	Name       string   `xml:"name"`
	Level      string   `xml:"level"`
	School     string   `xml:"school"`
	Time       string   `xml:"time"`
	Range      string   `xml:"range"`
	Components string   `xml:"components"`
	Duration   string   `xml:"duration"`
	Classes    string   `xml:"classes"`
	Text       []string `xml:"text"`
}

type XMLMonster struct {
	XMLName   xml.Name    `xml:"monster"`
	Name      string      `xml:"name"`
	Size      string      `xml:"size"`
	Type      string      `xml:"type"`
	Alignment string      `xml:"alignment"`
	AC        string      `xml:"ac"`
	HP        string      `xml:"hp"`
	Speed     string      `xml:"speed"`
	Str       string      `xml:"str"`
	Dex       string      `xml:"dex"`
	Con       string      `xml:"con"`
	Int       string      `xml:"int"`
	Wis       string      `xml:"wis"`
	Cha       string      `xml:"cha"`
	Skills    string      `xml:"skills"`
	Senses    string      `xml:"senses"`
	Passive   string      `xml:"passive"`
	Languages string      `xml:"languages"`
	CR        string      `xml:"cr"`
	Traits    []XMLTrait  `xml:"trait"`
	Actions   []XMLAction `xml:"action"`
}

type XMLTrait struct {
	Name string `xml:"name"`
	Text string `xml:"text"`
}

type XMLAction struct {
	Name   string `xml:"name"`
	Text   string `xml:"text"`
	Attack string `xml:"attack"`
}

type XMLItem struct {
	XMLName  xml.Name `xml:"item"`
	Name     string   `xml:"name"`
	Type     string   `xml:"type"`
	Value    string   `xml:"value"`
	Weight   string   `xml:"weight"`
	Dmg1     string   `xml:"dmg1"`
	DmgType  string   `xml:"dmgType"`
	Property string   `xml:"property"`
	Text     []string `xml:"text"`
}

type XMLRace struct {
	XMLName xml.Name   `xml:"race"`
	Name    string     `xml:"name"`
	Size    string     `xml:"size"`
	Speed   string     `xml:"speed"`
	Ability string     `xml:"ability"`
	Traits  []XMLTrait `xml:"trait"`
}

type XMLBackground struct {
	XMLName     xml.Name   `xml:"background"`
	Name        string     `xml:"name"`
	Proficiency string     `xml:"proficiency"`
	Feature     string     `xml:"feature"`
	Text        []string   `xml:"text"`
	Traits      []XMLTrait `xml:"trait"`
}

type XMLClass struct {
	XMLName      xml.Name `xml:"class"`
	Name         string   `xml:"name"`
	HD           string   `xml:"hd"`
	Proficiency  string   `xml:"proficiency"`
	SpellAbility string   `xml:"spellAbility"`
	Skills       string   `xml:"numSkills"`
	Armor        string   `xml:"armor"`
	Weapons      string   `xml:"weapons"`
	Tools        string   `xml:"tools"`
	Wealth       string   `xml:"wealth"`
}

// NewCompendiumParser creates a new CompendiumParser instance
func NewCompendiumParser() *CompendiumParser {
	return &CompendiumParser{
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
			"V":  "Versatile",
			"T":  "Thrown",
			"R":  "Reach",
			"A":  "Ammunition",
			"F":  "Finesse",
			"L":  "Light",
			"H":  "Heavy",
			"2":  "Two-handed",
			"LD": "Loading",
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
	}
}

// LoadCompendium loads a compendium XML file
func (cp *CompendiumParser) LoadCompendium(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read compendium file: %w", err)
	}

	var compendium CompendiumXML
	err = xml.Unmarshal(data, &compendium)
	if err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}

	// Parse spells
	for _, xmlSpell := range compendium.Spells {
		spell := cp.parseSpell(xmlSpell)
		cp.spells = append(cp.spells, spell)
	}

	// Parse monsters
	for _, xmlMonster := range compendium.Monsters {
		monster := cp.parseMonster(xmlMonster)
		cp.monsters = append(cp.monsters, monster)
	}

	// Parse items
	for _, xmlItem := range compendium.Items {
		item := cp.parseItem(xmlItem)
		cp.items = append(cp.items, item)
	}

	// Parse races
	for _, xmlRace := range compendium.Races {
		species := cp.parseRace(xmlRace)
		cp.species = append(cp.species, species)
	}

	// Parse backgrounds
	for _, xmlBackground := range compendium.Backgrounds {
		background := cp.parseBackground(xmlBackground)
		cp.backgrounds = append(cp.backgrounds, background)
	}

	// Parse classes
	for _, xmlClass := range compendium.Classes {
		class := cp.parseClass(xmlClass)
		cp.classes = append(cp.classes, class)
	}

	return nil
}

// parseSpell converts XMLSpell to Spell
func (cp *CompendiumParser) parseSpell(xmlSpell XMLSpell) Spell {
	level, _ := strconv.Atoi(xmlSpell.Level)
	school := cp.schoolMap[xmlSpell.School]
	if school == "" {
		school = xmlSpell.School
	}

	components := cp.parseComponents(xmlSpell.Components)

	description := ""
	higherLevels := ""
	if len(xmlSpell.Text) > 0 {
		description = xmlSpell.Text[0]
		if len(xmlSpell.Text) > 1 {
			higherLevels = xmlSpell.Text[1]
		}
	}

	classes := cp.parseClasses(xmlSpell.Classes)

	return Spell{
		Name:         xmlSpell.Name,
		Level:        level,
		School:       school,
		CastingTime:  xmlSpell.Time,
		Range:        xmlSpell.Range,
		Components:   components,
		Duration:     xmlSpell.Duration,
		Description:  description,
		HigherLevels: &higherLevels,
		Classes:      classes,
	}
}

// parseMonster converts XMLMonster to Monster
func (cp *CompendiumParser) parseMonster(xmlMonster XMLMonster) Monster {
	ac, acStr := cp.parseArmorClass(xmlMonster.AC)
	hp, hpStr := cp.parseHitPoints(xmlMonster.HP)

	stats := map[string]int{
		"STR": cp.parseAbility(xmlMonster.Str),
		"DEX": cp.parseAbility(xmlMonster.Dex),
		"CON": cp.parseAbility(xmlMonster.Con),
		"INT": cp.parseAbility(xmlMonster.Int),
		"WIS": cp.parseAbility(xmlMonster.Wis),
		"CHA": cp.parseAbility(xmlMonster.Cha),
	}

	skills := cp.parseSkills(xmlMonster.Skills)
	senses := cp.parseSenses(xmlMonster.Senses, xmlMonster.Passive)
	languages := cp.parseLanguages(xmlMonster.Languages)

	traits := cp.parseTraits(xmlMonster.Traits)
	actions := cp.parseActions(xmlMonster.Actions)

	return Monster{
		Name:             xmlMonster.Name,
		Size:             xmlMonster.Size,
		Type:             xmlMonster.Type,
		Alignment:        xmlMonster.Alignment,
		ArmorClass:       ac,
		ArmorClassStr:    acStr,
		HitPoints:        hp,
		HitPointsStr:     hpStr,
		Speed:            xmlMonster.Speed,
		Stats:            stats,
		Skills:           skills,
		Senses:           senses,
		Languages:        languages,
		Challenge:        xmlMonster.CR,
		SpecialAbilities: traits,
		Actions:          actions,
	}
}

// parseItem converts XMLItem to Item
func (cp *CompendiumParser) parseItem(xmlItem XMLItem) Item {
	description := ""
	if len(xmlItem.Text) > 0 {
		description = strings.Join(xmlItem.Text, " ")
	}

	damageType := cp.damageMap[xmlItem.DmgType]
	if damageType == "" {
		damageType = xmlItem.DmgType
	}

	return Item{
		Name:        xmlItem.Name,
		Description: description,
		Type:        xmlItem.Type,
		Value:       xmlItem.Value,
		Weight:      xmlItem.Weight,
		Damage:      xmlItem.Dmg1,
		DamageType:  damageType,
		Property:    xmlItem.Property,
	}
}

// parseRace converts XMLRace to Species
func (cp *CompendiumParser) parseRace(xmlRace XMLRace) Species {
	abilityBonuses := cp.parseAbilityBonus(xmlRace.Ability)

	traits := []string{}
	for _, trait := range xmlRace.Traits {
		traits = append(traits, trait.Text)
	}

	return Species{
		Name:           xmlRace.Name,
		AbilityBonuses: abilityBonuses,
		Speed:          xmlRace.Speed,
		Traits:         traits,
	}
}

// parseBackground converts XMLBackground to Background
func (cp *CompendiumParser) parseBackground(xmlBackground XMLBackground) Background {
	description := ""
	feature := ""
	personalityTraits := []string{}
	ideals := []string{}
	equipment := []string{}
	languages := []string{}

	// Use text field for description if available
	if len(xmlBackground.Text) > 0 {
		description = strings.Join(xmlBackground.Text, " ")
	}

	// Parse traits to extract different information
	for _, trait := range xmlBackground.Traits {
		traitLower := strings.ToLower(trait.Name)
		if traitLower == "description" {
			description = trait.Text
		} else if strings.HasPrefix(traitLower, "feature") {
			feature = trait.Text
		} else if strings.Contains(traitLower, "suggested characteristics") ||
			strings.Contains(traitLower, "personality trait") ||
			strings.Contains(traitLower, "ideal") ||
			strings.Contains(traitLower, "bond") ||
			strings.Contains(traitLower, "flaw") {
			// Add personality-related traits
			personalityTraits = append(personalityTraits, trait.Text)
		} else if strings.Contains(traitLower, "ideal") {
			ideals = append(ideals, trait.Text)
		}
	}

	// Parse proficiencies (comma-separated skills)
	skills := strings.Split(xmlBackground.Proficiency, ",")
	for i, skill := range skills {
		skills[i] = strings.TrimSpace(skill)
	}

	// Extract equipment and languages from description text
	if description != "" {
		// Look for equipment and languages in text
		lines := strings.Split(description, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "• Equipment:") {
				// Extract equipment items
				equipText := strings.TrimPrefix(line, "• Equipment:")
				equipmentItems := strings.Split(equipText, ",")
				for _, item := range equipmentItems {
					trimmed := strings.TrimSpace(item)
					if trimmed != "" {
						equipment = append(equipment, trimmed)
					}
				}
			} else if strings.HasPrefix(line, "• Languages:") {
				// Extract languages
				langText := strings.TrimPrefix(line, "• Languages:")
				langItems := strings.Split(langText, ",")
				for _, lang := range langItems {
					trimmed := strings.TrimSpace(lang)
					if trimmed != "" {
						languages = append(languages, trimmed)
					}
				}
			}
		}
	}

	return Background{
		Name:               xmlBackground.Name,
		Description:        description,
		SkillProficiencies: skills,
		ToolProficiencies:  []string{}, // Empty since not in XML structure
		Equipment:          equipment,
		Feature:            feature,
		PersonalityTraits:  personalityTraits,
		Ideals:             ideals,
	}
}

// parseClass converts XMLClass to Class
func (cp *CompendiumParser) parseClass(xmlClass XMLClass) Class {
	skillsCount, _ := strconv.Atoi(xmlClass.Skills)

	// Parse proficiencies - this field contains both saving throws AND skill choices
	allProfs := strings.Split(xmlClass.Proficiency, ",")
	savingThrows := []string{}
	weaponProficiencies := []string{}
	skillsChoices := []string{}

	// Common saving throw abilities
	savingThrowAbilities := map[string]bool{
		"Strength":     true,
		"Dexterity":    true,
		"Constitution": true,
		"Intelligence": true,
		"Wisdom":       true,
		"Charisma":     true,
	}

	// Common skill abilities
	skillAbilities := map[string]bool{
		"Acrobatics":      true,
		"Animal Handling": true,
		"Arcana":          true,
		"Athletics":       true,
		"Deception":       true,
		"History":         true,
		"Insight":         true,
		"Intimidation":    true,
		"Investigation":   true,
		"Medicine":        true,
		"Nature":          true,
		"Perception":      true,
		"Performance":     true,
		"Persuasion":      true,
		"Religion":        true,
		"Sleight of Hand": true,
		"Stealth":         true,
		"Survival":        true,
	}

	// Analyze proficiencies to separate saving throws from skill choices
	for _, prof := range allProfs {
		prof = strings.TrimSpace(prof)
		if savingThrowAbilities[prof] {
			savingThrows = append(savingThrows, prof)
		} else if skillAbilities[prof] {
			skillsChoices = append(skillsChoices, prof)
		} else {
			// These are actual proficiencies (shouldn't happen for classes, but just in case)
			weaponProficiencies = append(weaponProficiencies, prof)
		}
	}

	// Parse tools proficiencies
	tools := []string{}
	if xmlClass.Tools != "None" && xmlClass.Tools != "" {
		tools = strings.Split(xmlClass.Tools, ",")
		for i, t := range tools {
			tools[i] = strings.TrimSpace(t)
		}
	}

	// Parse armor proficiencies
	armorProfs := strings.Split(xmlClass.Armor, ",")
	for i, a := range armorProfs {
		armorProfs[i] = strings.TrimSpace(a)
	}

	// Parse weapons proficiencies
	weaponProfs := strings.Split(xmlClass.Weapons, ",")
	for i, w := range weaponProfs {
		weaponProfs[i] = strings.TrimSpace(w)
	}

	// Parse starting equipment from wealth field
	equipment := []string{}
	if xmlClass.Wealth != "" {
		// Convert wealth notation like "5d4x10" to equipment description
		if strings.Contains(xmlClass.Wealth, "d4") && strings.Contains(xmlClass.Wealth, "x10") {
			equipment = []string{fmt.Sprintf("Starting gold: %s gp", xmlClass.Wealth)}
		}
	}

	return Class{
		Name:                xmlClass.Name,
		HitDie:              xmlClass.HD,
		PrimaryAbility:      xmlClass.SpellAbility,
		SavingThrows:        savingThrows,
		ArmorProficiencies:  armorProfs,
		WeaponProficiencies: weaponProfs,
		ToolsProficiencies:  tools,
		SkillsCount:         skillsCount,
		SkillsChoices:       skillsChoices,
		Equipment:           equipment,
		Description:         "", // Could be parsed from other XML sources if available
	}
}

// Helper functions for parsing specific formats
func (cp *CompendiumParser) parseComponents(components string) []string {
	// Split components and handle material component
	result := []string{}
	for _, c := range strings.Split(components, ",") {
		c = strings.TrimSpace(c)
		if c == "M" {
			result = append(result, "Material")
		} else if c == "V" {
			result = append(result, "Verbal")
		} else if c == "S" {
			result = append(result, "Somatic")
		} else if strings.HasPrefix(c, "M (") {
			result = append(result, "Material")
		} else {
			result = append(result, c)
		}
	}
	return result
}

func (cp *CompendiumParser) parseClasses(classes string) []string {
	if classes == "" {
		return []string{}
	}
	result := strings.Split(classes, ",")
	for i, class := range result {
		result[i] = strings.TrimSpace(class)
	}
	return result
}

func (cp *CompendiumParser) parseArmorClass(ac string) (int, string) {
	// Parse AC like "17 (natural armor)" or "15"
	acInt := 0
	acStr := ""

	parts := strings.SplitN(ac, " ", 2)
	if len(parts) > 0 {
		if val, err := strconv.Atoi(parts[0]); err == nil {
			acInt = val
		}
	}
	if len(parts) > 1 {
		acStr = parts[1]
	}

	return acInt, acStr
}

func (cp *CompendiumParser) parseHitPoints(hp string) (int, string) {
	// Parse HP like "135 (18d10+36)" or "22"
	hpInt := 0
	hpStr := ""

	parts := strings.SplitN(hp, " ", 2)
	if len(parts) > 0 {
		if val, err := strconv.Atoi(parts[0]); err == nil {
			hpInt = val
		}
	}
	if len(parts) > 1 {
		hpStr = parts[1]
	}

	return hpInt, hpStr
}

func (cp *CompendiumParser) parseAbility(ability string) int {
	if val, err := strconv.Atoi(ability); err == nil {
		return val
	}
	return 10
}

func (cp *CompendiumParser) parseSkills(skills string) map[string]int {
	result := make(map[string]int)
	if skills == "" {
		return result
	}

	// Parse skills like "Perception +3, Stealth +6"
	for _, skill := range strings.Split(skills, ",") {
		skill = strings.TrimSpace(skill)
		parts := strings.Split(skill, " ")
		if len(parts) >= 2 {
			name := strings.Join(parts[:len(parts)-1], " ")
			if val, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				result[name] = val
			}
		}
	}
	return result
}

func (cp *CompendiumParser) parseSenses(senses, passive string) map[string]interface{} {
	result := make(map[string]interface{})

	if senses != "" {
		// Parse senses like "darkvision 60 ft., passive Perception 12"
		for _, sense := range strings.Split(senses, ",") {
			sense = strings.TrimSpace(sense)
			if sense != "" {
				result[sense] = true
			}
		}
	}

	if passive != "" {
		result["Passive"] = passive
	}

	return result
}

func (cp *CompendiumParser) parseLanguages(languages string) []string {
	if languages == "" {
		return []string{}
	}
	result := strings.Split(languages, ",")
	for i, lang := range result {
		result[i] = strings.TrimSpace(lang)
	}
	return result
}

func (cp *CompendiumParser) parseTraits(traits []XMLTrait) []MonsterAbility {
	result := []MonsterAbility{}
	for _, trait := range traits {
		result = append(result, MonsterAbility{
			Name:        trait.Name,
			Description: trait.Text,
		})
	}
	return result
}

func (cp *CompendiumParser) parseActions(actions []XMLAction) []MonsterAbility {
	result := []MonsterAbility{}
	for _, action := range actions {
		result = append(result, MonsterAbility{
			Name:        action.Name,
			Description: action.Text,
			Hit:         action.Attack,
		})
	}
	return result
}

func (cp *CompendiumParser) parseAbilityBonus(ability string) map[string]int {
	result := make(map[string]int)

	// Parse ability bonus like "Constitution 2, Wisdom 1"
	for _, bonus := range strings.Split(ability, ",") {
		bonus = strings.TrimSpace(bonus)
		parts := strings.Split(bonus, " ")
		if len(parts) == 2 {
			if val, err := strconv.Atoi(parts[1]); err == nil {
				result[parts[0]] = val
			}
		}
	}

	return result
}

// Get methods for retrieving data by name
func (cp *CompendiumParser) GetSpellByName(name string) (*Spell, error) {
	for _, spell := range cp.spells {
		if strings.EqualFold(spell.Name, name) {
			return &spell, nil
		}
	}
	return nil, fmt.Errorf("spell '%s' not found", name)
}

func (cp *CompendiumParser) GetMonsterByName(name string) (*Monster, error) {
	for _, monster := range cp.monsters {
		if strings.EqualFold(monster.Name, name) {
			return &monster, nil
		}
	}
	return nil, fmt.Errorf("monster '%s' not found", name)
}

func (cp *CompendiumParser) GetItemByName(name string) (*Item, error) {
	for _, item := range cp.items {
		if strings.EqualFold(item.Name, name) {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item '%s' not found", name)
}

func (cp *CompendiumParser) GetRaceByName(name string) (*Species, error) {
	for _, species := range cp.species {
		if strings.EqualFold(species.Name, name) {
			return &species, nil
		}
	}
	return nil, fmt.Errorf("species '%s' not found", name)
}

func (cp *CompendiumParser) GetBackgroundByName(name string) (*Background, error) {
	for _, background := range cp.backgrounds {
		if strings.EqualFold(background.Name, name) {
			return &background, nil
		}
	}
	return nil, fmt.Errorf("background '%s' not found", name)
}

func (cp *CompendiumParser) GetClassByName(name string) (*Class, error) {
	for _, class := range cp.classes {
		if strings.EqualFold(class.Name, name) {
			return &class, nil
		}
	}
	return nil, fmt.Errorf("class '%s' not found", name)
}

// Get methods for retrieving all names
func (cp *CompendiumParser) GetSpellNames() []string {
	names := make([]string, len(cp.spells))
	for i, spell := range cp.spells {
		names[i] = spell.Name
	}
	return names
}

func (cp *CompendiumParser) GetMonsterNames() []string {
	names := make([]string, len(cp.monsters))
	for i, monster := range cp.monsters {
		names[i] = monster.Name
	}
	return names
}

func (cp *CompendiumParser) GetItemNames() []string {
	names := make([]string, len(cp.items))
	for i, item := range cp.items {
		names[i] = item.Name
	}
	return names
}

func (cp *CompendiumParser) GetRaceNames() []string {
	names := make([]string, len(cp.species))
	for i, species := range cp.species {
		names[i] = species.Name
	}
	return names
}

func (cp *CompendiumParser) GetBackgroundNames() []string {
	names := make([]string, len(cp.backgrounds))
	for i, background := range cp.backgrounds {
		names[i] = background.Name
	}
	return names
}

func (cp *CompendiumParser) GetClassNames() []string {
	names := make([]string, len(cp.classes))
	for i, class := range cp.classes {
		names[i] = class.Name
	}
	return names
}

// ListAvailableCompendiumFiles returns available XML files in a directory
func ListAvailableCompendiumFiles(compendiumDir string) ([]string, error) {
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

	return names, nil
}
