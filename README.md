# D&D TUI Companion

An interactive terminal user interface for Dungeons & Dragons players and Dungeon Masters, built with Go and Bubble Tea. Features comprehensive Compendium XML data integration with rich D&D 5e content.

## Features

*   **Rich Compendium Data:** Comprehensive D&D 5e content from Compendium XML files with spells, monsters, items, races, backgrounds, and classes
*   **Multiple Data Sources:** Support for various Compendium editions (Complete, Core, WotC-only, etc.) with easy switching
*   **Interactive TUI:** Full-screen terminal interface for browsing and searching content with fuzzy matching
*   **Spell Browser:** Browse through all D&D 5e spells with detailed information including level, school, casting time, range, components, duration, and descriptions
*   **Monster Compendium:** Complete monster stat blocks with abilities, actions, skills, and challenge ratings
*   **Items Database:** Comprehensive item catalog with weapons, armor, and equipment details
*   **Character Creation Tools:** Race and background information with ability bonuses and traits
*   **Class Information:** Complete class details including proficiencies, equipment, and features
*   **Rules Reference:** Look up D&D 5e rules covering combat, conditions, ability checks, and more
*   **Dice Rolling:** Interactive dice rolling with standard D&D notation and visual feedback
*   **Fuzzy Search:** Real-time filtering across all content types - just start typing to narrow results
*   **Keyboard Navigation:** Full keyboard support with arrow keys, j/k navigation, search shortcuts, and intuitive controls

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

2.  **Build the application:**
    ```bash
    go mod tidy
    go build -o dnd .
    ```

That's it! The application includes comprehensive Compendium XML data files, so no additional data downloads are required. The repository comes with multiple Compendium editions pre-installed in the `Compendium/` directory.

## Usage

### Running the Application

1.  **Run with default Compendium:**
    ```bash
    ./dnd
    ```
    Uses `Complete_Compendium.xml` by default (most comprehensive content).

2.  **Run with specific Compendium file:**
    ```bash
    ./dnd --compendium Core_Rulebooks.xml
    ```
    Choose from available Compendium files in the `Compendium/` directory.

3.  **Run with custom Compendium directory:**
    ```bash
    ./dnd --path /path/to/custom/Compendium
    ```

### Available Compendium Files

The application includes multiple Compendium editions:

- **Complete_Compendium.xml** (default) - All content types with maximum variety
- **Core_Rulebooks.xml** - Official Wizards of the Coast content only
- **WotC_only.xml** - Pure WotC content without third-party material
- **WotC+PartneredContent.xml** - Official plus partnered content
- **System_Reference_Document.xml** - Basic SRD content

Launch the application and type `compendium list` to see all available files, or `compendium` to see which file is currently loaded.

### Optional: Install Globally

To run `dnd` from anywhere:

- Move to a PATH directory: `sudo mv dnd /usr/local/bin/`
- Or add to your PATH: Add `export PATH="$PATH:/path/to/dnd"` to `~/.bashrc` (replace with actual path).

### Command Line Options

- `--compendium <filename>`: Load specific Compendium file (e.g., `Core_Rulebooks.xml`)
- `--path <directory>`: Use custom Compendium directory instead of default `./Compendium/`

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

1.  **Spell Browser:**
   - Type `spell` at the main prompt
   - Browse complete spell list with level, school, and casting time
   - Real-time fuzzy filtering - just start typing!
   - Press Enter for full details: range, components, duration, description

2.  **Monster Compendium:**
   - Type `monster` at the main prompt  
   - View complete stat blocks with AC, HP, abilities, actions
   - Includes challenge ratings, skills, and special abilities

3.  **Items Database:**
   - Type `item` at the main prompt
   - Browse weapons, armor, and equipment with detailed stats
   - View damage types, properties, values, and weights

4.  **Character Creation Tools:**
   - `species`: Browse available races with ability bonuses and traits
   - `backgrounds`: View background features and proficiencies  
   - `classes`: Browse class information with proficiencies and features

5.  **Rules Reference:**
   - Type `rules` at the main prompt
   - Look up D&D 5e rules and mechanics

6.  **Compendium Management:**
   - Type `compendium` to see currently loaded file
   - Type `compendium list` to see all available Compendium files

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

# At the main prompt:
> spell fireball    # Look up specific spell directly
[Full spell details displayed]

> spell            # Opens spell browser/fuzzy finder
[Spell list appears]
> fire             # Type to filter spells in the list
[↓] Fireball       # Navigate with arrows and press Enter
[Full spell details displayed]
[Esc]              # Go back to main prompt

> search fire      # Global search across all content types
[Results: Spells, Monsters, Items...]
[Select any result for details]

> roll 1d20+5     # Roll dice
[Roll result: 18]   # Individual rolls and total shown
```

The TUI provides an immersive, terminal-native experience for comprehensive D&D 5e content from Compendium XML data. With rich spell descriptions, detailed monster stat blocks, complete item databases, and character creation tools, it's perfect for quick rule lookups, monster stats, spell references, and dice rolling during your D&D sessions.

## Contributing

### Recent Updates

- [x] **Integrate Compendium XML data source** - Replace JSON SRD with comprehensive XML parser
- [x] **Multiple Compendium support** - Support for Complete, Core, WotC, and other editions
- [x] **Rich data parsing** - Proper format conversion for schools, damage types, and properties
- [x] **TUI Compendium commands** - Add compendium file listing and status display
- [x] **Enhanced content coverage** - Access to extensive spell, monster, item, and race databases
- [x] **Flexible data loading** - Command-line options for Compendium file and directory selection
- [x] **Backward compatibility** - Maintain all existing TUI functionality with new data source

### Future Improvements

- [ ] Implement full character creation and management features
- [ ] Add campaign management tools
- [ ] Implement custom note-taking and bookmarking
- [ ] Add more Compendium file filters and categories
- [ ] Expand test coverage for new Compendium integration

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.