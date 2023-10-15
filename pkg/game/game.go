package game

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/d-nery/catorce/pkg/deck"
	"github.com/rs/zerolog"
)

type Game struct {
	Chat          int64
	Players       []*Player
	Deck          *deck.Deck
	State         GameState
	DrawCount     int
	CurrentCard   *deck.Card
	PlayerCatorce int
	Config        *Config

	TurnStarted time.Time

	// Current game stats, are added to overall when game is over
	Rounds              int
	P2Sequence          int
	P4Played            int
	LargestResponseTime time.Duration

	logger zerolog.Logger
	mx     sync.Mutex
}

func New(chat int64, logger zerolog.Logger, config *Config) *Game {
	logger.Trace().Int64("chat", chat).Msg("Creating new game")
	return &Game{
		Chat:          chat,
		Players:       []*Player{},
		Deck:          nil,
		State:         LOBBY,
		DrawCount:     0,
		CurrentCard:   nil,
		PlayerCatorce: 0,
		Rounds:        0,
		P2Sequence:    0,
		P4Played:      0,
		TurnStarted:   time.Time{},
		Config:        config,

		logger: logger.With().Int64("game_chat_id", chat).Logger(),
	}
}

func (g *Game) Lock() {
	g.mx.Lock()
}

func (g *Game) Unlock() {
	g.mx.Unlock()
}

func (g *Game) SetLogger(logger zerolog.Logger) {
	g.logger = logger.With().Int64("game_chat_id", g.Chat).Logger()
}

func (g *Game) CurrentPlayer() *Player {
	if len(g.Players) == 0 {
		return nil
	}

	return g.Players[0]
}

func (g *Game) PlayerList() []*Player {
	return g.Players
}

func (g *Game) GetPlayer(id int) *Player {
	for _, p := range g.Players {
		if p.ID == id {
			return p
		}
	}

	return nil
}

func (g *Game) PlayerAmount() int {
	return len(g.Players)
}

func (g *Game) AddPlayer(p *Player) {
	g.logger.Trace().Msg("Adding player")
	g.Players = append(g.Players, p)

	if g.logger.GetLevel() <= zerolog.TraceLevel {
		var out strings.Builder

		fmt.Fprint(&out, "Current order: ")
		for _, p := range g.Players {
			fmt.Fprintf(&out, "%d, ", p.ID)
		}
		g.logger.Trace().Msg(out.String())
	}
}

