package tui

import (
	"fmt"
	"io"

	"dnd-cli/internal/data"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// fuzzyModel handles fuzzy finding and selection for various data categories.
type fuzzyModel struct {
	list  list.Model
	mode  string // The category being searched (e.g., "spell", "monster")
	width int
}

// listItem represents a list item for fuzzy finder.
type listItem struct {
	title string
}

func (i listItem) FilterValue() string { return i.title }

// customDelegate implements list.ItemDelegate for custom rendering of list items.
type customDelegate struct{}

func (d customDelegate) Height() int { return 1 }

func (d customDelegate) Spacing() int { return 0 }

func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d customDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	li, ok := item.(listItem)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprint(w, selectedStyle.Render("> "+li.title))
	} else {
		fmt.Fprint(w, unselectedStyle.Render("  "+li.title))
	}
}

// newFuzzyModel creates a fuzzy finder model for the given mode.
func newFuzzyModel(mode string) fuzzyModel {
	var titles []string

	switch mode {
	case "spell":
		titles = getUniqueTitles(data.AllSpells, func(s data.Spell) string { return s.Name })
	case "monster":
		titles = getUniqueTitles(data.AllMonsters, func(m data.Monster) string { return m.Name })
	case "item":
		titles = getUniqueTitles(data.AllItems, func(i data.Item) string { return i.Name })
	case "race":
		titles = getUniqueTitles(data.AllSpecies, func(s data.Species) string { return s.Name })
	case "background":
		titles = getUniqueTitles(data.AllBackgrounds, func(b data.Background) string { return b.Name })
	case "class":
		titles = getUniqueTitles(data.AllClasses, func(c data.Class) string { return c.Name })
	case "rules":
		titles = []string{"combat", "conditions", "ability checks", "initiative", "actions"}
	case "global":
		titles = append(titles, getUniqueTitles(data.AllSpells, func(s data.Spell) string { return "Spell: " + s.Name })...)
		titles = append(titles, getUniqueTitles(data.AllMonsters, func(m data.Monster) string { return "Monster: " + m.Name })...)
		titles = append(titles, getUniqueTitles(data.AllItems, func(i data.Item) string { return "Item: " + i.Name })...)
		titles = append(titles, getUniqueTitles(data.AllSpecies, func(s data.Species) string { return "Race: " + s.Name })...)
		titles = append(titles, getUniqueTitles(data.AllBackgrounds, func(b data.Background) string { return "Background: " + b.Name })...)
		titles = append(titles, getUniqueTitles(data.AllClasses, func(c data.Class) string { return "Class: " + c.Name })...)
		titles = append(titles, "Rules: combat", "Rules: conditions", "Rules: ability checks", "Rules: initiative", "Rules: actions")
	}

	items := createListItems(titles)
	l := list.New(items, customDelegate{}, DefaultWidth, DefaultHeight-ListHeightPadding) // initial size
	l.Title = "Select " + mode
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	return fuzzyModel{list: l, mode: mode, width: DefaultWidth}
}

func (m fuzzyModel) Init() tea.Cmd {
	return nil
}

func (m fuzzyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		height := msg.Height - ListHeightPadding
		if height < 10 {
			height = 10
		}
		m.list.SetSize(msg.Width, height)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			selected := m.list.SelectedItem()
			if selected != nil {
				name := selected.(listItem).title
				return m, func() tea.Msg { return selectedMsg{mode: m.mode, name: name} }
			}
		} else if msg.Type == tea.KeyEsc {
			return m, func() tea.Msg { return switchModeMsg{"main"} }
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m fuzzyModel) View() string {
	return m.list.View()
}
