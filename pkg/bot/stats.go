package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/d-nery/catorce/pkg/game"
	"github.com/jedib0t/go-pretty/v6/table"
)

// OverallStats maps chat IDs to Stats
type OverallStats map[int64]*ChatStats

// ChatStats holds group and player's stats for the current group
type ChatStats struct {
	Group   GroupStats
	Players map[int]*PlayerStats
}

// GroupStats holds group stats for a specific chat
type GroupStats struct {
	GamesPlayed         int
	P2Sequence          int
	P4Played            int
	RoundsPlayed        int
	LargestResponseTime time.Duration
}

// PlayerStats holds player stats for a specific chat
type PlayerStats struct {
	Name            string
	GamesWon        int
	GamesPlayed     int
	Points          int
	CatorcesCalled  int
	CatorcesMissed  int
	CardsPlayed     int
	AvgResponseTime time.Duration
}

// ReadStatsFromFile loads OverallStats from a .json file
func ReadStatsFromFile(file string) (OverallStats, error) {
	s := OverallStats{}
	body, err := os.ReadFile(file)

	if err != nil {
		return s, err
	}

	err = json.Unmarshal(body, &s)
	return s, err
}

// Persist persists OverallStats into a .json file
func (s *OverallStats) Persist(file string) error {
	body, err := json.Marshal(s)

	if err != nil {
		return err
	}

	return os.WriteFile(file, body, 0644)
}

// Persist dumps OverallStats to the terminal
func (s *OverallStats) Dump() error {
	body, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

// AddGameStats adds stats from the game to the GroupStats
func (gs *GroupStats) AddGameStats(g *game.Game) {
	gs.GamesPlayed += 1

	if g.LargestResponseTime > gs.LargestResponseTime {
		gs.LargestResponseTime = g.LargestResponseTime
	}

	if g.P2Sequence > gs.P2Sequence {
		gs.P2Sequence = g.P2Sequence
	}

	gs.P4Played += g.P4Played
	gs.RoundsPlayed += g.Rounds
}

// AddPlayerStats adds stats from the player to the PlayerStats
func (ps *PlayerStats) AddPlayerStats(p *game.Player) {
	ps.GamesPlayed += 1
	if len(p.Hand) == 0 {
		ps.GamesWon += 1
	}

	// Cumulative Rolling Average
	if ps.CardsPlayed+p.CardsPlayed > 0 {
		ps.AvgResponseTime = time.Duration(
			(int64(ps.CardsPlayed)*ps.AvgResponseTime.Nanoseconds() + int64(p.CardsPlayed)*p.AvgRespTime.Nanoseconds()) /
				int64(ps.CardsPlayed+p.CardsPlayed))
	}

	ps.CardsPlayed += p.CardsPlayed
	ps.Points += p.CurrentHandPoints()
	ps.CatorcesCalled += p.CatorcesCalled
	ps.CatorcesMissed += p.CatorcesMissed
}

// SaveGameStats saves game and player's stats to the Bot's overall stats
// Should only be called after the game is finished
func (b *Bot) SaveGameStats(g *game.Game) {
	stats := b.stats[g.Chat]
	stats.Group.AddGameStats(g)

	for _, p := range g.PlayerList() {
		if _, ok := stats.Players[p.ID]; !ok {
			stats.Players[p.ID] = &PlayerStats{Name: p.Name}
		}

		stats.Players[p.ID].AddPlayerStats(p)
	}
}

// Report generates a Markdown formatted string with GroupStats report
func (gs *GroupStats) Report() string {
	var out strings.Builder

	fmt.Fprintf(&out, "*Estatísticas para esse grupo*\n\n")
	fmt.Fprintf(&out, "Total de Jogos: %d\n", gs.GamesPlayed)
	fmt.Fprintf(&out, "Total de Rounds: %d\n\n", gs.RoundsPlayed)
	fmt.Fprintf(&out, "Maior sequência de +2: +%d\n", gs.P2Sequence)
	fmt.Fprintf(&out, "Quantidade de +4 jogados: %d\n\n", gs.P4Played)
	fmt.Fprintf(&out, "Maior tempo de resposta: %s", gs.LargestResponseTime.Round(time.Minute))

	return out.String()
}

// Ranking generates a Markdown formatted table chat ranking
func (cs *ChatStats) Ranking() string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Nome", "Pontos", "Jogos", "Média"})

	for _, ps := range cs.Players {
		t.AppendRow(table.Row{ps.Name, ps.Points, ps.GamesPlayed, ps.Points / ps.GamesPlayed})
	}

	t.SortBy([]table.SortBy{
		{Name: "Média", Mode: table.AscNumeric},
	})

	// t.SetStyle(table.StyleLight)
	return "```\n" + t.Render() + "\n```"
}

// Report generates a Markdown formatted string with PlayerStats report
func (ps *PlayerStats) Report() string {
	var out strings.Builder

	fmt.Fprintf(&out, "*Suas Estatísticas*\n\n")
	fmt.Fprintf(&out, "Total de jogos: %d\n", ps.GamesPlayed)
	fmt.Fprintf(&out, "Total de jogos vencidos: %d\n", ps.GamesWon)
	fmt.Fprintf(&out, "Total de pontos (menos é melhor): %d\n", ps.Points)
	fmt.Fprintf(&out, "Total de cartas jogadas: %d\n", ps.CardsPlayed)
	fmt.Fprintf(&out, "Catorces: %d/%d\n\n", ps.CatorcesCalled, ps.CatorcesCalled+ps.CatorcesMissed)
	fmt.Fprintf(&out, "Tempo médio de resposta: %s", ps.AvgResponseTime.Round(time.Second))

	return out.String()
}
