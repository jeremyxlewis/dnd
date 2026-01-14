# D&D TUI Companion

An interactive terminal user interface for Dungeons & Dragons players and Dungeon Masters, built with Go and Bubble Tea. This is a **TUI-only** application - no CLI commands or arguments.

## Features

*   **Interactive TUI:** Full-screen terminal interface for browsing and searching content with fuzzy matching.
*   **Spell Browser:** Browse through all D&D 5e spells with detailed information including level, school, casting time, range, components, duration, and descriptions.
*   **Monster Compendium:** Complete monster stat blocks with abilities, actions, skills, and challenge ratings.
*   **Rules Reference:** Look up D&D 5e rules covering combat, conditions, ability checks, and more.
*   **Dice Rolling:** Interactive dice rolling with standard D&D notation and visual feedback.
*   **Fuzzy Search:** Real-time filtering across all content types - just start typing to narrow results.
*   **Keyboard Navigation:** Full keyboard support with arrow keys, j/k navigation, search shortcuts, and intuitive controls.

## Installation

### Prerequisites

*   Go (version 1.25 or higher)
*   Git

### Steps

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/jeremyxlewis/dnd.git
    cd dnd
    ```

2.  **Get D&D 5e SRD data (required):**
    ```bash
    # Clone the official D&D 5e SRD repository
    git clone https://github.com/adrjian/5e-srd.git dnd-5e-srd
    
    # Or download the 5esrd.json file directly to the dnd-5e-srd/ directory
    mkdir -p dnd-5e-srd
    # Download from: https://github.com/adrjian/5e-srd/blob/master/5esrd.json
    ```
    
    **IMPORTANT:** The SRD data file (`dnd-5e-srd/5esrd.json`) is **required** for the application to run. The old `data/` JSON files are no longer used.

3.  **Build the application:**
    ```bash
    go mod tidy
    go build -o dnd .
    ```

## Usage

3.  **Run the application:**
    ```bash
    ./dnd
    ```
    
    The application will launch directly into the interactive terminal interface. Make sure you have the SRD data available at `dnd-5e-srd/5esrd.json`.

### Optional: Install Globally

To run `dnd` from anywhere:

- Move to a PATH directory: `sudo mv dnd /usr/local/bin/`
- Or add to your PATH: Add `export PATH="$PATH:/path/to/dnd"` to `~/.bashrc` (replace with actual path).

## TUI Features and Usage

### Main Interface

When you launch `./dnd`, you'll see:

- **Welcome Screen:** Interactive command interface with fuzzy search
- **Content Browser:** Access to spells, monsters, items, species, backgrounds, classes, and rules
- **Dice Roller:** Interactive dice interface with visual feedback
- **Help System:** Built-in help and command reference

### Navigation Controls

- **Arrow Keys (↑↓)** or **j/k**: Navigate through lists and menus
- **Enter**: Select an item or confirm an action
- **Esc**: Go back to previous screen or exit
- **/**: Search within current content list
- **Ctrl+C**: Quick exit from any screen

### Content Browsing

1. **Spell Browser:**
   - Type `spell` at the main prompt
   - Browse complete spell list with level, school, and casting time
   - Real-time fuzzy filtering - just start typing!
   - Press Enter for full details: range, components, duration, description

2. **Monster Compendium:**
   - Type `monster` at the main prompt  
   - View complete stat blocks with AC, HP, abilities, actions
   - Includes challenge ratings, skills, and special abilities

3. **Other Content:**
   - `species`: Browse available races and racial traits
   - `backgrounds`: View background features and proficiencies  
   - `classes`: Browse class information and features
   - `rules`: Look up D&D 5e rules and mechanics

### Rules Reference

Look up D&D 5e rules by typing `rules` at the main prompt:

- **Combat:** Attack rolls, damage, critical hits, and combat flow
- **Conditions:** All condition effects and duration rules  
- **Ability Checks:** How to make ability checks and saving throws
- **Skill Usage:** Guidelines for using skills in various situations

Type `rules combat`, `rules conditions`, or browse the complete ruleset.

### Dice Rolling

Type `roll` followed by dice notation:

- `roll 1d20` - Single d20 roll
- `roll 2d6+3` - Multiple dice with modifier
- `roll 4d6` - Multiple dice
- Visual feedback shows individual rolls and total

### Search and Filtering

All content browsers support:
- **Real-time filtering:** Start typing to instantly filter results
- **Case-insensitive:** Works regardless of capitalization
- **Partial matching:** "fire" finds "Fireball", "Fire Bolt", etc.
- **Fuzzy matching:** Finds close matches for typos

### Example Workflow

```bash
./dnd

# At main prompt:
> spell          # Enter spell browser
> fire           # Type to filter for fire spells
[↓] Fireball    # Select and press Enter
[Full spell details displayed]
[Esc]            # Go back to main
> roll 1d20+5    # Roll dice with advantage
[Roll result shown]
```

The TUI provides an immersive, terminal-native experience for all D&D 5e content with fuzzy search, themed styling, and intuitive keyboard controls. Perfect for quick rule lookups, monster stats, spell references, and dice rolling during your D&D sessions.

## Contributing

### Future Improvements

- [x] Replace deprecated `ioutil` functions with `os` equivalents in character handling
- [x] Optimize random seeding to avoid reseeding on every operation
- [ ] Implement full character creation and management features
- [x] Expand test coverage for core functionality
- [x] Refactor error message generation to reduce code duplication
- [x] Remove legacy CLI components and transition to pure TUI architecture
- [x] Clean up unused data directories and scripts
- [x] Update documentation to reflect current TUI-only structure
- [x] Remove outdated character creation references from README

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
