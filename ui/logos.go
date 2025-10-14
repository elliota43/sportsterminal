package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Terminal capability detection
func detectTerminalCapabilities() (supportsImages bool) {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	// Check for Kitty terminal
	if strings.Contains(term, "kitty") || termProgram == "kitty" {
		return true
	}

	// Check for iTerm2
	if termProgram == "iTerm.app" {
		return true
	}

	// Check for WezTerm
	if strings.Contains(term, "wezterm") {
		return true
	}

	return false
}

// Display team logo if terminal supports it
func displayTeamLogo(logoURL, teamName string) string {
	supportsImages := detectTerminalCapabilities()

	if !supportsImages {
		// Fallback to team abbreviation or name
		return getTeamEmoji(teamName)
	}

	// For now, return emoji fallback
	// In the future, we could implement actual image display:
	// - Kitty: use kitty icat command
	// - iTerm2: use inline image escape codes
	// - WezTerm: use image display functions
	return getTeamEmoji(teamName)
}

// Get appropriate emoji for team names
func getTeamEmoji(teamName string) string {
	name := strings.ToLower(teamName)

	// NBA teams (using team colors/identities)
	nbaEmojis := map[string]string{
		"lakers": "🟡", "warriors": "🏀", "celtics": "🟢", "bulls": "🔴",
		"heat": "🔥", "spurs": "⚫", "pistons": "🔵", "cavaliers": "🏹", // Cavalier = knight with bow
		"knicks": "🟠", "nets": "⚫", "76ers": "🔵", "raptors": "🔴",
		"hawks": "🔴", "hornets": "🟣", "magic": "🔵", "wizards": "🔴",
		"bucks": "🟢", "pacers": "🟡", "rockets": "🔴", "mavericks": "🔵",
		"grizzlies": "🔵", "pelicans": "🟣", "suns": "🟡", "jazz": "🟡",
		"nuggets": "🔵", "timberwolves": "🟢", "thunder": "🟡", "blazers": "🔴",
		"kings": "🟣", "clippers": "🔵",
	}

	// NFL teams
	nflEmojis := map[string]string{
		"patriots": "🔴", "bills": "🔴", "dolphins": "🔵", "jets": "🟢",
		"steelers": "🟡", "ravens": "🟣", "browns": "🟠", "bengals": "🟠",
		"texans": "🔴", "colts": "🔵", "jaguars": "🟢", "titans": "🔵",
		"chiefs": "🔴", "raiders": "⚫", "chargers": "🔵", "broncos": "🟠",
		"cowboys": "🔵", "eagles": "🟢", "giants": "🔵", "commanders": "🔴",
		"packers": "🟢", "vikings": "🟣", "bears": "🟠", "lions": "🔵",
		"falcons": "🔴", "panthers": "🔵", "saints": "🟣", "buccaneers": "🔴",
		"cardinals": "🔴", "49ers": "🔴", "seahawks": "🟢", "rams": "🟡",
	}

	// MLB teams
	mlbEmojis := map[string]string{
		"yankees": "🔵", "red sox": "🔴", "blue jays": "🔵", "orioles": "🟠",
		"rays": "🔵", "astros": "🟠", "angels": "🔴", "athletics": "🟢",
		"mariners": "🔵", "rangers": "🔴", "twins": "🔵", "white sox": "⚫",
		"guardians": "🔵", "tigers": "🟠", "royals": "🔵", "braves": "🔴",
		"mets": "🔵", "phillies": "🔴", "marlins": "🔵", "nationals": "🔴",
		"cubs": "🔵", "cardinals": "🔴", "brewers": "🟡", "pirates": "⚫",
		"reds": "🔴", "dodgers": "🔵", "giants": "🟠", "padres": "🟡",
		"diamondbacks": "🔴", "rockies": "🟣",
	}

	// Check all leagues
	for keyword, emoji := range nbaEmojis {
		if strings.Contains(name, keyword) {
			return emoji
		}
	}
	for keyword, emoji := range nflEmojis {
		if strings.Contains(name, keyword) {
			return emoji
		}
	}
	for keyword, emoji := range mlbEmojis {
		if strings.Contains(name, keyword) {
			return emoji
		}
	}

	// Default fallback
	return "🏆"
}

// Future implementation for actual image display
func displayImageWithKitty(imageURL string) error {
	// This would use: kitty icat --place 32x32@0x0 image.png
	cmd := exec.Command("kitty", "icat", "--place", "32x32@0x0", imageURL)
	return cmd.Run()
}

func displayImageWithITerm2(imageURL string) error {
	// This would use iTerm2's inline image escape codes
	// \033]1337;File=size=%d;inline=1:%s\007
	fmt.Printf("\033]1337;File=inline=1:%s\007", imageURL)
	return nil
}
