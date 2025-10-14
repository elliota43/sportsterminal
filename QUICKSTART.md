# 🚀 Quick Start Guide

## Running the Application

```bash
# Build and run in one command
make run

# Or build first, then run
make build
./sportsterminal
```

## Navigation Flow

1. **Sport Selection** (First Screen)
   - Use `↑/↓` or `k/j` to navigate
   - Press `Enter` to select a sport
   - Available sports: Football 🏈, Basketball 🏀, Baseball ⚾, Hockey 🏒, Soccer ⚽

2. **League Selection** (Second Screen)
   - Choose from available leagues within the selected sport
   - Examples: NFL, NBA, Premier League, Champions League
   - Press `Esc` to go back to sport selection

3. **Games View** (Third Screen)
   - See all games with live scores
   - Live games are marked with 🔴 LIVE
   - Game information includes:
     - Team names and scores
     - Game status (Final, In Progress, Scheduled)
     - Venue location
     - Date and time
   - Press `Enter` to view detailed game information
   - Press `r` to manually refresh scores
   - Auto-refreshes every 30 seconds for live games
   - Press `Esc` to go back to league selection

4. **Game Detail View** (Fourth Screen)
   - View comprehensive game information:
     - Live score with period and clock (for live games)
     - Team records
     - Venue and attendance
     - Game leaders (top performers)
     - Team statistics (box score)
     - Recent plays and scoring plays (marked with 🎯)
   - Scroll through content with `↑/↓` or `k/j`
   - Press `Esc` to return to games list

## Keyboard Controls Quick Reference

```
Navigation:
  ↑/k        Move up
  ↓/j        Move down
  Enter/→/l  Select / Go forward
  Esc/←      Go back
  
Actions:
  r          Refresh scores (on games view)
  q          Quit application
  Ctrl+C     Force quit
```

## Example Session

```
1. Start the app: ./sportsterminal
2. Navigate to "Soccer" and press Enter
3. Select "Premier League" and press Enter
4. View live scores and upcoming matches
5. Press 'r' to refresh scores manually
6. Press Esc to try another league/sport
7. Press 'q' to quit
```

## Tips

- **Live Games**: The app automatically refreshes scores every 30 seconds when there are live games
- **Keyboard Shortcuts**: Works with both arrow keys and Vim-style hjkl keys
- **Navigation**: You can quickly navigate back using Esc or Backspace
- **Multiple Sports**: Easily switch between different sports and leagues

## Troubleshooting

**No games showing?**
- Some leagues may not have games scheduled today
- Try another league or sport
- Press 'r' to manually refresh

**API Errors?**
- Check your internet connection
- ESPN API might be temporarily unavailable
- Try refreshing with 'r'

**Terminal too small?**
- Increase your terminal window size for better viewing
- Minimum recommended: 80x24 characters

Enjoy checking sports scores! 🏆

