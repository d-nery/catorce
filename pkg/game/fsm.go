package game

import (
	"errors"
	"fmt"

	"github.com/d-nery/catorce/pkg/deck"
)

type GameState string

const (
	LOBBY        GameState = "LOBBY"
	CHOOSE_CARD  GameState = "CHOOSE_CARD"
	DREW         GameState = "DREW"
	CHOOSE_COLOR GameState = "CHOOSE_COLOR"
	CATORCE      GameState = "CATORCE"
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
	Color  *deck.Color
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
	g.logger.Debug().Str("event", fmt.Sprintf("%T", evt)).Str("current_state", string(g.state)).Msg("New event received")

	switch e := evt.(type) {
	case *EvtStartGame:
		if g.state != LOBBY {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtStartGame")
			return ErrEventNotCovered
		}

		if g.PlayerAmount() < 2 {
			g.logger.Trace().Msg("ErrNotEnoughPlayers for EvtStartGame")
			return ErrNotEnoughPlayers
		}

		g.deck.Shuffle()
		g.ShufflePlayers()
		g.DistributeCards()
		g.PlayFirstCard()

		g.logger.Debug().Str("from", string(g.state)).Str("to", "CHOOSE_CARD").Msg("Changing state")
		g.state = CHOOSE_CARD
		return nil

	case *EvtAddPlayer:
		if g.state != LOBBY {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtAddPlayer")
			return ErrEventNotCovered
		}

		if g.PlayerAmount() >= 10 {
			g.logger.Trace().Msg("ErrMaxPlayers for EvtAddPlayer")
			return ErrMaxPlayers
		}

		g.AddPlayer(e.Player)
		return nil

	case *EvtCardPlayed:
		if g.state != CHOOSE_CARD && g.state != DREW {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtCardPlayed")
			return ErrEventNotCovered
		}

		if e.Player != g.CurrentPlayer() {
			g.logger.Trace().Msg("ErrWrongPlayer for EvtCardPlayed")
			return ErrWrongPlayer
		}

		c := e.Card
		if !c.CanPlayOnTop(g.CurrentCard(), g.DrawCounter() > 0) {
			g.logger.Trace().Str("card", c.String()).Str("current", g.CurrentCard().String()).Msg("ErrCantPlayCard for EvtCardPlayed")
			return ErrCantPlayCard
		}

		e.Player.RemoveCard(c)
		g.PlayCard(c)

		return nil

	case *EvtDrawCard:
		if g.state != CHOOSE_CARD {
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
		if g.state != DREW {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtPass")
			return ErrEventNotCovered
		}

		g.EndTurn()
		return nil

	case *EvtColorChosen:
		if g.state != CHOOSE_COLOR {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtColorChosen")
			return ErrEventNotCovered
		}

		if !g.CurrentCard().IsSpecial() {
			g.logger.Trace().Msg("ErrCantChooseColor for EvtColorChosen")
			return ErrCantChooseColor
		}

		g.current_card.SetColor(e.Color)
		if !g.EndTurn() {
			g.logger.Debug().Str("from", string(g.state)).Str("to", "CHOOSE_CARD").Msg("Changing state")
			g.state = CHOOSE_CARD
		}
		return nil

	case *EvtCatorce:
		if g.state != CHOOSE_CARD {
			g.logger.Trace().Msg("ErrEventNotCovered for EvtCatorce")
			return ErrEventNotCovered
		}

		if !g.HasPendingCatorce() {
			g.logger.Trace().Msg("ErrNoCatorcePending for EvtCatorce")
			return ErrNoCatorcePending
		}

		g.player_catorce = nil
		return nil

	default:
		g.logger.Trace().Msg("ErrUnknownEvent")
		return ErrUnknownEvent
	}
}
