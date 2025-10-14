# 🏆 Sports Scores Terminal

A beautiful and functional terminal interface for checking live sports scores and games. Built with Go and inspired by the lazygit UI.

![Sports Scores Demo](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## ✨ Features

- 🎨 **Beautiful TUI** - Modern terminal interface with colors and styling
- ⚡ **Live Updates** - Auto-refreshing scores for live games
- 🏅 **Multiple Sports** - Support for NFL, NBA, MLB, NHL, Soccer, and more
- 🌍 **International Leagues** - Premier League, La Liga, Serie A, Bundesliga, Champions League
- 📊 **Real-time Data** - Powered by ESPN's public API
- ⌨️ **Keyboard Navigation** - Vim-style keybindings (hjkl) and arrow keys
- 🎯 **Easy to Use** - Intuitive navigation between sports, leagues, and games

## 📦 Installation

### Homebrew (macOS/Linux) - Recommended

```bash
brew install elliota43/tap/sportsterminal
```

### Go Install

```bash
go install github.com/elliota43/sportsterminal@latest
```

### Build from Source

**Prerequisites:** Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/elliota43/sportsterminal.git
cd sportsterminal

# Install dependencies
go mod download

# Build the application
go build -o sportsterminal

# Run the application
./sportsterminal
```

## 🚀 Usage

Simply run the application:

```bash
./sportsterminal
```

### Keyboard Controls

#### General Navigation
- `↑/k` - Move cursor up
- `↓/j` - Move cursor down
- `Enter/→/l` - Select item / Go forward
- `Esc/Backspace` - Go back to previous screen
- `q/Ctrl+C` - Quit application

#### Games View
- `r` - Manually refresh scores
- Auto-refresh every 30 seconds for live games

## 🎮 Sports & Leagues Supported

### 🏈 Football
- NFL (National Football League)
- College Football

### 🏀 Basketball
- NBA (National Basketball Association)
- WNBA (Women's National Basketball Association)
- NCAA Men's Basketball
- NCAA Women's Basketball

### ⚾ Baseball
- MLB (Major League Baseball)
- College Baseball

### 🏒 Hockey
- NHL (National Hockey League)

### ⚽ Soccer
- Premier League (England)
- La Liga (Spain)
- Serie A (Italy)
- Bundesliga (Germany)
- MLS (Major League Soccer)
- UEFA Champions League

## 🛠️ Technology Stack

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Style definitions for terminal layouts
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - Common TUI components
- **ESPN API** - Real-time sports data

## 📁 Project Structure

```
sportsterminal/
├── main.go           # Application entry point
├── api/
│   └── sports.go     # ESPN API client and data models
├── ui/
│   └── model.go      # TUI logic and rendering
├── go.mod            # Go module dependencies
└── README.md         # This file
```

## 🎨 Screenshots

The application features:
- A clean, modern interface with color-coded elements
- Live game indicators (🔴 LIVE)
- Real-time score updates
- Game status and venue information
- Smooth navigation between views

## 🚀 Publishing to Homebrew

Want to distribute this via Homebrew? See [PUBLISH.md](PUBLISH.md) for a quick guide or [HOMEBREW.md](HOMEBREW.md) for detailed instructions.

## 🤝 Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests
- Add support for more sports/leagues

## 📝 License

MIT License - feel free to use this project however you'd like!

## 🙏 Acknowledgments

- Inspired by [lazygit](https://github.com/jesseduffield/lazygit)
- Built with [Charm](https://charm.sh/) libraries
- Sports data provided by ESPN

## 📮 Contact

For questions or feedback, please open an issue on GitHub.

---

**Enjoy checking your sports scores! 🏆**

