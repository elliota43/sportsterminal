package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	espnAPIBase = "https://site.api.espn.com/apis/site/v2/sports"
)

type Sport struct {
	Name string
	ID   string
	Leagues []League
}

type League struct {
	Name string
	ID   string
}

type Game struct {
	ID        string
	Name      string
	ShortName string
	Date      time.Time
	Status    string
	HomeTeam  Team
	AwayTeam  Team
	IsLive    bool
	Venue     string
}

type Team struct {
	Name      string
	ShortName string
	Score     string
	Logo      string
}

type ESPNResponse struct {
	Events []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ShortName   string `json:"shortName"`
		Date        string `json:"date"`
		Competitions []struct {
			ID        string `json:"id"`
			Venue     struct {
				FullName string `json:"fullName"`
			} `json:"venue"`
			Status struct {
				Type struct {
					State       string `json:"state"`
					Completed   bool   `json:"completed"`
					Description string `json:"description"`
				} `json:"type"`
			} `json:"status"`
			Competitors []struct {
				ID       string `json:"id"`
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Team     struct {
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Logo             string `json:"logo"`
				} `json:"team"`
				Score string `json:"score"`
			} `json:"competitors"`
		} `json:"competitions"`
	} `json:"events"`
}

var AvailableSports = []Sport{
	{
		Name: "Football",
		ID:   "football",
		Leagues: []League{
			{Name: "NFL", ID: "nfl"},
			{Name: "College Football", ID: "college-football"},
		},
	},
	{
		Name: "Basketball",
		ID:   "basketball",
		Leagues: []League{
			{Name: "NBA", ID: "nba"},
			{Name: "WNBA", ID: "wnba"},
			{Name: "College Basketball (Men)", ID: "mens-college-basketball"},
			{Name: "College Basketball (Women)", ID: "womens-college-basketball"},
		},
	},
	{
		Name: "Baseball",
		ID:   "baseball",
		Leagues: []League{
			{Name: "MLB", ID: "mlb"},
			{Name: "College Baseball", ID: "college-baseball"},
		},
	},
	{
		Name: "Hockey",
		ID:   "hockey",
		Leagues: []League{
			{Name: "NHL", ID: "nhl"},
		},
	},
	{
		Name: "Soccer",
		ID:   "soccer",
		Leagues: []League{
			{Name: "Premier League", ID: "eng.1"},
			{Name: "La Liga", ID: "esp.1"},
			{Name: "Serie A", ID: "ita.1"},
			{Name: "Bundesliga", ID: "ger.1"},
			{Name: "MLS", ID: "usa.1"},
			{Name: "Champions League", ID: "uefa.champions"},
		},
	},
}

func GetGames(sport string, league string) ([]Game, error) {
	url := fmt.Sprintf("%s/%s/%s/scoreboard", espnAPIBase, sport, league)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch games: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var espnResp ESPNResponse
	if err := json.Unmarshal(body, &espnResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	games := make([]Game, 0, len(espnResp.Events))
	for _, event := range espnResp.Events {
		if len(event.Competitions) == 0 {
			continue
		}

		comp := event.Competitions[0]
		game := Game{
			ID:        event.ID,
			Name:      event.Name,
			ShortName: event.ShortName,
			Status:    comp.Status.Type.Description,
			IsLive:    comp.Status.Type.State == "in",
			Venue:     comp.Venue.FullName,
		}

		// Parse date
		if t, err := time.Parse(time.RFC3339, event.Date); err == nil {
			game.Date = t
		}

		// Extract team information
		for _, competitor := range comp.Competitors {
			team := Team{
				Name:      competitor.Team.DisplayName,
				ShortName: competitor.Team.ShortDisplayName,
				Score:     competitor.Score,
				Logo:      competitor.Team.Logo,
			}

			if competitor.HomeAway == "home" {
				game.HomeTeam = team
			} else {
				game.AwayTeam = team
			}
		}

		games = append(games, game)
	}

	return games, nil
}

