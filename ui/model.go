package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elliota43/sportsterminal/api"
)

type viewState int

const (
	sportView viewState = iota
	leagueView
	gamesView
)

type Model struct {
	state          viewState
	selectedSport  *api.Sport
	selectedLeague *api.League
	games          []api.Game
	sportCursor    int
	leagueCursor   int
	gameCursor     int
	width          int
	height         int
	loading        bool
	err            error
	lastUpdate     time.Time
	autoRefresh    bool
}

type gamesLoadedMsg struct {
	games []api.Game
	err   error
}

type tickMsg time.Time

func NewModel() Model {
	return Model{
		state:       sportView,
		autoRefresh: true,
		lastUpdate:  time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func loadGamesCmd(sport, league string) tea.Cmd {
	return func() tea.Msg {
		games, err := api.GetGames(sport, league)
		return gamesLoadedMsg{games: games, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "backspace":
			switch m.state {
			case leagueView:
				m.state = sportView
				m.leagueCursor = 0
			case gamesView:
				m.state = leagueView
				m.gameCursor = 0
				m.games = nil
			}
			return m, nil

		case "r":
			// Manual refresh
			if m.state == gamesView && m.selectedSport != nil && m.selectedLeague != nil {
				m.loading = true
				return m, loadGamesCmd(m.selectedSport.ID, m.selectedLeague.ID)
			}
			return m, nil

		case "up", "k":
			switch m.state {
			case sportView:
				if m.sportCursor > 0 {
					m.sportCursor--
				}
			case leagueView:
				if m.leagueCursor > 0 {
					m.leagueCursor--
				}
			case gamesView:
				if m.gameCursor > 0 {
					m.gameCursor--
				}
			}
			return m, nil

		case "down", "j":
			switch m.state {
			case sportView:
				if m.sportCursor < len(api.AvailableSports)-1 {
					m.sportCursor++
				}
			case leagueView:
				if m.selectedSport != nil && m.leagueCursor < len(m.selectedSport.Leagues)-1 {
					m.leagueCursor++
				}
			case gamesView:
				if m.gameCursor < len(m.games)-1 {
					m.gameCursor++
				}
			}
			return m, nil

		case "enter", "right", "l":
			switch m.state {
			case sportView:
				if m.sportCursor < len(api.AvailableSports) {
					m.selectedSport = &api.AvailableSports[m.sportCursor]
					m.state = leagueView
					m.leagueCursor = 0
				}
			case leagueView:
				if m.selectedSport != nil && m.leagueCursor < len(m.selectedSport.Leagues) {
					m.selectedLeague = &m.selectedSport.Leagues[m.leagueCursor]
					m.state = gamesView
					m.loading = true
					m.gameCursor = 0
					return m, loadGamesCmd(m.selectedSport.ID, m.selectedLeague.ID)
				}
			}
			return m, nil
		}

	case gamesLoadedMsg:
		m.loading = false
		m.games = msg.games
		m.err = msg.err
		m.lastUpdate = time.Now()
		return m, nil

	case tickMsg:
		// Auto-refresh live games
		if m.autoRefresh && m.state == gamesView && m.selectedSport != nil && m.selectedLeague != nil {
			hasLiveGames := false
			for _, game := range m.games {
				if game.IsLive {
					hasLiveGames = true
					break
				}
			}
			if hasLiveGames {
				return m, tea.Batch(
					loadGamesCmd(m.selectedSport.ID, m.selectedLeague.ID),
					tickCmd(),
				)
			}
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string

	switch m.state {
	case sportView:
		content = m.renderSportView()
	case leagueView:
		content = m.renderLeagueView()
	case gamesView:
		content = m.renderGamesView()
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		content,
	)
}

func (m Model) renderSportView() string {
	title := titleStyle.Render("🏆 Sports Scores")
	subtitle := subtitleStyle.Render("Select a sport")

	var items string
	for i, sport := range api.AvailableSports {
		cursor := "  "
		style := itemStyle
		if i == m.sportCursor {
			cursor = "❯ "
			style = selectedItemStyle
		}
		icon := getSportIcon(sport.ID)
		items += style.Render(fmt.Sprintf("%s%s %s", cursor, icon, sport.Name)) + "\n"
	}

	help := helpStyle.Render("↑/k up • ↓/j down • enter select • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		items,
		"",
		help,
	)
}

func (m Model) renderLeagueView() string {
	if m.selectedSport == nil {
		return "No sport selected"
	}

	title := titleStyle.Render(fmt.Sprintf("🏆 %s", m.selectedSport.Name))
	subtitle := subtitleStyle.Render("Select a league")

	var items string
	for i, league := range m.selectedSport.Leagues {
		cursor := "  "
		style := itemStyle
		if i == m.leagueCursor {
			cursor = "❯ "
			style = selectedItemStyle
		}
		items += style.Render(fmt.Sprintf("%s%s", cursor, league.Name)) + "\n"
	}

	help := helpStyle.Render("↑/k up • ↓/j down • enter select • esc back • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		items,
		"",
		help,
	)
}

func (m Model) renderGamesView() string {
	if m.selectedLeague == nil {
		return "No league selected"
	}

	title := titleStyle.Render(fmt.Sprintf("🏆 %s - %s", m.selectedSport.Name, m.selectedLeague.Name))
	
	var statusText string
	if m.loading {
		statusText = subtitleStyle.Render("Loading games...")
	} else {
		lastUpdate := m.lastUpdate.Format("3:04 PM")
		statusText = subtitleStyle.Render(fmt.Sprintf("Last updated: %s", lastUpdate))
	}

	if m.err != nil {
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		help := helpStyle.Render("r refresh • esc back • q quit")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", errorMsg, "", help)
	}

	if len(m.games) == 0 && !m.loading {
		noGames := itemStyle.Render("No games scheduled")
		help := helpStyle.Render("r refresh • esc back • q quit")
		return lipgloss.JoinVertical(lipgloss.Left, title, statusText, "", noGames, "", help)
	}

	var items string
	for i, game := range m.games {
		cursor := "  "
		if i == m.gameCursor {
			cursor = "❯ "
		}
		items += cursor + m.renderGame(game, i == m.gameCursor) + "\n\n"
	}

	help := helpStyle.Render("↑/k up • ↓/j down • r refresh • esc back • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		statusText,
		"",
		items,
		help,
	)
}

func (m Model) renderGame(game api.Game, selected bool) string {
	boxStyle := gameBoxStyle
	if selected {
		boxStyle = selectedGameBoxStyle
	}

	status := game.Status
	if game.IsLive {
		status = liveStyle.Render("🔴 LIVE - " + status)
	} else {
		status = statusStyle.Render(status)
	}

	awayScore := game.AwayTeam.Score
	homeScore := game.HomeTeam.Score
	
	if awayScore == "" {
		awayScore = "-"
	}
	if homeScore == "" {
		homeScore = "-"
	}

	gameTime := game.Date.Local().Format("Mon Jan 2, 3:04 PM")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		status,
		"",
		teamStyle.Render(fmt.Sprintf("%-30s %3s", game.AwayTeam.Name, awayScore)),
		teamStyle.Render(fmt.Sprintf("%-30s %3s", game.HomeTeam.Name, homeScore)),
		"",
		venueStyle.Render(fmt.Sprintf("📍 %s", game.Venue)),
		venueStyle.Render(fmt.Sprintf("🕐 %s", gameTime)),
	)

	return boxStyle.Render(content)
}

func getSportIcon(sportID string) string {
	icons := map[string]string{
		"football":   "🏈",
		"basketball": "🏀",
		"baseball":   "⚾",
		"hockey":     "🏒",
		"soccer":     "⚽",
	}
	if icon, ok := icons[sportID]; ok {
		return icon
	}
	return "🏃"
}

var (
	primaryColor = lipgloss.Color("#7C3AED")
	accentColor  = lipgloss.Color("#F59E0B")
	liveColor    = lipgloss.Color("#EF4444")
	textColor    = lipgloss.Color("#E5E7EB")
	dimColor     = lipgloss.Color("#9CA3AF")
	bgColor      = lipgloss.Color("#1F2937")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(1, 2).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Padding(0, 2)

	itemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor).
				Padding(0, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Padding(1, 2)

	gameBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dimColor).
			Padding(1, 2).
			Width(60)

	selectedGameBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				Width(60)

	teamStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	liveStyle = lipgloss.NewStyle().
			Foreground(liveColor).
			Bold(true)

	venueStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(liveColor).
			Padding(0, 2)
)

