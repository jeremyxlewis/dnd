package tui

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// getRandomMessage returns a random message from the given slice.
func getRandomMessage(messages []string) string {
	return messages[rand.Intn(len(messages))]
}

// getRandomErrorMessage returns a random error message.
func getRandomErrorMessage() string {
	return getRandomMessage(errorMessages)
}

// getRandomSpellErrorMessage returns a random spell error message.
func getRandomSpellErrorMessage(name string) string {
	return fmt.Sprintf(getRandomMessage(spellErrorMessages), name)
}

// getRandomMonsterErrorMessage returns a random monster error message.
func getRandomMonsterErrorMessage(name string) string {
	return fmt.Sprintf(getRandomMessage(monsterErrorMessages), name)
}

// getRandomItemErrorMessage returns a random item error message.
func getRandomItemErrorMessage(name string) string {
	return fmt.Sprintf(getRandomMessage(itemErrorMessages), name)
}

// getRandomSpeciesErrorMessage returns a random species error message.
func getRandomSpeciesErrorMessage(name string) string {
	return fmt.Sprintf(getRandomMessage(speciesErrorMessages), name)
}

// getRandomBackgroundErrorMessage returns a random background error message.
func getRandomBackgroundErrorMessage(name string) string {
	return fmt.Sprintf(getRandomMessage(backgroundErrorMessages), name)
}

// getRandomClassErrorMessage returns a random class error message.
func getRandomClassErrorMessage(name string) string {
	return fmt.Sprintf(getRandomMessage(classErrorMessages), name)
}

// getRandomPrompt returns a random prompt message.
func getRandomPrompt() string {
	return getRandomMessage(prompts)
}

// createListItems creates a slice of list.Item from a slice of strings.
func createListItems(titles []string) []list.Item {
	items := make([]list.Item, len(titles))
	for i, title := range titles {
		items[i] = listItem{title: title}
	}
	return items
}

// getUniqueTitles extracts unique titles from a slice of any type.
func getUniqueTitles[T any](slice []T, getTitle func(T) string) []string {
	seen := make(map[string]bool)
	var titles []string
	for _, item := range slice {
		title := getTitle(item)
		if !seen[title] {
			titles = append(titles, title)
			seen[title] = true
		}
	}
	return titles
}

// formatDescription formats long descriptions for better readability.
func formatDescription(desc string) string {
	// Clean up existing formatting first
	desc = strings.ReplaceAll(desc, "\n", " ")

	// Add line breaks after complete sentences
	desc = strings.ReplaceAll(desc, ". ", ".\n\n")

	// Add line breaks before ability names (common monster ability pattern)
	re := regexp.MustCompile(`([.!?]) ([A-Z][a-z]+(?: [A-Z][a-z]+)*:)`)
	desc = re.ReplaceAllString(desc, "$1\n\n$2")

	// Add line breaks after sentence fragments followed by uppercase letters (not ending in colon)
	re2 := regexp.MustCompile(`([.!?]) ([A-Z][a-z]*(?: [a-z]+)*)`)
	desc = re2.ReplaceAllString(desc, "$1\n\n$2")

	// Add line breaks before common section headers and abilities
	headers := []string{
		"Hit Points", "Proficiencies", "Armor", "Weapons", "Tools", "Saving Throws", "Skills", "Equipment",
		"Spellcasting", "Cantrips", "Spellbook", "Preparing and Casting Spells", "Arcane Recovery",
		"Arcane Tradition", "Ability Score Improvement", "Creating a", "Quick Build", "The Wizard Level",
		"Class Features", "Creature Type", "Size", "Speed", "Flight", "Talons", "Wind Caller",
		"Dive Attack", "Melee Weapon Attack", "Ranged Weapon Attack", "Summon Air Elemental", "Hit:",
		"Description", "Personality Trait", "Ideal", "Bond", "Flaw", "Backstory",
	}
	for _, h := range headers {
		// Add line breaks for headers with various patterns
		desc = strings.ReplaceAll(desc, " "+h, "\n\n"+h)
		desc = strings.ReplaceAll(desc, h+":", "\n\n"+h+":")
		if strings.HasSuffix(h, ":") {
			desc = strings.ReplaceAll(desc, strings.TrimSuffix(h, ":"), "\n\n"+h)
		}
	}

	// Clean up extra whitespace and ensure proper line breaks
	lines := strings.Split(desc, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			// Don't start a line with just punctuation
			if !strings.HasPrefix(trimmed, ".") && !strings.HasPrefix(trimmed, "!") && !strings.HasPrefix(trimmed, "?") {
				cleanLines = append(cleanLines, trimmed)
			}
		}
	}

	return strings.Join(cleanLines, "\n")
}

// renderPropertiesTable renders a simple table for key-value pairs.
func renderPropertiesTable(properties map[string]interface{}, order []string) string {
	var rows []string
	for _, key := range order {
		if v, ok := properties[key]; ok {
			keyStyled := lipgloss.NewStyle().Bold(true).Render(key + ":")
			valueStyled := outputStyle.Render(fmt.Sprintf("%v", v))
			row := lipgloss.JoinHorizontal(lipgloss.Top, keyStyled, valueStyled)
			rows = append(rows, row)
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
