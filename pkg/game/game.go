package game

import (
	"container/ring"
	"fmt"
	"math/rand"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/d-nery/catorce/pkg/deck"
	"github.com/rs/zerolog"
)

type Game struct {
	chat           int64
	players        *ring.Ring
	deck           *deck.Deck
	reversed       bool
	state          GameState
	draw_count     int
	current_card   *deck.Card
	player_catorce *Player
	rounds         int

	logger zerolog.Logger
}

func New(chat int64, logger zerolog.Logger) *Game {
	logger.Trace().Int64("chat", chat).Msg("Creating new game")
	return &Game{
		chat:           chat,
		players:        nil,
		reversed:       false,
		deck:           deck.New(),
		state:          LOBBY,
		draw_count:     0,
		current_card:   nil,
		player_catorce: nil,
		rounds:         0,

		logger: logger.With().Int64("game_chat_id", chat).Logger(),
	}
}

func (g *Game) CurrentPlayer() *Player {
	return g.players.Value.(*Player)
}

func (g *Game) PlayerList() []*Player {
	if g.players == nil {
		return nil
	}

	players := make([]*Player, 0, g.players.Len())

	g.players.Do(func(i interface{}) {
		players = append(players, i.(*Player))
	})

	return players
}

func (g *Game) GetPlayer(id int) *Player {
	var player *Player

	g.players.Do(func(i interface{}) {
		if i.(*Player).ID == id {
			player = i.(*Player)
		}
	})

	return player
}

func (g *Game) PlayerAmount() int {
	return g.players.Len()
}

