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
	gameDetailView
)

type Model struct {
	state              viewState
	selectedSport      *api.Sport
	selectedLeague     *api.League
	games              []api.Game
	selectedGameDetail *api.GameDetail
	sportCursor        int
	leagueCursor       int
	gameCursor         int
	gameScrollOffset   int
	detailScrollOffset int
	width              int
	height             int
	loading            bool
	loadingDetail      bool
	err                error
	lastUpdate         time.Time
	autoRefresh        bool
}

type gamesLoadedMsg struct {
	games []api.Game
	err   error
}

type gameDetailLoadedMsg struct {
	detail *api.GameDetail
	err    error
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

func loadGameDetailCmd(sport, league, eventID string) tea.Cmd {
	return func() tea.Msg {
		detail, err := api.GetGameDetail(sport, league, eventID)
		return gameDetailLoadedMsg{detail: detail, err: err}
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
				m.gameScrollOffset = 0
				m.games = nil
			case gameDetailView:
				m.state = gamesView
				m.selectedGameDetail = nil
				m.detailScrollOffset = 0
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
					// Scroll up if cursor moves above visible area
					if m.gameCursor < m.gameScrollOffset {
						m.gameScrollOffset = m.gameCursor
					}
				}
			case gameDetailView:
				if m.detailScrollOffset > 0 {
					m.detailScrollOffset--
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
					// Calculate visible games and scroll down if needed
					visibleGames := m.calculateVisibleGames()
					if m.gameCursor >= m.gameScrollOffset+visibleGames {
						m.gameScrollOffset = m.gameCursor - visibleGames + 1
					}
				}
			case gameDetailView:
				m.detailScrollOffset++
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
					m.gameScrollOffset = 0
					return m, loadGamesCmd(m.selectedSport.ID, m.selectedLeague.ID)
				}
			case gamesView:
				if m.gameCursor < len(m.games) {
					m.state = gameDetailView
					m.loadingDetail = true
					m.detailScrollOffset = 0
					selectedGame := m.games[m.gameCursor]
					return m, loadGameDetailCmd(m.selectedSport.ID, m.selectedLeague.ID, selectedGame.ID)
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

	case gameDetailLoadedMsg:
		m.loadingDetail = false
		m.selectedGameDetail = msg.detail
		m.err = msg.err
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
	case gameDetailView:
		content = m.renderGameDetailView()
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

	// Calculate visible games window
	visibleGames := m.calculateVisibleGames()
	startIdx := m.gameScrollOffset
	endIdx := startIdx + visibleGames
	if endIdx > len(m.games) {
		endIdx = len(m.games)
	}

	var items string
	for i := startIdx; i < endIdx; i++ {
		cursor := "  "
		if i == m.gameCursor {
			cursor = "❯ "
		}
		items += cursor + m.renderGame(m.games[i], i == m.gameCursor) + "\n\n"
	}

	// Add scroll indicator
	scrollIndicator := ""
	if len(m.games) > visibleGames {
		scrollIndicator = lipgloss.NewStyle().Foreground(dimColor).Render(fmt.Sprintf(" (Showing %d-%d of %d games)", startIdx+1, endIdx, len(m.games)))
	}

	help := helpStyle.Render("↑/k up • ↓/j down • enter details • r refresh • esc back • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		statusText+scrollIndicator,
		"",
		items,
		help,
	)
}

// calculateVisibleGames calculates how many game cards can fit in the viewport
func (m Model) calculateVisibleGames() int {
	// Each game card takes roughly 11 lines (including spacing):
	// - Status (1 line)
	// - Empty line (1)
	// - Away team (1)
	// - Home team (1)
	// - Empty line (1)
	// - Venue (1)
	// - Time (1)
	// - Box padding/borders (2)
	// - Spacing between cards (2)
	const linesPerGame = 11
	
	// Reserve space for: title (3), status (1), empty (1), help (2), margins (4)
	const reservedLines = 11
	
	availableHeight := m.height - reservedLines
	if availableHeight < linesPerGame {
		return 1 // Always show at least 1 game
	}
	
	return availableHeight / linesPerGame
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

func (m Model) renderGameDetailView() string {
	if m.loadingDetail {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("🏆 Game Details"),
			subtitleStyle.Render("Loading game details..."),
		)
	}

	if m.err != nil {
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		help := helpStyle.Render("esc back • q quit")
		return lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("🏆 Game Details"), "", errorMsg, "", help)
	}

	if m.selectedGameDetail == nil {
		return "No game details available"
	}

	detail := m.selectedGameDetail

	// Build content lines
	var contentLines []string

	// Header with scores
	status := detail.Status
	if detail.IsLive {
		status = liveStyle.Render(fmt.Sprintf("🔴 LIVE - %s", detail.Status))
		if detail.Period != "" && detail.Clock != "" {
			status += statusStyle.Render(fmt.Sprintf(" • %s %s", detail.Period, detail.Clock))
		}
	} else {
		status = statusStyle.Render(status)
	}
	contentLines = append(contentLines, status)
	contentLines = append(contentLines, "")

	// Score summary
	awayRecord := ""
	homeRecord := ""
	if detail.AwayTeam.Record != "" {
		awayRecord = fmt.Sprintf(" (%s)", detail.AwayTeam.Record)
	}
	if detail.HomeTeam.Record != "" {
		homeRecord = fmt.Sprintf(" (%s)", detail.HomeTeam.Record)
	}

	contentLines = append(contentLines, 
		teamStyle.Render(fmt.Sprintf("%-35s %5s", detail.AwayTeam.Name+awayRecord, detail.AwayTeam.Score)))
	contentLines = append(contentLines, 
		teamStyle.Render(fmt.Sprintf("%-35s %5s", detail.HomeTeam.Name+homeRecord, detail.HomeTeam.Score)))
	contentLines = append(contentLines, "")

	// Venue and attendance
	if detail.Venue != "" {
		contentLines = append(contentLines, venueStyle.Render(fmt.Sprintf("📍 %s", detail.Venue)))
	}
	if detail.Attendance != "" {
		contentLines = append(contentLines, venueStyle.Render(fmt.Sprintf("👥 Attendance: %s", detail.Attendance)))
	}
	contentLines = append(contentLines, "")

	// Team Leaders
	if len(detail.Leaders) > 0 {
		contentLines = append(contentLines, lipgloss.NewStyle().Bold(true).Foreground(accentColor).Render("⭐ Game Leaders"))
		contentLines = append(contentLines, "")
		for _, leader := range detail.Leaders {
			leaderLine := fmt.Sprintf("  %s: %s (%s) - %s", 
				leader.Category, 
				leader.Athlete, 
				leader.Team, 
				leader.Value)
			contentLines = append(contentLines, itemStyle.Render(leaderLine))
		}
		contentLines = append(contentLines, "")
	}

	// Team Statistics
	if len(detail.HomeTeam.Statistics) > 0 || len(detail.AwayTeam.Statistics) > 0 {
		contentLines = append(contentLines, lipgloss.NewStyle().Bold(true).Foreground(accentColor).Render("📊 Team Statistics"))
		contentLines = append(contentLines, "")
		
		// Create columns for away and home stats
		maxStats := len(detail.AwayTeam.Statistics)
		if len(detail.HomeTeam.Statistics) > maxStats {
			maxStats = len(detail.HomeTeam.Statistics)
		}
		
		for i := 0; i < maxStats && i < 10; i++ { // Limit to 10 stats
			var statLine string
			if i < len(detail.AwayTeam.Statistics) {
				stat := detail.AwayTeam.Statistics[i]
				statLine = fmt.Sprintf("  %-20s: %8s", stat.Label, stat.Value)
			}
			if i < len(detail.HomeTeam.Statistics) {
				stat := detail.HomeTeam.Statistics[i]
				statLine += fmt.Sprintf("    |    %-20s: %8s", stat.Label, stat.Value)
			}
			contentLines = append(contentLines, statusStyle.Render(statLine))
		}
		contentLines = append(contentLines, "")
	}

	// Recent Plays
	if len(detail.Plays) > 0 {
		contentLines = append(contentLines, lipgloss.NewStyle().Bold(true).Foreground(accentColor).Render("📝 Recent Plays"))
		contentLines = append(contentLines, "")
		
		playCount := 0
		for _, play := range detail.Plays {
			if playCount >= 15 { // Limit display
				break
			}
			
			playPrefix := "  "
			playTextStyle := statusStyle
			if play.ScoringPlay {
				playPrefix = "🎯 "
				playTextStyle = liveStyle
			}
			
			clockInfo := ""
			if play.Period != "" && play.Clock != "" {
				clockInfo = fmt.Sprintf("[%s %s] ", play.Period, play.Clock)
			}
			
			playText := fmt.Sprintf("%s%s%s", playPrefix, clockInfo, play.Text)
			contentLines = append(contentLines, playTextStyle.Render(playText))
			playCount++
		}
	}

	// Calculate visible content
	availableHeight := m.height - 8 // Reserve for title, help, margins
	startLine := m.detailScrollOffset
	endLine := startLine + availableHeight
	if endLine > len(contentLines) {
		endLine = len(contentLines)
	}
	if startLine >= len(contentLines) {
		startLine = 0
	}

	visibleContent := lipgloss.JoinVertical(lipgloss.Left, contentLines[startLine:endLine]...)

	// Scroll indicator
	scrollInfo := ""
	if len(contentLines) > availableHeight {
		scrollInfo = lipgloss.NewStyle().Foreground(dimColor).Render(
			fmt.Sprintf(" (Scroll: %d/%d lines)", startLine+1, len(contentLines)))
	}

	title := titleStyle.Render("🏆 Game Details")
	help := helpStyle.Render("↑/k up • ↓/j down • esc back • q quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title+scrollInfo,
		"",
		visibleContent,
		"",
		help,
	)
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

