package game

import (
	"errors"
	"fmt"

	"github.com/d-nery/catorce/pkg/deck"
)

type GameState string

const (
	LOBBY         GameState = "LOBBY"
	CHOOSE_CARD   GameState = "CHOOSE_CARD"
	DREW          GameState = "DREW"
	CHOOSE_COLOR  GameState = "CHOOSE_COLOR"
	CHOOSE_PLAYER GameState = "CHOOSE_PLAYER"
	CATORCE       GameState = "CATORCE"
)

type EvtStartGame struct{}

type EvtCatorce struct {
	Player *Player
}

type EvtAddPlayer struct {
	Player *Player
}

type EvtCardPlayed struct {
	Player *Player
	Card   *deck.Card
}

type EvtColorChosen struct {
	Player *Player
	Color  deck.Color
}

type EvtPlayerSwapChosen struct {
	Player *Player
	Target int
}

type EvtDrawCard struct {
	Player *Player
}

type EvtPass struct {
	Player *Player
}

type EventError error

// Possible Event Errors
var (
	ErrNotEnoughPlayers EventError = errors.New("fsm: not enough players")
	ErrEventNotCovered  EventError = errors.New("fsm: event not covered in current state")
	ErrMaxPlayers       EventError = errors.New("fsm: maximum number of players reached")
	ErrWrongPlayer      EventError = errors.New("fsm: it's not this player turn")
	ErrCantPlayCard     EventError = errors.New("fsm: illegal card")
	ErrCantChooseColor  EventError = errors.New("fsm: current card is not special, can't change color")
	ErrNoCatorcePending EventError = errors.New("fsm: no catorces pending")
	ErrUnknownEvent     EventError = errors.New("fsm: unknown event")
)

func (g *Game) FireEvent(evt interface{}) EventError {
	g.logger.Debug().Str("event", fmt.Sprintf("%T", evt)).Str("current_state", string(g.State)).Msg("New event received")

	switch e := evt.(type) {
	case *EvtStartGame:
		if g.State != LOBBY {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtStartGame")
			return ErrEventNotCovered
		}

		if g.PlayerAmount() < 2 {
			g.logger.Trace().Msg("ErrNotEnoughPlayers for EvtStartGame")
			return ErrNotEnoughPlayers
		}

		g.ResetDeck()
		g.Deck.Shuffle()
		g.ShufflePlayers()
		g.DistributeCards()
		g.PlayFirstCard()

		g.logger.Debug().Str("from", string(g.State)).Str("to", "CHOOSE_CARD").Msg("Changing state")
		g.State = CHOOSE_CARD
		return nil

	case *EvtAddPlayer:
		if g.State != LOBBY {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtAddPlayer")
			return ErrEventNotCovered
		}

		g.AddPlayer(e.Player)
		return nil

	case *EvtCardPlayed:
		if g.State != CHOOSE_CARD && g.State != DREW {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtCardPlayed")
			return ErrEventNotCovered
		}

		if e.Player != g.CurrentPlayer() {
			g.logger.Trace().Msg("ErrWrongPlayer for EvtCardPlayed")
			return ErrWrongPlayer
		}

		c := e.Card
		if !c.CanPlayOnTop(g.GetCurrentCard(), g.DrawCounter() > 0, g.Config.StackConfig) {
			g.logger.Trace().Str("card", c.String()).Str("current", g.GetCurrentCard().String()).Msg("ErrCantPlayCard for EvtCardPlayed")
			return ErrCantPlayCard
		}

		e.Player.RemoveCard(c)
		g.PlayCard(c)

		return nil

	case *EvtDrawCard:
		if g.State != CHOOSE_CARD {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtDrawCard")
			return ErrEventNotCovered
		}

		p := e.Player

		if p != g.CurrentPlayer() {
			g.logger.Trace().Msg("ErrWrongPlayer for EvtDrawCard")
			return ErrWrongPlayer
		}

		g.DrawCard()
		return nil

	case *EvtPass:
		if g.State != DREW {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtPass")
			return ErrEventNotCovered
		}

		g.EndTurn(false, CHOOSE_CARD)
		return nil

	case *EvtColorChosen:
		if g.State != CHOOSE_COLOR {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtColorChosen")
			return ErrEventNotCovered
		}

		if !g.GetCurrentCard().IsSpecial() {
			g.logger.Trace().Msg("ErrCantChooseColor for EvtColorChosen")
			return ErrCantChooseColor
		}

		g.CurrentCard.SetColor(e.Color)
		g.EndTurn(false, CHOOSE_CARD)
		return nil

	case *EvtPlayerSwapChosen:
		if g.State != CHOOSE_PLAYER {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtColorChosen")
			return ErrEventNotCovered
		}

		target := g.GetPlayer(e.Target)
		g.SwapHands(g.CurrentPlayer(), target)

		g.EndTurn(false, CHOOSE_CARD)
		return nil

	case *EvtCatorce:
		if g.State != CHOOSE_CARD {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtCatorce")
			return ErrEventNotCovered
		}

		if !g.HasPendingCatorce() {
			g.logger.Trace().Msg("ErrNoCatorcePending for EvtCatorce")
			return ErrNoCatorcePending
		}

		if e.Player.ID != g.PlayerCatorce {
			g.logger.Trace().Msg("ErrWrongPlayer for EvtCatorce")
			return ErrWrongPlayer
		}

		g.PlayerCatorce = 0
		return nil

	default:
		g.logger.Trace().Msg("ErrUnknownEvent")
		return ErrUnknownEvent
	}
}