func (g *Game) AddPlayer(p *Player) {
	g.logger.Trace().Msg("Adding player")
	if g.state != LOBBY {
		g.logger.Trace().Msg("Failed: can't add player outside of lobby")
		return
	}

	if g.players == nil {
		g.logger.Trace().Int("id", p.ID).Msg("No players, adding first")
		g.players = ring.New(1)
		g.players.Value = p
		return
	}

	r := ring.New(1)
	r.Value = p

	g.players = g.players.Prev()
	g.players = g.players.Link(r)

	if g.logger.GetLevel() <= zerolog.TraceLevel {
		var out strings.Builder

		fmt.Fprint(&out, "Current order: ")
		g.players.Do(func(i interface{}) {
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

	g.players.Value = players[0]
	i := 1
	for p := g.players.Next(); p != g.players; p = p.Next() {
		p.Value = players[i]
		i += 1
	}

	if g.logger.GetLevel() <= zerolog.TraceLevel {
		var out strings.Builder

		fmt.Fprint(&out, "Current order: ")
		g.players.Do(func(i interface{}) {
			fmt.Fprintf(&out, "%d, ", i.(*Player).ID)
		})
		g.logger.Trace().Msg(out.String())
	}
}

func (g *Game) PlayFirstCard() {
	g.logger.Trace().Msg("Playing first card")
	if g.current_card != nil {
		return
	}

	g.current_card = g.deck.Draw()

	// Can't start with a special card
	for g.current_card.IsSpecial() {
		g.logger.Trace().Str("card", g.current_card.String()).Msg("Got special card, redrawing")
		g.deck.Discard(g.current_card)
		g.current_card = g.deck.Draw()
	}

	g.logger.Trace().Str("card", g.current_card.String()).Msg("First card played")

	switch g.current_card.Value() {
	case deck.SKIP:
		g.EndTurn()
	case deck.DRAW:
		g.draw_count += 2
	case deck.REVERSE:
		if g.PlayerAmount() != 2 {
			g.Reverse()
		} else {
			g.EndTurn()
		}
	}
}

func (g *Game) DistributeCards() {
	g.logger.Trace().Msg("Distributing cards")
	if g.state != LOBBY {
		g.logger.Trace().Msg("Failed, can't distribute cards outside lobby")
		return
	}

	for i := 0; i < 2; i++ {
		g.players.Do(func(i interface{}) {
			p := i.(*Player)
			card := g.deck.Draw()
			p.AddCard(card)
		})
	}
}

func (g *Game) PlayCard(c *deck.Card) {
	g.logger.Trace().Msg("Playing card")

	if g.HasPendingCatorce() {
		g.logger.Trace().Msg("There's a pending catorce!")
		for i := 0; i < 4; i++ {
			card := g.deck.Draw()
			g.player_catorce.AddCard(card)
		}
		g.player_catorce = nil
	}

	g.deck.Discard(g.current_card)
	g.current_card = c

	if c.IsSpecial() {
		switch c.Special() {
		case deck.DFOUR:
			g.draw_count += 4
		}

		g.logger.Debug().Str("from", string(g.state)).Str("to", "CHOOSE_COLOR").Msg("Changing state")
		g.state = CHOOSE_COLOR
		return
	}

	switch c.Value() {
	case deck.SKIP:
		g.EndTurn()
	case deck.DRAW:
		g.draw_count += 2
	case deck.REVERSE:
		if g.PlayerAmount() != 2 {
			g.Reverse()
		} else {
			g.EndTurn()
		}
	}

	g.EndTurn()
}

func (g *Game) DrawCard() {
	g.logger.Trace().Msg("Drawing a card")

	if g.draw_count == 0 {
		card := g.deck.Draw()
		g.CurrentPlayer().AddCard(card)
		g.logger.Debug().Str("from", string(g.state)).Str("to", "DREW").Msg("Changing state")
		g.state = DREW
		return
	}

	for i := 0; i < g.draw_count; i++ {
		card := g.deck.Draw()
		g.CurrentPlayer().AddCard(card)
	}

	g.draw_count = 0
	g.EndTurn()
}

// EndTurn finishes the turn, returns true if the game is over
func (g *Game) EndTurn() bool {
	g.rounds += 1
	g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Int("rounds", g.rounds).Msg("Ending turn")
	if len(g.CurrentPlayer().Hand) == 0 {
		g.state = LOBBY
		return true
	}

	if len(g.CurrentPlayer().Hand) == 1 {
		g.player_catorce = g.CurrentPlayer()
	}

	g.players = g.NextPlayer()

	g.logger.Debug().Str("from", string(g.state)).Str("to", "CHOOSE_CARD").Msg("Changing state")
	g.state = CHOOSE_CARD

	return false
}

func (g *Game) Reverse() {
	g.logger.Trace().Msg("Reversing game")
	g.reversed = !g.reversed
}

func (g *Game) ChooseColor(c *deck.Color) {
	// We change the card color to the chosen color, this only
	// affects special cards, so we don't see it as they are always black
	g.logger.Trace().Str("color", string(*c)).Msg("Setting card color")
	g.current_card.SetColor(c)
	g.EndTurn()
}

func (g *Game) NextPlayer() *ring.Ring {
	g.logger.Trace().Msg("Moving to next player")
	if g.reversed {
		return g.players.Prev()
	}

	return g.players.Next()
}

func (g *Game) CurrentCard() *deck.Card {
	return g.current_card
}

func (g *Game) GetDeck() *deck.Deck {
	return g.deck
}

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) DrawCounter() int {
	return g.draw_count
}

func (g *Game) ResetDrawCounter() {
	g.logger.Trace().Msg("Resetting draw counter")
	g.draw_count = 0
}

func (g *Game) CurrentCardSticker() *tb.Sticker {
	return &tb.Sticker{
		File: tb.File{FileID: g.CurrentCard().Sticker()},
	}
}

func (g *Game) HasPendingCatorce() bool {
	return g.player_catorce != nil
}

func (g *Game) CatorcePlayer() *Player {
	return g.player_catorce
}

func (g *Game) GameInfo() string {
	var out strings.Builder
	fmt.Fprintf(&out, "Jogador atual: %s\n", g.CurrentPlayer().NameWithMention())
	fmt.Fprintf(&out, "Última carta: %s\n", g.CurrentCard().StringPretty())
	fmt.Fprint(&out, "Próximos Jogadores:\n")

	for _, p := range g.PlayerList() {
		if p == g.CurrentPlayer() {
			continue
		}

		fmt.Fprintf(&out, " • %s \\[%d carta(s)]\n", p.Name, len(p.Hand))
	}

	fmt.Fprintf(&out, "Cartas na pilha: %d", g.deck.Available())

	return out.String()
}
