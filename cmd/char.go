package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dnd-cli/internal/character"
	"dnd-cli/internal/data"

	"github.com/spf13/cobra"
)

// charCmd represents the char command
var charCmd = &cobra.Command{
	Use:   "char",
	Short: "Manage D&D characters (create, view, levelup, hp, spells, inventory)",
	Long: `The char command provides subcommands to manage your D&D characters. For guided creation, use the TUI with 'dnd tui'.
 
Use 'dnd char create <name>' to make a new character.
Use 'dnd char view <name>' to see a character's sheet.
Use 'dnd char levelup <name>' to advance a character's level.
Use 'dnd char hp <name> <action> <amount>' to manage HP.
Use 'dnd char spells <name> <action> <level> <amount>' to manage spell slots.
Use 'dnd char inventory <name> <action> <item>' to manage inventory.
Use 'dnd char condition <name> <action> <condition>' to manage conditions.
Use 'dnd char edit <name> <field> <value>' to edit character details.`,
}

func init() {
	RootCmd.AddCommand(charCmd)

	// Add 'create' subcommand
	var createCharCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new D&D character",
		Long:  `Guides you through the process of creating a new D&D character.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]

			// Check if data is loaded
			if len(data.AllSpecies) == 0 || len(data.AllClasses) == 0 || len(data.AllBackgrounds) == 0 {
				fmt.Printf("Hark! The ancient tomes are not loaded. Ensure the data submodule is initialized: git submodule update --init --recursive\n")
				return
			}

			// Check if character already exists
			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}
			if _, err := os.Stat(charFilePath); err == nil {
				fmt.Printf("Hark! A hero named '%s' already exists in the annals! Choose another name, adventurer.\n", charName)
				return
			}

			reader := bufio.NewReader(os.Stdin)

			// Species selection
			fmt.Printf("Choose a species (e.g., Human, Elf, Dwarf). Available: %s\n", getSpeciesNames())
			fmt.Print("Enter Species: ")
			speciesInput, _ := reader.ReadString('\n')
			speciesInput = strings.TrimSpace(speciesInput)
			if _, err := data.GetSpeciesByName(speciesInput); err != nil {
				fmt.Printf("Hark! That species is unknown to these lands. %v\n", err)
				return
			}

			// Class selection
			fmt.Printf("Choose a class (e.g., Fighter, Wizard, Rogue). Available: %s\n", getClassNames())
			fmt.Print("Enter Class: ")
			classInput, _ := reader.ReadString('\n')
			classInput = strings.TrimSpace(classInput)
			if _, err := data.GetClassByName(classInput); err != nil {
				fmt.Printf("Hark! That class is not in our teachings. %v\n", err)
				return
			}

			// Background selection
			fmt.Printf("Choose a background (e.g., Acolyte, Soldier, Criminal). Available: %s\n", getBackgroundNames())
			fmt.Print("Enter Background: ")
			backgroundInput, _ := reader.ReadString('\n')
			backgroundInput = strings.TrimSpace(backgroundInput)
			if _, err := data.GetBackgroundByName(backgroundInput); err != nil {
				fmt.Printf("Hark! That background is not etched in our lore. %v\n", err)
				return
			}

			// Create character with base stats (standard array)
			scores := [6]int{15, 14, 13, 12, 10, 8}
			newChar := character.NewCharacter(charName, speciesInput, classInput, backgroundInput, "", 1, scores[0], scores[1], scores[2], scores[3], scores[4], scores[5])

			// Apply racial, class, and background traits
			newChar.ApplyRacialTraits()
			newChar.ApplyClassTraits()
			newChar.ApplyBackgroundTraits()

			// Save character
			err = character.SaveCharacter(newChar, charFilePath)
			if err != nil {
				fmt.Printf("Hark! The scribe's quill faltered! Failed to save character: %v\n", err)
				return
			}

			fmt.Printf("\nVerily! Character '%s' created and saved to '%s'.\n", charName, charFilePath)
		},
	}
	charCmd.AddCommand(createCharCmd)

	// Add 'view' subcommand
	var viewCharCmd = &cobra.Command{
		Use:   "view [name]",
		Short: "View a D&D character's sheet",
		Long:  `Displays the details of a saved D&D character.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			fmt.Printf("\n--- Character Sheet: %s ---\n", char.Name)
			fmt.Printf("Species: %s\n", char.Species)
			fmt.Printf("Class: %s\n", char.Class)
			if char.Subclass != "" {
				fmt.Printf("Subclass: %s\n", char.Subclass)
			}
			fmt.Printf("Level: %d\n", char.Level)
			fmt.Printf("Background: %s\n", char.Background)
			fmt.Printf("Alignment: %s\n", char.Alignment)
			fmt.Printf("Experience: %d\n", char.Experience)
			fmt.Println("\n--- Ability Scores ---")
			fmt.Printf("STR: %d (%+d)\n", char.Strength, (char.Strength-10)/2)
			fmt.Printf("DEX: %d (%+d)\n", char.Dexterity, (char.Dexterity-10)/2)
			fmt.Printf("CON: %d (%+d)\n", char.Constitution, (char.Constitution-10)/2)
			fmt.Printf("INT: %d (%+d)\n", char.Intelligence, (char.Intelligence-10)/2)
			fmt.Printf("WIS: %d (%+d)\n", char.Wisdom, (char.Wisdom-10)/2)
			fmt.Printf("CHA: %d (%+d)\n", char.Charisma, (char.Charisma-10)/2)
			fmt.Println("\n--- Derived Stats ---")
			fmt.Printf("HP: %d/%d", char.CurrentHP, char.HitPoints)
			if char.TempHP > 0 {
				fmt.Printf(" (+%d temp)", char.TempHP)
			}
			fmt.Println()
			fmt.Printf("AC: %d\n", char.ArmorClass)
			fmt.Printf("Speed: %d ft.\n", char.Speed)
			fmt.Printf("Proficiency Bonus: +%d\n", char.ProficiencyBonus)
			fmt.Printf("Hit Dice: %s\n", char.HitDice)
			if char.Inspiration {
				fmt.Println("Inspiration: Yes")
			}
			if len(char.Conditions) > 0 {
				fmt.Printf("Conditions: %s\n", strings.Join(char.Conditions, ", "))
			}
			if len(char.Languages) > 0 {
				fmt.Printf("Languages: %s\n", strings.Join(char.Languages, ", "))
			}
			if len(char.ArmorProficiencies) > 0 {
				fmt.Printf("Armor Proficiencies: %s\n", strings.Join(char.ArmorProficiencies, ", "))
			}
			if len(char.WeaponProficiencies) > 0 {
				fmt.Printf("Weapon Proficiencies: %s\n", strings.Join(char.WeaponProficiencies, ", "))
			}
			if len(char.ToolProficiencies) > 0 {
				fmt.Printf("Tool Proficiencies: %s\n", strings.Join(char.ToolProficiencies, ", "))
			}
			if len(char.SkillProficiencies) > 0 {
				fmt.Printf("Skill Proficiencies: %s\n", strings.Join(char.SkillProficiencies, ", "))
			}
			if len(char.Features) > 0 {
				fmt.Printf("Features: %s\n", strings.Join(char.Features, ", "))
			}
			if char.SpellcastingAbility != "" {
				fmt.Printf("Spellcasting Ability: %s\n", char.SpellcastingAbility)
				if len(char.SpellSlots) > 0 {
					fmt.Print("Spell Slots: ")
					for level := 1; level <= 9; level++ {
						if count, ok := char.SpellSlots[level]; ok {
							used := char.UsedSpellSlots[level]
							fmt.Printf("%d: %d/%d ", level, count-used, count)
						}
					}
					fmt.Println()
				}
			}
			if len(char.Equipment) > 0 {
				fmt.Printf("Equipment: %s\n", strings.Join(char.Equipment, ", "))
			}
			fmt.Print("---------------------------\n")
		},
	}
	charCmd.AddCommand(viewCharCmd)

	// Add 'levelup' subcommand
	var levelUpCharCmd = &cobra.Command{
		Use:   "levelup [name]",
		Short: "Level up a D&D character",
		Long:  `Increases the level of a saved D&D character and updates basic stats.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			fmt.Printf("\nVerily! '%s' is now level %d.\n", char.Name, char.Level+1)
			char.LevelUp()

			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! The scribe's quill faltered! Failed to save character: %v\n", err)
				return
			}

			fmt.Printf("Character '%s' leveled up to %d! HP: %d, Proficiency Bonus: +%d.\n", char.Name, char.Level, char.HitPoints, char.ProficiencyBonus)
		},
	}
	charCmd.AddCommand(levelUpCharCmd)

	// Add 'hp' subcommand
	var hpCmd = &cobra.Command{
		Use:   "hp [name] [action] [amount]",
		Short: "Manage character HP",
		Long:  `Manage a character's hit points. Actions: damage, heal, set.`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]
			action := strings.ToLower(args[1])
			amountStr := args[2]

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			amount, err := strconv.Atoi(amountStr)
			if err != nil {
				fmt.Printf("Hark! '%s' is not a valid number. %v\n", amountStr, err)
				return
			}

			switch action {
			case "damage":
				char.CurrentHP -= amount
				if char.CurrentHP < 0 {
					char.CurrentHP = 0
				}
				fmt.Printf("Dealt %d damage to %s. Current HP: %d/%d\n", amount, charName, char.CurrentHP, char.HitPoints)
			case "heal":
				char.CurrentHP += amount
				if char.CurrentHP > char.HitPoints {
					char.CurrentHP = char.HitPoints
				}
				fmt.Printf("Healed %s for %d HP. Current HP: %d/%d\n", charName, amount, char.CurrentHP, char.HitPoints)
			case "set":
				char.CurrentHP = amount
				if char.CurrentHP > char.HitPoints {
					char.CurrentHP = char.HitPoints
				} else if char.CurrentHP < 0 {
					char.CurrentHP = 0
				}
				fmt.Printf("Set %s HP to %d. Current HP: %d/%d\n", charName, amount, char.CurrentHP, char.HitPoints)
			default:
				fmt.Printf("Hark! Unknown action '%s'. Use damage, heal, or set.\n", action)
				return
			}

			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! Failed to save character: %v\n", err)
			}
		},
	}
	charCmd.AddCommand(hpCmd)

	// Add 'spells' subcommand
	var spellsCmd = &cobra.Command{
		Use:   "spells [name] [action] [level] [amount]",
		Short: "Manage character spell slots",
		Long:  `Manage a character's spell slots. Actions: use, restore.`,
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]
			action := strings.ToLower(args[1])
			levelStr := args[2]
			amountStr := args[3]

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			level, err := strconv.Atoi(levelStr)
			if err != nil {
				fmt.Printf("Hark! '%s' is not a valid spell level. %v\n", levelStr, err)
				return
			}

			amount, err := strconv.Atoi(amountStr)
			if err != nil {
				fmt.Printf("Hark! '%s' is not a valid number. %v\n", amountStr, err)
				return
			}

			switch action {
			case "use":
				if char.UsedSpellSlots[level]+amount > char.SpellSlots[level] {
					fmt.Printf("Hark! Not enough spell slots at level %d. Available: %d, Used: %d\n", level, char.SpellSlots[level], char.UsedSpellSlots[level])
					return
				}
				char.UsedSpellSlots[level] += amount
				fmt.Printf("Used %d spell slot(s) at level %d for %s. Used: %d/%d\n", amount, level, charName, char.UsedSpellSlots[level], char.SpellSlots[level])
			case "restore":
				char.UsedSpellSlots[level] -= amount
				if char.UsedSpellSlots[level] < 0 {
					char.UsedSpellSlots[level] = 0
				}
				fmt.Printf("Restored %d spell slot(s) at level %d for %s. Used: %d/%d\n", amount, level, charName, char.UsedSpellSlots[level], char.SpellSlots[level])
			default:
				fmt.Printf("Hark! Unknown action '%s'. Use use or restore.\n", action)
				return
			}

			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! Failed to save character: %v\n", err)
			}
		},
	}
	charCmd.AddCommand(spellsCmd)

	// Add 'inventory' subcommand
	var inventoryCmd = &cobra.Command{
		Use:   "inventory [name] [action] [item]",
		Short: "Manage character inventory",
		Long:  `Manage a character's inventory. Actions: add, remove.`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]
			action := strings.ToLower(args[1])
			item := strings.Join(args[2:], " ")

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			switch action {
			case "add":
				char.Equipment = append(char.Equipment, item)
				fmt.Printf("Added '%s' to %s's inventory.\n", item, charName)
			case "remove":
				for i, it := range char.Equipment {
					if strings.EqualFold(it, item) {
						char.Equipment = append(char.Equipment[:i], char.Equipment[i+1:]...)
						fmt.Printf("Removed '%s' from %s's inventory.\n", item, charName)
						goto save
					}
				}
				fmt.Printf("Hark! '%s' not found in %s's inventory.\n", item, charName)
				return
			default:
				fmt.Printf("Hark! Unknown action '%s'. Use add or remove.\n", action)
				return
			}

		save:
			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! Failed to save character: %v\n", err)
			}
		},
	}
	charCmd.AddCommand(inventoryCmd)

	// Add 'condition' subcommand
	var conditionCmd = &cobra.Command{
		Use:   "condition [name] [action] [condition]",
		Short: "Manage character conditions",
		Long:  `Manage a character's conditions. Actions: add, remove.`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]
			action := strings.ToLower(args[1])
			condition := strings.Join(args[2:], " ")

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			switch action {
			case "add":
				char.Conditions = append(char.Conditions, condition)
				fmt.Printf("Added condition '%s' to %s.\n", condition, charName)
			case "remove":
				for i, cond := range char.Conditions {
					if strings.EqualFold(cond, condition) {
						char.Conditions = append(char.Conditions[:i], char.Conditions[i+1:]...)
						fmt.Printf("Removed condition '%s' from %s.\n", condition, charName)
						goto save
					}
				}
				fmt.Printf("Hark! '%s' not found in %s's conditions.\n", condition, charName)
				return
			default:
				fmt.Printf("Hark! Unknown action '%s'. Use add or remove.\n", action)
				return
			}

		save:
			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! Failed to save character: %v\n", err)
			}
		},
	}
	charCmd.AddCommand(conditionCmd)

	// Add 'edit' subcommand
	var editCmd = &cobra.Command{
		Use:   "edit [name] [field] [value]",
		Short: "Edit character details",
		Long:  `Edit a character's details. Fields: alignment, backstory.`,
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			charName := args[0]
			field := strings.ToLower(args[1])
			value := strings.Join(args[2:], " ")

			charFilePath, err := character.GetCharacterFilePath(charName)
			if err != nil {
				fmt.Printf("Hark! A parchment error: %v\n", err)
				return
			}

			char, err := character.LoadCharacter(charFilePath)
			if err != nil {
				fmt.Printf("Hark! The hero '%s' is not found in the archives! %v\n", charName, err)
				return
			}

			switch field {
			case "alignment":
				char.Alignment = value
				fmt.Printf("Set %s's alignment to '%s'.\n", charName, value)
			case "backstory":
				char.Backstory = value
				fmt.Printf("Set %s's backstory.\n", charName)
			default:
				fmt.Printf("Hark! Unknown field '%s'. Use alignment or backstory.\n", field)
				return
			}

			err = character.SaveCharacter(char, charFilePath)
			if err != nil {
				fmt.Printf("Hark! Failed to save character: %v\n", err)
			}
		},
	}
	charCmd.AddCommand(editCmd)
}

// Helper to get available species names
func getSpeciesNames() string {
	names := make([]string, len(data.AllSpecies))
	for i, s := range data.AllSpecies {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}

// Helper to get available background names
func getBackgroundNames() string {
	names := make([]string, len(data.AllBackgrounds))
	for i, b := range data.AllBackgrounds {
		names[i] = b.Name
	}
	return strings.Join(names, ", ")
}

// Helper to get available class names
func getClassNames() string {
	names := make([]string, len(data.AllClasses))
	for i, c := range data.AllClasses {
		names[i] = c.Name
	}
	return strings.Join(names, ", ")
}
