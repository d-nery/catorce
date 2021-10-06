package game

import "github.com/d-nery/catorce/pkg/deck"

// Config holds game configuration, like extra cards, time limits etc
type Config struct {
	// Stacking +4 cards is possible in this game
	CanStackPlus4 bool

	// Swap hands cards is available in this game
	UseSpecialSwap bool

	DeckConfig deck.DeckConfig
}

func DefaultConfig() *Config {
	return &Config{
		CanStackPlus4:  false,
		UseSpecialSwap: false,
		DeckConfig: deck.DeckConfig{
			AmountOfJoker: 4,
			AmountOfDraw4: 4,

			AmountOfDraw2:   2,
			AmountOfSwap:    1,
			AmountOfSkip:    2,
			AmountOfReverse: 2,
		},
	}
}

func (g *Game) SetConfig(config *Config) {
	g.Config = config
}

func (g *Game) SetCanStackPlus4(can bool) {
	g.Config.CanStackPlus4 = can
}

func (g *Game) SetUseSpecialSwap(can bool) {
	g.Config.UseSpecialSwap = can
}

func (g *Game) CanStackPlus4() bool {
	return g.Config.CanStackPlus4
}

func (g *Game) UseSpecialSwap() bool {
	return g.Config.UseSpecialSwap
}
