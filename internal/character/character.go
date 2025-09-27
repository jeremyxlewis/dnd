package character

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Character represents a D&D 5e character
type Character struct {
	Name       string `json:"name"`
	Species    string `json:"species"`
	Class      string `json:"class"`
	Level      int    `json:"level"`
	Background string `json:"background"`
	Subclass   string `json:"subclass,omitempty"`

	// Ability Scores
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`

	// Derived Stats
	HitPoints        int `json:"hit_points"`
	CurrentHP        int `json:"current_hp"`
	TempHP           int `json:"temp_hp"`
	ArmorClass       int `json:"armor_class"`
	ProficiencyBonus int `json:"proficiency_bonus"`
	Speed            int `json:"speed"`

	// Proficiencies
	ArmorProficiencies  []string `json:"armor_proficiencies,omitempty"`
	WeaponProficiencies []string `json:"weapon_proficiencies,omitempty"`
	ToolProficiencies   []string `json:"tool_proficiencies,omitempty"`
	SkillProficiencies  []string `json:"skill_proficiencies,omitempty"`
	Languages           []string `json:"languages,omitempty"`

	// Features and Traits
	Features []string `json:"features,omitempty"`

	// Spellcasting (for spellcasters)
	SpellcastingAbility string      `json:"spellcasting_ability,omitempty"`
	SpellSlots          map[int]int `json:"spell_slots,omitempty"`      // level -> count
	UsedSpellSlots      map[int]int `json:"used_spell_slots,omitempty"` // level -> used
	SpellsKnown         []string    `json:"spells_known,omitempty"`
	SpellsPrepared      []string    `json:"spells_prepared,omitempty"`

	// Equipment and Inventory
	Equipment []string `json:"equipment,omitempty"`

	// Other
	HitDice     string   `json:"hit_dice"`
	Alignment   string   `json:"alignment,omitempty"`
	Experience  int      `json:"experience"`
	Inspiration bool     `json:"inspiration"`
	Conditions  []string `json:"conditions,omitempty"`
	Backstory   string   `json:"backstory,omitempty"`
}

// NewCharacter creates a new character with default values
func NewCharacter(name, species, class, background, alignment string, level int, str, dex, con, intelligence, wisdom, charisma int) *Character {
	return &Character{
		Name:                name,
		Species:             species,
		Class:               class,
		Level:               level,
		Background:          background,
		Alignment:           alignment,
		Strength:            str,
		Dexterity:           dex,
		Constitution:        con,
		Intelligence:        intelligence,
		Wisdom:              wisdom,
		Charisma:            charisma,
		HitPoints:           10, // Placeholder, will be updated with racial/class logic
		CurrentHP:           10,
		TempHP:              0,
		ArmorClass:          10, // Placeholder
		ProficiencyBonus:    2,
		Speed:               30, // Default bipedal speed
		ArmorProficiencies:  []string{},
		WeaponProficiencies: []string{},
		ToolProficiencies:   []string{},
		SkillProficiencies:  []string{},
		Languages:           []string{"Common"},
		Features:            []string{},
		SpellSlots:          make(map[int]int),
		UsedSpellSlots:      make(map[int]int),
		SpellsKnown:         []string{},
		SpellsPrepared:      []string{},
		Equipment:           []string{},
		HitDice:             "1d8", // Placeholder, class-dependent
		Experience:          0,
		Inspiration:         false,
		Conditions:          []string{},
		Backstory:           "",
	}
}

// SaveCharacter saves a character to a JSON file
func SaveCharacter(char *Character, filePath string) error {
	data, err := json.MarshalIndent(char, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write character file: %w", err)
	}
	return nil
}

// LoadCharacter loads a character from a JSON file
func LoadCharacter(filePath string) (*Character, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	var char Character
	err = json.Unmarshal(data, &char)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal character: %w", err)
	}
	return &char, nil
}

// GetCharacterFilePath returns the standard path for a character file
func GetCharacterFilePath(charName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create a .dnd-cli directory in the user's home directory
	appDir := filepath.Join(homeDir, ".dnd-cli")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err = os.Mkdir(appDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create application directory: %w", err)
		}
	}

	return filepath.Join(appDir, fmt.Sprintf("%s.json", charName)), nil
}

// ApplyRacialTraits applies racial ability bonuses, speed, languages, and features
func (c *Character) ApplyRacialTraits() {
	switch c.Species {
	case "Human":
		c.Strength++
		c.Dexterity++
		c.Constitution++
		c.Intelligence++
		c.Wisdom++
		c.Charisma++
		c.Languages = append(c.Languages, "One extra language")
		c.Features = append(c.Features, "Versatile")
	case "Elf", "High Elf", "Wood Elf", "Dark Elf":
		c.Dexterity += 2
		c.Speed = 30
		c.Languages = append(c.Languages, "Elvish")
		c.Features = append(c.Features, "Darkvision", "Fey Ancestry", "Trance")
		if c.Species == "High Elf" {
			c.Intelligence++
			c.WeaponProficiencies = append(c.WeaponProficiencies, "Longsword", "Shortsword", "Shortbow", "Longbow")
			c.SpellsKnown = append(c.SpellsKnown, "One wizard cantrip")
		} else if c.Species == "Wood Elf" {
			c.Wisdom++
			c.WeaponProficiencies = append(c.WeaponProficiencies, "Longsword", "Shortsword", "Shortbow", "Longbow")
			c.Speed = 35
			c.Features = append(c.Features, "Mask of the Wild")
		} else if c.Species == "Dark Elf" {
			c.Charisma++
			c.WeaponProficiencies = append(c.WeaponProficiencies, "Rapier", "Shortsword", "Hand Crossbow")
			c.Features = append(c.Features, "Sunlight Sensitivity", "Drow Magic")
		}
	case "Dwarf", "Hill Dwarf", "Mountain Dwarf":
		c.Constitution += 2
		c.Speed = 25
		c.Languages = append(c.Languages, "Dwarvish")
		c.Features = append(c.Features, "Darkvision", "Dwarven Resilience", "Stonecunning")
		if c.Species == "Hill Dwarf" {
			c.Wisdom++
			c.HitPoints += 1 // per level
		} else if c.Species == "Mountain Dwarf" {
			c.Strength++
			c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor")
		}
	case "Halfling", "Lightfoot", "Stout":
		c.Dexterity += 2
		c.Speed = 25
		c.Languages = append(c.Languages, "Halfling")
		c.Features = append(c.Features, "Lucky", "Brave", "Halfling Nimbleness")
		if c.Species == "Lightfoot" {
			c.Charisma++
			c.Features = append(c.Features, "Naturally Stealthy")
		} else if c.Species == "Stout" {
			c.Constitution++
			c.Features = append(c.Features, "Stout Resilience")
		}
	case "Half-Elf":
		c.Charisma += 2
		// +1 to two other abilities (placeholder: +1 Str, +1 Dex)
		c.Strength++
		c.Dexterity++
		c.Speed = 30
		c.Languages = append(c.Languages, "Elvish", "One extra language")
		c.Features = append(c.Features, "Darkvision", "Fey Ancestry", "Skill Versatility")
	case "Half-Orc":
		c.Strength += 2
		c.Constitution++
		c.Speed = 30
		c.Languages = append(c.Languages, "Orc")
		c.Features = append(c.Features, "Darkvision", "Menacing", "Relentless Endurance", "Savage Attacks")
	case "Tiefling":
		c.Intelligence++
		c.Charisma += 2
		c.Speed = 30
		c.Languages = append(c.Languages, "Infernal")
		c.Features = append(c.Features, "Darkvision", "Hellish Resistance", "Infernal Legacy")
	case "Dragonborn":
		c.Strength += 2
		c.Charisma++
		c.Speed = 30
		c.Languages = append(c.Languages, "Draconic")
		c.Features = append(c.Features, "Draconic Ancestry", "Breath Weapon", "Damage Resistance")
	case "Gnome", "Forest Gnome", "Rock Gnome":
		c.Intelligence += 2
		c.Speed = 25
		c.Languages = append(c.Languages, "Gnomish")
		c.Features = append(c.Features, "Darkvision", "Gnome Cunning")
		if c.Species == "Forest Gnome" {
			c.Dexterity++
			c.Features = append(c.Features, "Natural Illusionist", "Speak with Small Beasts")
		} else if c.Species == "Rock Gnome" {
			c.Constitution++
			c.Features = append(c.Features, "Artificer's Lore", "Tinker")
		}
	case "Tabaxi":
		c.Dexterity += 2
		c.Charisma++
		c.Speed = 30
		c.Languages = append(c.Languages, "One extra language")
		c.Features = append(c.Features, "Darkvision", "Feline Agility", "Cat's Claws", "Cat's Talent")
	case "Genasi", "Air Genasi", "Earth Genasi", "Fire Genasi", "Water Genasi":
		c.Constitution += 2
		c.Speed = 30
		c.Languages = append(c.Languages, "Primordial")
		if c.Species == "Air Genasi" {
			c.Dexterity++
			c.Features = append(c.Features, "Unending Breath", "Mingle with the Wind")
		} else if c.Species == "Earth Genasi" {
			c.Strength++
			c.Features = append(c.Features, "Earth Walk", "Merge with Stone")
		} else if c.Species == "Fire Genasi" {
			c.Intelligence++
			c.Features = append(c.Features, "Darkvision", "Fire Resistance", "Reach to the Blaze")
		} else if c.Species == "Water Genasi" {
			c.Wisdom++
			c.Features = append(c.Features, "Acid Resistance", "Amphibious", "Swim", "Call to the Wave")
		}
	case "Aarakocra":
		c.Dexterity += 2
		c.Wisdom++
		c.Speed = 25
		c.Languages = append(c.Languages, "Aarakocra", "Auran")
		c.Features = append(c.Features, "Flight", "Talons")
	case "Aasimar", "Protector Aasimar", "Scourge Aasimar", "Fallen Aasimar":
		c.Wisdom++
		c.Charisma += 2
		c.Speed = 30
		c.Languages = append(c.Languages, "Celestial")
		c.Features = append(c.Features, "Darkvision", "Celestial Resistance", "Healing Hands")
		if c.Species == "Protector Aasimar" {
			c.Features = append(c.Features, "Radiant Soul")
		} else if c.Species == "Scourge Aasimar" {
			c.Features = append(c.Features, "Radiant Consumption")
		} else if c.Species == "Fallen Aasimar" {
			c.Strength++
			c.Features = append(c.Features, "Necrotic Shroud")
		}
	case "Bugbear":
		c.Strength += 2
		c.Dexterity++
		c.Speed = 30
		c.Languages = append(c.Languages, "Goblin")
		c.Features = append(c.Features, "Darkvision", "Long-Limbed", "Powerful Build", "Sneaky", "Surprise Attack")
	case "Centaur":
		c.Strength += 2
		c.Wisdom++
		c.Speed = 40
		c.Languages = append(c.Languages, "Sylvan")
		c.Features = append(c.Features, "Charge", "Hooves", "Equine Build", "Survivor")
	case "Changeling":
		c.Charisma += 2
		c.Speed = 30
		c.Languages = append(c.Languages, "Two extra languages")
		c.Features = append(c.Features, "Shapechanger")
	case "Duergar":
		c.Strength += 2
		c.Constitution++
		c.Speed = 25
		c.Languages = append(c.Languages, "Dwarvish", "Undercommon")
		c.Features = append(c.Features, "Superior Darkvision", "Duergar Resilience", "Duergar Magic", "Sunlight Sensitivity")
	case "Eladrin":
		c.Dexterity += 2
		c.Charisma++
		c.Speed = 30
		c.Languages = append(c.Languages, "Elvish")
		c.Features = append(c.Features, "Darkvision", "Fey Ancestry", "Trance", "Fey Step")
	case "Fairy":
		c.Dexterity += 2
		c.Charisma++
		c.Speed = 30
		c.Languages = append(c.Languages, "Sylvan")
		c.Features = append(c.Features, "Flight", "Fairy Magic")
	case "Firbolg":
		c.Wisdom += 2
		c.Strength++
		c.Speed = 30
		c.Languages = append(c.Languages, "Elvish", "Giant")
		c.Features = append(c.Features, "Firbolg Magic", "Hidden Step", "Powerful Build", "Speech of Beast and Leaf")
	case "Gith", "Githyanki", "Githzerai":
		c.Intelligence += 2
		c.Speed = 30
		c.Languages = append(c.Languages, "Gith")
		c.Features = append(c.Features, "Martial Prodigy", "Decadent Mastery")
		if c.Species == "Githyanki" {
			c.Strength++
			c.Features = append(c.Features, "Astronaut", "Githyanki Psionics")
		} else if c.Species == "Githzerai" {
			c.Wisdom++
			c.Features = append(c.Features, "Mental Discipline", "Githzerai Psionics")
		}
	case "Goblin":
		c.Dexterity += 2
		c.Constitution++
		c.Speed = 30
		c.Languages = append(c.Languages, "Goblin")
		c.Features = append(c.Features, "Darkvision", "Fury of the Small", "Nimble Escape")
	case "Goliath":
		c.Strength += 2
		c.Constitution++
		c.Speed = 30
		c.Languages = append(c.Languages, "Giant")
		c.Features = append(c.Features, "Stone's Endurance", "Powerful Build", "Mountain Born")
	case "Harengon":
		c.Dexterity += 2
		c.Wisdom++
		c.Speed = 30
		c.Languages = append(c.Languages, "One extra language")
		c.Features = append(c.Features, "Hare-Trigger", "Leporine Senses", "Lucky Footwork", "Rabbit Hop")
	case "Hobgoblin":
		c.Constitution += 2
		c.Intelligence++
		c.Speed = 30
		c.Languages = append(c.Languages, "Goblin")
		c.Features = append(c.Features, "Darkvision", "Martial Training", "Saving Face")
	case "Kenku":
		c.Dexterity += 2
		c.Wisdom++
		c.Speed = 30
		c.Languages = append(c.Languages, "Auran")
		c.Features = append(c.Features, "Expert Forgery", "Kenku Training", "Mimicry")
	case "Kobold":
		c.Dexterity += 2
		c.Strength--
		c.Speed = 30
		c.Languages = append(c.Languages, "Draconic")
		c.Features = append(c.Features, "Darkvision", "Grovel, Cower, and Beg", "Pack Tactics", "Sunlight Sensitivity")
	case "Lizardfolk":
		c.Constitution += 2
		c.Wisdom++
		c.Speed = 30
		c.Languages = append(c.Languages, "Draconic")
		c.Features = append(c.Features, "Bite", "Cunning Artisan", "Hold Breath", "Hunter's Lore", "Natural Armor", "Hungry Jaws")
	case "Minotaur":
		c.Strength += 2
		c.Constitution++
		c.Speed = 30
		c.Languages = append(c.Languages, "Minotaur")
		c.Features = append(c.Features, "Horns", "Goring Rush", "Hammering Horns", "Imposing Presence")
	case "Orc":
		c.Strength += 2
		c.Constitution++
		c.Speed = 30
		c.Languages = append(c.Languages, "Orc")
		c.Features = append(c.Features, "Darkvision", "Aggressive", "Menacing", "Powerful Build", "Primal Intuition")
	case "Satyr":
		c.Charisma += 2
		c.Dexterity++
		c.Speed = 35
		c.Languages = append(c.Languages, "Sylvan")
		c.Features = append(c.Features, "Magic Resistance", "Mirthful Leaps", "Reveler", "Satyr Magic")
	case "Sea Elf":
		c.Constitution += 2
		c.Dexterity++
		c.Speed = 30
		c.Languages = append(c.Languages, "Elvish", "Aquan")
		c.Features = append(c.Features, "Darkvision", "Fey Ancestry", "Trance", "Child of the Sea", "Friend of the Sea")
	case "Shifter", "Beasthide", "Cliffwalk", "Longstride", "Longtooth", "Razorclaw", "Wildhunt":
		c.Dexterity += 2
		c.Strength++
		c.Speed = 30
		c.Languages = append(c.Languages, "One extra language")
		c.Features = append(c.Features, "Darkvision", "Shifting")
		if c.Species == "Beasthide" {
			c.Constitution++
		} else if c.Species == "Cliffwalk" {
			c.Dexterity++
		} else if c.Species == "Longstride" {
			c.Dexterity++
		} else if c.Species == "Longtooth" {
			c.Strength++
		} else if c.Species == "Razorclaw" {
			c.Dexterity++
		} else if c.Species == "Wildhunt" {
			c.Wisdom++
		}
	case "Tortle":
		c.Strength += 2
		c.Wisdom++
		c.Speed = 30
		c.Languages = append(c.Languages, "Aquan")
		c.Features = append(c.Features, "Claws", "Hold Breath", "Natural Armor", "Shell Defense", "Survival Instinct")
	case "Triton":
		c.Strength++
		c.Constitution++
		c.Charisma++
		c.Speed = 30
		c.Languages = append(c.Languages, "Primordial")
		c.Features = append(c.Features, "Amphibious", "Control Air and Water", "Emissary of the Sea", "Guardians of the Depths")
	case "Yuan-ti Pureblood":
		c.Charisma += 2
		c.Intelligence++
		c.Speed = 30
		c.Languages = append(c.Languages, "Abyssal", "Draconic")
		c.Features = append(c.Features, "Darkvision", "Innate Spellcasting", "Magic Resistance", "Poison Immunity")
	default:
		// No changes for unknown species
	}
}

// ApplyClassTraits applies class-specific starting proficiencies, hit dice, and features
func (c *Character) ApplyClassTraits() {
	switch c.Class {
	case "Barbarian":
		c.HitDice = "1d12"
		c.HitPoints = 12 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor", "Shields")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Martial weapons")
		c.ToolProficiencies = append(c.ToolProficiencies, "One artisan's tools or one musical instrument")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Animal Handling, Athletics, Intimidation, Nature, Perception, Survival")
		c.Features = append(c.Features, "Rage", "Unarmored Defense")
	case "Bard":
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Hand crossbows", "Longswords", "Rapiers", "Shortswords")
		c.ToolProficiencies = append(c.ToolProficiencies, "Three musical instruments")
		c.SkillProficiencies = append(c.SkillProficiencies, "Three from: Acrobatics, Animal Handling, Arcana, Athletics, Deception, History, Insight, Intimidation, Investigation, Medicine, Nature, Perception, Performance, Persuasion, Religion, Sleight of Hand, Stealth, Survival")
		c.SpellcastingAbility = "Charisma"
		c.Features = append(c.Features, "Bardic Inspiration", "Spellcasting")
		c.SpellSlots = map[int]int{1: 2}
	case "Cleric":
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor", "Shields")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: History, Insight, Medicine, Persuasion, Religion")
		c.SpellcastingAbility = "Wisdom"
		c.Features = append(c.Features, "Divine Domain", "Spellcasting")
		c.SpellSlots = map[int]int{1: 2}
	case "Druid":
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor", "Shields (non-metal)")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Clubs", "Daggers", "Darts", "Javelins", "Maces", "Quarterstaffs", "Scimitars", "Sickles", "Slings", "Spears")
		c.ToolProficiencies = append(c.ToolProficiencies, "Herbalism kit")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Arcana, Animal Handling, Insight, Medicine, Nature, Perception, Religion, Survival")
		c.SpellcastingAbility = "Wisdom"
		c.Languages = append(c.Languages, "Druidic")
		c.Features = append(c.Features, "Druidcraft", "Spellcasting")
		c.SpellSlots = map[int]int{1: 2}
	case "Fighter":
		c.HitDice = "1d10"
		c.HitPoints = 10 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor", "Heavy armor", "Shields")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Martial weapons")
		c.ToolProficiencies = append(c.ToolProficiencies, "One artisan's tools")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Acrobatics, Animal Handling, Athletics, History, Insight, Intimidation, Perception, Survival")
		c.Features = append(c.Features, "Fighting Style", "Second Wind")
	case "Monk":
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Shortswords")
		c.ToolProficiencies = append(c.ToolProficiencies, "One artisan's tools or one musical instrument")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Acrobatics, Athletics, History, Insight, Religion, Stealth")
		c.Features = append(c.Features, "Unarmored Defense", "Martial Arts")
	case "Paladin":
		c.HitDice = "1d10"
		c.HitPoints = 10 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor", "Heavy armor", "Shields")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Martial weapons")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Athletics, Insight, Intimidation, Medicine, Persuasion, Religion")
		c.SpellcastingAbility = "Charisma"
		c.Features = append(c.Features, "Divine Sense", "Lay on Hands")
	case "Ranger":
		c.HitDice = "1d10"
		c.HitPoints = 10 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor", "Medium armor", "Shields")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Martial weapons")
		c.SkillProficiencies = append(c.SkillProficiencies, "Three from: Animal Handling, Athletics, Insight, Investigation, Nature, Perception, Stealth, Survival")
		c.Languages = append(c.Languages, "One extra language")
		c.SpellcastingAbility = "Wisdom"
		c.Features = append(c.Features, "Favored Enemy", "Natural Explorer")
	case "Rogue":
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons", "Hand crossbows", "Longswords", "Rapiers", "Shortswords")
		c.ToolProficiencies = append(c.ToolProficiencies, "Thieves' tools")
		c.SkillProficiencies = append(c.SkillProficiencies, "Four from: Acrobatics, Athletics, Deception, Insight, Intimidation, Investigation, Perception, Performance, Persuasion, Sleight of Hand, Stealth")
		c.Features = append(c.Features, "Expertise", "Sneak Attack", "Thieves' Cant")
	case "Sorcerer":
		c.HitDice = "1d6"
		c.HitPoints = 6 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Daggers", "Darts", "Slings", "Quarterstaffs", "Light crossbows")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Arcana, Deception, Insight, Intimidation, Persuasion, Religion")
		c.SpellcastingAbility = "Charisma"
		c.Features = append(c.Features, "Spellcasting", "Sorcerous Origin")
		c.SpellSlots = map[int]int{1: 2}
	case "Warlock":
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.ArmorProficiencies = append(c.ArmorProficiencies, "Light armor")
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Simple weapons")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Arcana, Deception, History, Intimidation, Investigation, Nature, Religion")
		c.SpellcastingAbility = "Charisma"
		c.Features = append(c.Features, "Otherworldly Patron", "Pact Magic")
		c.SpellSlots = map[int]int{1: 1}
	case "Wizard":
		c.HitDice = "1d6"
		c.HitPoints = 6 + (c.Constitution-10)/2
		c.CurrentHP = c.HitPoints
		c.WeaponProficiencies = append(c.WeaponProficiencies, "Daggers", "Darts", "Slings", "Quarterstaffs", "Light crossbows")
		c.SkillProficiencies = append(c.SkillProficiencies, "Two from: Arcana, History, Insight, Investigation, Medicine, Religion")
		c.SpellcastingAbility = "Intelligence"
		c.Features = append(c.Features, "Spellcasting", "Arcane Recovery")
		c.SpellSlots = map[int]int{1: 2}
	default:
		// Default to d8 hit die
		c.HitDice = "1d8"
		c.HitPoints = 8 + (c.Constitution-10)/2
	}
}

// ApplyBackgroundTraits applies background proficiencies and features
func (c *Character) ApplyBackgroundTraits() {
	switch c.Background {
	case "Acolyte":
		c.SkillProficiencies = append(c.SkillProficiencies, "Insight", "Religion")
		c.Languages = append(c.Languages, "Two extra languages")
		c.Equipment = append(c.Equipment, "Holy symbol", "Prayer book", "5 candles", "Tinderbox", "Alms box", "2 blocks of incense", "Censer", "Vestments", "2 rations", "Waterskin")
		c.Features = append(c.Features, "Shelter of the Faithful")
	case "Charlatan":
		c.SkillProficiencies = append(c.SkillProficiencies, "Deception", "Sleight of Hand")
		c.ToolProficiencies = append(c.ToolProficiencies, "Disguise kit", "Forgery kit")
		c.Equipment = append(c.Equipment, "Fine clothes", "Disguise kit", "Con tools", "15 gp")
		c.Features = append(c.Features, "False Identity")
	case "Criminal":
		c.SkillProficiencies = append(c.SkillProficiencies, "Deception", "Stealth")
		c.ToolProficiencies = append(c.ToolProficiencies, "One gaming set", "Thieves' tools")
		c.Equipment = append(c.Equipment, "Crowbar", "Dark common clothes", "15 gp")
		c.Features = append(c.Features, "Criminal Contact")
	case "Entertainer":
		c.SkillProficiencies = append(c.SkillProficiencies, "Acrobatics", "Performance")
		c.ToolProficiencies = append(c.ToolProficiencies, "Disguise kit", "One musical instrument")
		c.Equipment = append(c.Equipment, "Musical instrument", "Favor of an admirer", "Costume", "15 gp")
		c.Features = append(c.Features, "By Popular Demand")
	case "Folk Hero":
		c.SkillProficiencies = append(c.SkillProficiencies, "Animal Handling", "Survival")
		c.ToolProficiencies = append(c.ToolProficiencies, "One artisan's tools", "Vehicles (land)")
		c.Equipment = append(c.Equipment, "Artisan's tools", "Shovel", "Iron pot", "Common clothes", "10 gp")
		c.Features = append(c.Features, "Rustic Hospitality")
	case "Guild Artisan":
		c.SkillProficiencies = append(c.SkillProficiencies, "Insight", "Persuasion")
		c.Languages = append(c.Languages, "One extra language")
		c.ToolProficiencies = append(c.ToolProficiencies, "One artisan's tools")
		c.Equipment = append(c.Equipment, "Artisan's tools", "Letter of introduction", "Traveler's clothes", "15 gp")
		c.Features = append(c.Features, "Guild Membership")
	case "Hermit":
		c.SkillProficiencies = append(c.SkillProficiencies, "Medicine", "Religion")
		c.Languages = append(c.Languages, "One extra language")
		c.ToolProficiencies = append(c.ToolProficiencies, "Herbalism kit")
		c.Equipment = append(c.Equipment, "Scroll case of notes", "Winter blanket", "Common clothes", "5 gp")
		c.Features = append(c.Features, "Discovery")
	case "Noble":
		c.SkillProficiencies = append(c.SkillProficiencies, "History", "Persuasion")
		c.Languages = append(c.Languages, "One extra language")
		c.ToolProficiencies = append(c.ToolProficiencies, "One gaming set")
		c.Equipment = append(c.Equipment, "Fine clothes", "Signet ring", "Scroll of pedigree", "25 gp")
		c.Features = append(c.Features, "Position of Privilege")
	case "Outlander":
		c.SkillProficiencies = append(c.SkillProficiencies, "Athletics", "Survival")
		c.Languages = append(c.Languages, "One extra language")
		c.ToolProficiencies = append(c.ToolProficiencies, "One musical instrument")
		c.Equipment = append(c.Equipment, "Staff", "Hunting trap", "Traveler's clothes", "10 gp")
		c.Features = append(c.Features, "Wanderer")
	case "Sage":
		c.SkillProficiencies = append(c.SkillProficiencies, "Arcana", "History")
		c.Languages = append(c.Languages, "Two extra languages")
		c.Equipment = append(c.Equipment, "Bottle of ink", "Quill", "Small knife", "Letter from colleague", "Common clothes", "10 gp")
		c.Features = append(c.Features, "Researcher")
	case "Sailor":
		c.SkillProficiencies = append(c.SkillProficiencies, "Athletics", "Perception")
		c.ToolProficiencies = append(c.ToolProficiencies, "Navigator's tools", "Vehicles (water)")
		c.Equipment = append(c.Equipment, "Belaying pin", "50 feet of silk rope", "Lucky charm", "Common clothes", "10 gp")
		c.Features = append(c.Features, "Bad Reputation")
	case "Soldier":
		c.SkillProficiencies = append(c.SkillProficiencies, "Athletics", "Intimidation")
		c.ToolProficiencies = append(c.ToolProficiencies, "One gaming set", "Vehicles (land)")
		c.Equipment = append(c.Equipment, "Insignia of rank", "Trophy from fallen enemy", "Bone dice or deck of cards", "Common clothes", "10 gp")
		c.Features = append(c.Features, "Military Rank")
	case "Urchin":
		c.SkillProficiencies = append(c.SkillProficiencies, "Sleight of Hand", "Stealth")
		c.ToolProficiencies = append(c.ToolProficiencies, "Disguise kit", "Thieves' tools")
		c.Equipment = append(c.Equipment, "Small knife", "Map of home city", "Pet mouse", "Token of parents", "Common clothes", "10 gp")
		c.Features = append(c.Features, "City Secrets")
	default:
		// No changes
	}
}

// LevelUp increases character level and applies class-specific changes
func (c *Character) LevelUp() {
	c.Level++
	// HP gain: use average hit die + Con modifier
	hitDieAvg := 5 // default
	switch c.Class {
	case "Barbarian":
		hitDieAvg = 7
	case "Bard", "Cleric", "Druid", "Monk", "Rogue", "Warlock":
		hitDieAvg = 5
	case "Fighter", "Paladin", "Ranger":
		hitDieAvg = 6
	case "Sorcerer", "Wizard":
		hitDieAvg = 4
	}
	c.HitPoints += hitDieAvg + (c.Constitution-10)/2

	// Proficiency bonus
	if c.Level == 5 || c.Level == 9 || c.Level == 13 || c.Level == 17 {
		c.ProficiencyBonus++
	}

	// Class-specific features and spell slots
	switch c.Class {
	case "Barbarian":
		if c.Level == 2 {
			c.Features = append(c.Features, "Reckless Attack", "Danger Sense")
		} else if c.Level == 3 {
			c.Subclass = "Berserker" // placeholder
			c.Features = append(c.Features, "Primal Path")
		}
		// Rage increases, etc.
	case "Bard":
		if c.Level == 2 {
			c.Features = append(c.Features, "Jack of All Trades", "Song of Rest")
		} else if c.Level == 3 {
			c.Subclass = "College of Valor" // placeholder
			c.Features = append(c.Features, "Bard College")
		}
		updateSpellSlots(c)
	case "Cleric":
		if c.Level == 2 {
			c.Features = append(c.Features, "Channel Divinity", "Divine Domain feature")
		}
		updateSpellSlots(c)
	case "Druid":
		if c.Level == 2 {
			c.Features = append(c.Features, "Wild Shape")
		} else if c.Level == 3 {
			c.Subclass = "Circle of the Moon" // placeholder
			c.Features = append(c.Features, "Druid Circle")
		}
		updateSpellSlots(c)
	case "Fighter":
		if c.Level == 2 {
			c.Features = append(c.Features, "Action Surge")
		} else if c.Level == 3 {
			c.Subclass = "Champion" // placeholder
			c.Features = append(c.Features, "Martial Archetype")
		}
	case "Monk":
		if c.Level == 2 {
			c.Features = append(c.Features, "Ki", "Unarmored Movement")
		} else if c.Level == 3 {
			c.Subclass = "Way of the Open Hand" // placeholder
			c.Features = append(c.Features, "Monastic Tradition")
		}
	case "Paladin":
		if c.Level == 2 {
			c.Features = append(c.Features, "Divine Smite", "Fighting Style")
		} else if c.Level == 3 {
			c.Subclass = "Oath of Devotion" // placeholder
			c.Features = append(c.Features, "Divine Health", "Sacred Oath")
		}
		if c.Level >= 2 {
			updateSpellSlots(c)
		}
	case "Ranger":
		if c.Level == 2 {
			c.Features = append(c.Features, "Fighting Style", "Spellcasting")
		} else if c.Level == 3 {
			c.Subclass = "Hunter" // placeholder
			c.Features = append(c.Features, "Ranger Archetype")
		}
		if c.Level >= 2 {
			updateSpellSlots(c)
		}
	case "Rogue":
		if c.Level == 2 {
			c.Features = append(c.Features, "Cunning Action")
		} else if c.Level == 3 {
			c.Subclass = "Thief" // placeholder
			c.Features = append(c.Features, "Roguish Archetype")
		}
	case "Sorcerer":
		if c.Level == 2 {
			c.Features = append(c.Features, "Font of Magic")
		}
		updateSpellSlots(c)
	case "Warlock":
		if c.Level == 2 {
			c.Features = append(c.Features, "Eldritch Invocations")
		}
		updateSpellSlots(c)
	case "Wizard":
		if c.Level == 2 {
			c.Features = append(c.Features, "Arcane Tradition")
		}
		updateSpellSlots(c)
	}
}

// updateSpellSlots updates spell slots based on level for spellcasters
func updateSpellSlots(c *Character) {
	if c.SpellcastingAbility == "" {
		return
	}
	// Spell slot tables (simplified, full caster progression)
	switch c.Class {
	case "Bard", "Cleric", "Druid", "Sorcerer", "Wizard":
		switch c.Level {
		case 1:
			c.SpellSlots[1] = 2
		case 2:
			c.SpellSlots[1] = 3
		case 3:
			c.SpellSlots[1] = 4
			c.SpellSlots[2] = 2
		case 4:
			c.SpellSlots[1] = 4
			c.SpellSlots[2] = 3
		case 5:
			c.SpellSlots[1] = 4
			c.SpellSlots[2] = 3
			c.SpellSlots[3] = 2
			// Add more levels as needed
		}
	case "Paladin", "Ranger":
		// Half caster
		switch c.Level {
		case 2:
			c.SpellSlots[1] = 2
		case 3:
			c.SpellSlots[1] = 3
		case 5:
			c.SpellSlots[1] = 4
			c.SpellSlots[2] = 2
			// Add more
		}
	case "Warlock":
		// Pact magic
		switch c.Level {
		case 1:
			c.SpellSlots[1] = 1
		case 2:
			c.SpellSlots[1] = 2
		case 3:
			c.SpellSlots[2] = 2
		case 5:
			c.SpellSlots[3] = 2
			// Add more
		}
	}
}
