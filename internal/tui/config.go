package tui

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme for the TUI.
type Theme struct {
	FocusedColor    string `json:"focused_color"`
	BlurColor       string `json:"blur_color"`
	CursorColor     string `json:"cursor_color"`
	PromptColor     string `json:"prompt_color"`
	QuitColor       string `json:"quit_color"`
	OutputColor     string `json:"output_color"`
	ErrorColor      string `json:"error_color"`
	SelectedColor   string `json:"selected_color"`
	UnselectedColor string `json:"unselected_color"`
	InfoCardColor   string `json:"info_card_color"`
	RollColor       string `json:"roll_color"`
	HeaderColor     string `json:"header_color"`
}

// Config holds TUI configuration.
type Config struct {
	Theme Theme `json:"theme"`
}

// DefaultTheme returns the default theme.
func DefaultTheme() Theme {
	return Theme{
		FocusedColor:    "#FFB6C1", // Pastel Pink
		BlurColor:       "#E6E6FA", // Lavender
		CursorColor:     "#FFB6C1", // Pastel Pink
		PromptColor:     "#A8DADC", // Pastel Mint
		QuitColor:       "#D3D3D3", // Light Gray
		OutputColor:     "#F8F9FA", // Off-White
		ErrorColor:      "#F8BBD9", // Pastel Rose
		SelectedColor:   "#FFB6C1", // Pastel Pink
		UnselectedColor: "#D3D3D3", // Light Gray
		InfoCardColor:   "#B19CD9", // Pastel Purple
		RollColor:       "#A8DADC", // Pastel Mint
		HeaderColor:     "#FFB6C1", // Pastel Pink
	}
}

// LoadConfig loads the config from file, or returns defaults.
func LoadConfig() Config {
	configDir := filepath.Join(os.Getenv("HOME"), ".dnd-cli")
	configPath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Config{Theme: DefaultTheme()}
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{Theme: DefaultTheme()}
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{Theme: DefaultTheme()}
	}

	return config
}

// SaveConfig saves the config to file.
func SaveConfig(config Config) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".dnd-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(config)
}

// ApplyTheme applies the theme to global styles.
func ApplyTheme(theme Theme) {
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FocusedColor))
	blurStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.BlurColor))
	cursorStyle = focusedStyle
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.PromptColor)).Bold(true)
	quitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.QuitColor)).Italic(true)
	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.OutputColor)).Padding(0, 1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ErrorColor)).Padding(0, 1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SelectedColor)).Bold(true)
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.UnselectedColor))
	infoCardStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).BorderForeground(lipgloss.Color(theme.InfoCardColor))
	rollStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderForeground(lipgloss.Color(theme.RollColor))
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.HeaderColor)).Background(lipgloss.Color(theme.HeaderColor)).Padding(0, 1).Bold(true)
	viewStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.OutputColor))
}
