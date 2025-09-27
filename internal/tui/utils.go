package tui

import (
	"fmt"
	"math/rand"
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
	// Add line breaks after sentences
	desc = strings.ReplaceAll(desc, ". ", ".\n\n")
	// Add line breaks before common section headers
	headers := []string{
		"Hit Points", "Proficiencies", "Armor", "Weapons", "Tools", "Saving Throws", "Skills", "Equipment",
		"Spellcasting", "Cantrips", "Spellbook", "Preparing and Casting Spells", "Arcane Recovery",
		"Arcane Tradition", "Ability Score Improvement", "Creating a", "Quick Build", "The Wizard Level",
		"Class Features", "Creature Type", "Size", "Speed", "Flight", "Talons", "Wind Caller",
		"Description", "Personality Trait", "Ideal", "Bond", "Flaw", "Backstory",
	}
	for _, h := range headers {
		desc = strings.ReplaceAll(desc, h, "\n\n"+h)
	}
	return desc
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
