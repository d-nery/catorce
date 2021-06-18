package game

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/d-nery/catorce/pkg/deck"
	"github.com/rs/zerolog"
)

type Game struct {
	Chat          int64
	Players       *ring.Ring
	Deck          *deck.Deck
	Reversed      bool
	State         GameState
	DrawCount     int
	CurrentCard   *deck.Card
	PlayerCatorce *Player

	TurnStarted time.Time

	// Current game stats, are added to overall when game is over
	Rounds              int
	P2Sequence          int
	P4Played            int
	LargestResponseTime time.Duration

	logger zerolog.Logger
}

func New(chat int64, logger zerolog.Logger) *Game {
	logger.Trace().Int64("chat", chat).Msg("Creating new game")
	return &Game{
		Chat:          chat,
		Players:       nil,
		Reversed:      false,
		Deck:          deck.New(),
		State:         LOBBY,
		DrawCount:     0,
		CurrentCard:   nil,
		PlayerCatorce: nil,
		Rounds:        0,
		P2Sequence:    0,
		P4Played:      0,
		TurnStarted:   time.Time{},

		logger: logger.With().Int64("game_chat_id", chat).Logger(),
	}
}

func (g *Game) SetLogger(logger zerolog.Logger) {
	g.logger = logger.With().Int64("game_chat_id", g.Chat).Logger()
}

func (g *Game) CurrentPlayer() *Player {
	return g.Players.Value.(*Player)
}

func (g *Game) PlayerList() []*Player {
	if g.Players == nil {
		return nil
	}

	players := make([]*Player, 0, g.Players.Len())

	g.Players.Do(func(i interface{}) {
		players = append(players, i.(*Player))
	})

	return players
}

func (g *Game) GetPlayer(id int) *Player {
	var player *Player

	g.Players.Do(func(i interface{}) {
		if i.(*Player).ID == id {
			player = i.(*Player)
		}
	})

	return player
}

func (g *Game) PlayerAmount() int {
	return g.Players.Len()
}

func (g *Game) AddPlayer(p *Player) {
	g.logger.Trace().Msg("Adding player")
	if g.State != LOBBY {
		g.logger.Trace().Msg("Failed: can't add player outside of lobby")
		return
	}

	if g.Players == nil {
		g.logger.Trace().Int("id", p.ID).Msg("No players, adding first")
		g.Players = ring.New(1)
		g.Players.Value = p
		return
	}

	r := ring.New(1)
	r.Value = p

	g.Players = g.Players.Prev()
	g.Players = g.Players.Link(r)

	if g.logger.GetLevel() <= zerolog.TraceLevel {
		var out strings.Builder

		fmt.Fprint(&out, "Current order: ")
		g.Players.Do(func(i interface{}) {
			fmt.Fprintf(&out, "%d, ", i.(*Player).ID)
		})
		g.logger.Trace().Msg(out.String())
	}
}

func (g *Game) ShufflePlayers() {
	g.logger.Trace().Msg("Shuffling players")
	players := g.PlayerList()

	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	g.Players.Value = players[0]
	i := 1
	for p := g.Players.Next(); p != g.Players; p = p.Next() {
		p.Value = players[i]
		i += 1
	}

	if g.logger.GetLevel() <= zerolog.TraceLevel {
		var out strings.Builder

		fmt.Fprint(&out, "Current order: ")
		g.Players.Do(func(i interface{}) {
			fmt.Fprintf(&out, "%d, ", i.(*Player).ID)
		})
		g.logger.Trace().Msg(out.String())
	}
}

func (g *Game) PlayFirstCard() {
	g.logger.Trace().Msg("Playing first card")
	if g.CurrentCard != nil {
		return
	}

	g.CurrentCard = g.Deck.Draw()

	// Can't start with a special card
	for g.CurrentCard.IsSpecial() {
		g.logger.Trace().Str("card", g.CurrentCard.String()).Msg("Got special card, redrawing")
		g.Deck.Discard(g.CurrentCard)
		g.CurrentCard = g.Deck.Draw()
	}

	g.logger.Trace().Str("card", g.CurrentCard.String()).Msg("First card played")

	switch g.CurrentCard.GetValue() {
	case deck.SKIP:
		g.EndTurn(false)
	case deck.DRAW:
		g.DrawCount += 2
	case deck.REVERSE:
		if g.PlayerAmount() != 2 {
			g.Reverse()
		} else {
			g.EndTurn(false)
		}
	}

	g.TurnStarted = time.Now()
}

func (g *Game) DistributeCards() {
	g.logger.Trace().Msg("Distributing cards")
	if g.State != LOBBY {
		g.logger.Trace().Msg("Failed, can't distribute cards outside lobby")
		return
	}

	for i := 0; i < 7; i++ {
		g.Players.Do(func(i interface{}) {
			p := i.(*Player)
			card := g.Deck.Draw()
			p.AddCard(card)
		})
	}
}

func (g *Game) PlayCard(c *deck.Card) {
	g.logger.Trace().Msg("Playing card")

	if g.HasPendingCatorce() {
		g.logger.Trace().Str("player_name", g.PlayerCatorce.Name).Msg("There's a pending catorce!")
		g.PlayerCatorce.CatorcesMissed += 1

		for i := 0; i < 4; i++ {
			card := g.Deck.Draw()
			g.PlayerCatorce.AddCard(card)
		}
		g.PlayerCatorce = nil
	}

	g.CurrentPlayer().CardsPlayed += 1
	turnDuration := time.Since(g.TurnStarted)
	g.CurrentPlayer().AddDuration(turnDuration)

	if turnDuration > g.LargestResponseTime {
		g.LargestResponseTime = turnDuration
	}

	g.Deck.Discard(g.CurrentCard)
	g.CurrentCard = c

	if c.IsSpecial() {
		switch c.GetSpecial() {
		case deck.DFOUR:
			g.P4Played += 1
			g.DrawCount += 4
		}

		g.logger.Debug().Str("from", string(g.State)).Str("to", "CHOOSE_COLOR").Msg("Changing state")
		g.State = CHOOSE_COLOR
		return
	}

	jump := false

	switch c.GetValue() {
	case deck.SKIP:
		jump = true
	case deck.DRAW:
		g.DrawCount += 2
	case deck.REVERSE:
		if g.PlayerAmount() != 2 {
			g.Reverse()
		} else {
			jump = true
		}
	}

	g.EndTurn(jump)
}

