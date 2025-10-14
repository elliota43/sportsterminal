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

type GameDetail struct {
	ID          string
	Name        string
	Status      string
	StatusDetail string
	IsLive      bool
	HomeTeam    TeamDetail
	AwayTeam    TeamDetail
	Venue       string
	Attendance  string
	Plays       []Play
	Leaders     []Leader
	Period      string
	Clock       string
}

type TeamDetail struct {
	Name       string
	ShortName  string
	Score      string
	Record     string
	Logo       string
	Statistics []Statistic
}

type Statistic struct {
	Label string
	Value string
}

type Play struct {
	Period      string
	Clock       string
	Text        string
	ScoringPlay bool
	Team        string
}

type Leader struct {
	Category string
	Team     string
	Athlete  string
	Value    string
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

func GetGameDetail(sport string, league string, eventID string) (*GameDetail, error) {
	url := fmt.Sprintf("%s/%s/%s/summary?event=%s", espnAPIBase, sport, league, eventID)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	detail := &GameDetail{
		ID: eventID,
	}

	// Extract header info
	if header, ok := result["header"].(map[string]interface{}); ok {
		if competitions, ok := header["competitions"].([]interface{}); ok && len(competitions) > 0 {
			comp := competitions[0].(map[string]interface{})
			
			// Status
			if status, ok := comp["status"].(map[string]interface{}); ok {
				if statusType, ok := status["type"].(map[string]interface{}); ok {
					detail.Status = getString(statusType, "description")
					detail.StatusDetail = getString(statusType, "detail")
					detail.IsLive = getString(statusType, "state") == "in"
				}
				detail.Period = getString(status, "period")
				detail.Clock = getString(status, "displayClock")
			}

			// Venue
			if venue, ok := comp["venue"].(map[string]interface{}); ok {
				detail.Venue = getString(venue, "fullName")
			}

			// Attendance
			if attendance, ok := comp["attendance"].(float64); ok {
				detail.Attendance = fmt.Sprintf("%.0f", attendance)
			}

			// Teams
			if competitors, ok := comp["competitors"].([]interface{}); ok {
				for _, c := range competitors {
					competitor := c.(map[string]interface{})
					teamDetail := parseTeamDetail(competitor)
					
					if getString(competitor, "homeAway") == "home" {
						detail.HomeTeam = teamDetail
					} else {
						detail.AwayTeam = teamDetail
					}
				}
			}
		}
	}

	// Extract box score
	if boxscore, ok := result["boxscore"].(map[string]interface{}); ok {
		if teams, ok := boxscore["teams"].([]interface{}); ok {
			for _, t := range teams {
				team := t.(map[string]interface{})
				teamID := getString(team, "team", "id")
				
				stats := []Statistic{}
				if statistics, ok := team["statistics"].([]interface{}); ok {
					for _, s := range statistics {
						stat := s.(map[string]interface{})
						stats = append(stats, Statistic{
							Label: getString(stat, "label"),
							Value: getString(stat, "displayValue"),
						})
					}
				}

				// Match to home/away team
				if detail.HomeTeam.Name != "" && teamID != "" {
					detail.HomeTeam.Statistics = stats
				}
				if detail.AwayTeam.Name != "" && teamID != "" {
					detail.AwayTeam.Statistics = stats
				}
			}
		}
	}

	// Extract plays (last 20 significant plays)
	if plays, ok := result["plays"].([]interface{}); ok {
		playCount := 0
		for i := len(plays) - 1; i >= 0 && playCount < 20; i-- {
			play := plays[i].(map[string]interface{})
			
			// Only include significant plays
			if scoringPlay, _ := play["scoringPlay"].(bool); scoringPlay || 
				getString(play, "type", "text") != "" {
				
				detail.Plays = append([]Play{{
					Period:      getString(play, "period", "displayValue"),
					Clock:       getString(play, "clock", "displayValue"),
					Text:        getString(play, "text"),
					ScoringPlay: scoringPlay,
					Team:        getString(play, "team", "shortDisplayName"),
				}}, detail.Plays...)
				playCount++
			}
		}
	}

	// Extract leaders
	if leaders, ok := result["leaders"].([]interface{}); ok {
		for _, l := range leaders {
			leader := l.(map[string]interface{})
			if leaders, ok := leader["leaders"].([]interface{}); ok && len(leaders) > 0 {
				topLeader := leaders[0].(map[string]interface{})
				detail.Leaders = append(detail.Leaders, Leader{
					Category: getString(leader, "displayName"),
					Team:     getString(topLeader, "team", "shortDisplayName"),
					Athlete:  getString(topLeader, "athlete", "displayName"),
					Value:    getString(topLeader, "displayValue"),
				})
			}
		}
	}

	if result["gameInfo"] != nil {
		detail.Name = getString(result, "gameInfo", "venue", "fullName")
	}

	return detail, nil
}

func parseTeamDetail(competitor map[string]interface{}) TeamDetail {
	td := TeamDetail{}
	
	if team, ok := competitor["team"].(map[string]interface{}); ok {
		td.Name = getString(team, "displayName")
		td.ShortName = getString(team, "shortDisplayName")
		td.Logo = getString(team, "logo")
	}
	
	td.Score = getString(competitor, "score")
	
	if records, ok := competitor["records"].([]interface{}); ok && len(records) > 0 {
		record := records[0].(map[string]interface{})
		td.Record = getString(record, "summary")
	}
	
	return td
}

func getString(m map[string]interface{}, keys ...string) string {
	current := m
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

