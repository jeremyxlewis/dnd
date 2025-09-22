# D&D CLI Companion

A command-line companion for Dungeons & Dragons players and Dungeon Masters, built with Go and Cobra.

## Features

*   **Dice Roller:** Roll dice using standard D&D notation (e.g., `2d6`, `1d20+5`), with support for advantage and disadvantage.
*   **Information Lookup:** Quickly look up details for spells, monsters, and items from the D&D 5e SRD.
*   **NPC Generator:** Generate random NPC names, species, and backgrounds.
*   **Interactive TUI:** Full-screen terminal interface for browsing and searching content with fuzzy matching.

## Installation

### Prerequisites

*   Go (version 1.16 or higher)
*   Git

### Steps

1.  **Clone the repository:**
    ```bash
    git clone --recurse-submodules https://github.com/YOUR_GITHUB_USERNAME/dnd-cli.git # Replace with your actual repo URL
    cd dnd-cli
    ```

2.  **The D&D data is already included** in the `data/` directory as a git submodule.

3.  **Build the application:**
    ```bash
    go mod tidy
    go build -o dnd .
    ```

## Usage

Run the `dnd` executable from the project root. Use `./dnd` if it's not in your PATH.

### Optional: Install Globally

To run `dnd` from anywhere without `./`:

- Move to a PATH directory: `sudo mv dnd /usr/local/bin/`
- Or add to your PATH: Add `export PATH="$PATH:/path/to/dnd-cli"` to `~/.bashrc` (replace with actual path).

Then use `dnd` or `dnd tui` directly.

### Dice Roller

Roll dice with standard notation:

```bash
dnd roll 1d20
dnd roll 2d6+3
```

Roll with advantage:

```bash
dnd roll 1d20 --advantage
dnd roll 1d20 -a
```

Roll with disadvantage:

```bash
dnd roll 1d20 --disadvantage
dnd roll 1d20 -d
```

### Spell Lookup

Look up a spell by name:

```bash
dnd spell "Fireball"
dnd spell "eldritch blast"
```

### Monster Lookup

Look up a monster by name:

```bash
dnd monster "Goblin"
dnd monster "Ancient Red Dragon"
```

### Item Lookup

Look up an item by name:

```bash
dnd item "Potion of Healing"
dnd item "Longsword"
```

### NPC Generator

Generate a random NPC:

```bash
dnd npc generate
# or simply
dnd npc
```

## TUI (Terminal User Interface)

Launch the interactive TUI for a more immersive experience:

```bash
dnd tui
```

### Features in TUI

- **Command Input:** Type commands at the prompt "What is thy command, adventurer?"
- **Fuzzy Search:** Type `spell`, `monster`, or `item` to browse lists with real-time filtering. Start typing to narrow down results.
- **Navigation:** Use arrow keys to scroll lists, Enter to select, Esc to cancel.
- **Scrolling:** Use ↑/↓ keys to scroll through long content in the output area.
- **Help:** Type `help` or `?` for available commands.
- **Quitting:** Press Esc or Ctrl+C to exit.

### Example TUI Workflow

1. Run `dnd tui`
2. Type `spell` and press Enter to browse spells.
3. Type "fire" to filter to fireball-related spells.
4. Use arrows to select, Enter to view details.
5. Press Esc to return to main prompt.
6. Type `roll 1d20` for dice rolls.

The TUI provides themed error messages and a clean, scrollable interface for all CLI features.

## Contributing



## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