func (g *Game) DrawCard() {
	g.logger.Trace().Msg("Drawing a card")

	if g.HasPendingCatorce() {
		g.logger.Trace().Str("player_name", g.PlayerCatorce.Name).Msg("There's a pending catorce!")
		for i := 0; i < 4; i++ {
			card := g.Deck.Draw()
			g.PlayerCatorce.AddCard(card)
		}
		g.PlayerCatorce = nil
	}

	if g.DrawCount == 0 {
		card := g.Deck.Draw()
		g.CurrentPlayer().AddCard(card)
		g.logger.Debug().Str("from", string(g.State)).Str("to", "DREW").Msg("Changing state")
		g.State = DREW
		return
	}

	if g.DrawCount > g.P2Sequence {
		g.P2Sequence = g.DrawCount
	}

	for i := 0; i < g.DrawCount; i++ {
		card := g.Deck.Draw()
		g.CurrentPlayer().AddCard(card)
	}

	g.DrawCount = 0
	g.EndTurn(false)
}

// EndTurn finishes the turn, returns true if the game is over
func (g *Game) EndTurn(jump bool) bool {
	g.Rounds += 1
	g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Int("rounds", g.Rounds).Msg("Ending turn")
	if len(g.CurrentPlayer().Hand) == 0 {
		g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Msg("Player has 0 cards")
		g.State = LOBBY
		return true
	}

	if len(g.CurrentPlayer().Hand) == 1 {
		g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Msg("Player has 1 card left, setting catorce")
		g.PlayerCatorce = g.CurrentPlayer()
	}

	g.Players = g.NextPlayer()
	if jump {
		g.Players = g.NextPlayer()
	}

	g.logger.Debug().Str("from", string(g.State)).Str("to", "CHOOSE_CARD").Msg("Changing state")
	g.State = CHOOSE_CARD

	g.TurnStarted = time.Now()
	return false
}

func (g *Game) Reverse() {
	g.logger.Trace().Msg("Reversing game")
	g.Reversed = !g.Reversed
}

func (g *Game) ChooseColor(c deck.Color) {
	// We change the card color to the chosen color, this only
	// affects special cards, so we don't see it as they are always black
	g.logger.Trace().Str("color", string(c)).Msg("Setting card color")
	g.CurrentCard.SetColor(c)
	g.EndTurn(false)
}

func (g *Game) NextPlayer() *ring.Ring {
	g.logger.Trace().Msg("Moving to next player")
	if g.Reversed {
		return g.Players.Prev()
	}

	return g.Players.Next()
}

func (g *Game) GetCurrentCard() *deck.Card {
	return g.CurrentCard
}

func (g *Game) GetDeck() *deck.Deck {
	return g.Deck
}

func (g *Game) GetState() GameState {
	return g.State
}

func (g *Game) DrawCounter() int {
	return g.DrawCount
}

func (g *Game) ResetDrawCounter() {
	g.logger.Trace().Msg("Resetting draw counter")
	g.DrawCount = 0
}

func (g *Game) CurrentCardSticker() *tb.Sticker {
	return &tb.Sticker{
		File: tb.File{FileID: g.GetCurrentCard().Sticker()},
	}
}

func (g *Game) HasPendingCatorce() bool {
	return g.PlayerCatorce != nil
}

func (g *Game) CatorcePlayer() *Player {
	return g.PlayerCatorce
}

func (g *Game) GameInfo() string {
	var out strings.Builder
	fmt.Fprintf(&out, "Jogador atual: %s \\[%d]\n", g.CurrentPlayer().NameWithMention(), len(g.CurrentPlayer().Hand))
	fmt.Fprintf(&out, "Última carta: %s\n", g.GetCurrentCard().StringPretty())
	fmt.Fprint(&out, "Próximos Jogadores:\n")

	for _, p := range g.PlayerList() {
		if p == g.CurrentPlayer() {
			continue
		}

		fmt.Fprintf(&out, " • %s \\[%d]\n", p.Name, len(p.Hand))
	}

	fmt.Fprintf(&out, "Cartas na pilha: %d\n", g.Deck.Available())
	if g.Reversed {
		fmt.Fprintf(&out, "*Invertido*")
	}

	return out.String()
}

func (g *Game) UnmarshalJSON(data []byte) error {
	type Alias Game

	v := &struct {
		Players []*Player
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for _, p := range v.Players {
		if g.Players == nil {
			g.Players = ring.New(1)
			g.Players.Value = p
			continue
		}

		r := ring.New(1)
		r.Value = p

		g.Players = g.Players.Prev()
		g.Players = g.Players.Link(r)
	}

	g.PlayerCatorce = nil
	return nil
}

func (g *Game) MarshalJSON() ([]byte, error) {
	type Alias Game
	return json.Marshal(&struct {
		Players []*Player
		*Alias
	}{
		Players: g.PlayerList(),
		Alias:   (*Alias)(g),
	})
}