func (g *Game) ShufflePlayers() {
	g.logger.Trace().Msg("Shuffling players")

	rand.Shuffle(len(g.Players), func(i, j int) {
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	})

	if g.logger.GetLevel() <= zerolog.TraceLevel {
		var out strings.Builder

		fmt.Fprint(&out, "Current order: ")
		for _, p := range g.Players {
			fmt.Fprintf(&out, "%d, ", p.ID)
		}
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

	if g.CurrentCard.Type.Has(deck.DRAW) {
		g.DrawCount += g.CurrentCard.Value
	}

	if g.CurrentCard.Type.Has(deck.SKIP) {
		g.EndTurn(false)
	}

	if g.CurrentCard.Type.Has(deck.REVERSE) {
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
		for _, p := range g.Players {
			p.AddCard(g.Deck.Draw())
		}
	}
}

func (g *Game) PlayCard(c *deck.Card) {
	g.logger.Trace().Msg("Playing card")

	if g.HasPendingCatorce() {
		p := g.GetPlayer(g.PlayerCatorce)
		g.logger.Trace().Str("player_name", p.Name).Msg("There's a pending catorce!")
		p.CatorcesMissed += 1

		for i := 0; i < 4; i++ {
			card := g.Deck.Draw()
			p.AddCard(card)
		}
		g.PlayerCatorce = 0
	}

	g.CurrentPlayer().CardsPlayed += 1
	turnDuration := time.Since(g.TurnStarted)
	g.CurrentPlayer().AddDuration(turnDuration)

	if turnDuration > g.LargestResponseTime {
		g.LargestResponseTime = turnDuration
	}

	g.Deck.Discard(g.CurrentCard)
	g.CurrentCard = c

	jump := false

	if c.Type.Has(deck.WILD) {
		g.logger.Debug().Str("from", string(g.State)).Str("to", "CHOOSE_COLOR").Msg("Changing state")
		g.State = CHOOSE_COLOR
	}

	if c.Type.Has(deck.SKIP) {
		jump = true
	}

	if c.Type.Has(deck.DRAW) {
		g.DrawCount += c.Value
	}

	if c.Type.Has(deck.SWAP) {
		// Don't enter swap state if the game will be over
		if len(g.CurrentPlayer().Hand) != 0 {
			g.logger.Debug().Str("from", string(g.State)).Str("to", "CHOOSE_PLAYER").Msg("Changing state")
			g.State = CHOOSE_PLAYER // TODO: Possible conflict in states if a card is a WILD SWAP, check
		}
	}

	if c.Type.Has(deck.REVERSE) {
		if g.PlayerAmount() != 2 {
			g.Reverse()
		} else {
			jump = true
		}
	}

	// TODO: SKIPALL DISCARDALL etc

	g.EndTurn(jump)
}

func (g *Game) DrawCard() {
	g.logger.Trace().Msg("Drawing a card")

	if g.HasPendingCatorce() {
		p := g.GetPlayer(g.PlayerCatorce)
		g.logger.Trace().Str("player_name", p.Name).Msg("There's a pending catorce!")
		for i := 0; i < 4; i++ {
			card := g.Deck.Draw()
			p.AddCard(card)
		}
		g.PlayerCatorce = 0
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
func (g *Game) EndTurn(skip bool) bool {
	g.Rounds += 1
	g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Int("rounds", g.Rounds).Msg("Ending turn")
	if len(g.CurrentPlayer().Hand) == 0 {
		g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Msg("Player has 0 cards")
		g.State = LOBBY
		return true
	}

	if len(g.CurrentPlayer().Hand) == 1 && g.CurrentCard.Type != deck.SWAP {
		g.logger.Trace().Int("pid", g.CurrentPlayer().ID).Msg("Player has 1 card left, setting catorce")
		g.PlayerCatorce = g.CurrentPlayer().ID
	}

	g.NextPlayer()
	if skip {
		g.NextPlayer()
	}

	g.logger.Debug().Str("from", string(g.State)).Str("to", "CHOOSE_CARD").Msg("Changing state")
	g.State = CHOOSE_CARD

	g.TurnStarted = time.Now()
	return false
}

func (g *Game) SwapHands(p1, p2 *Player) {
	p1.Hand, p2.Hand = p2.Hand, p1.Hand
}

func (g *Game) Reverse() {
	g.logger.Trace().Msg("Reversing game")
	slices.Reverse(g.Players)
}

func (g *Game) ChooseColor(c deck.Color) {
	// We change the card color to the chosen color, this only
	// affects special cards, so we don't see it as they are always black
	g.logger.Trace().Str("color", string(c)).Msg("Setting card color")
	g.CurrentCard.SetColor(c)
	g.EndTurn(false)
}

func (g *Game) NextPlayer() {
	g.logger.Trace().Msg("Moving to next player")
	p := g.Players[0]
	g.Players = g.Players[1:]
	g.Players = append(g.Players, p)
}

func (g *Game) GetCurrentCard() *deck.Card {
	return g.CurrentCard
}

func (g *Game) GetDeck() *deck.Deck {
	return g.Deck
}

func (g *Game) ResetDeck() {
	g.logger.Trace().Msg("Resetting deck")
	g.Deck = deck.New(g.Config.DeckConfig, false)
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
	return g.PlayerCatorce != 0
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

	return out.String()
}
