package game

import "github.com/d-nery/catorce/pkg/deck"

// Config holds game configuration
type Config struct {
	DeckConfig  deck.DeckConfig
	StackConfig deck.StackConfig
}

func DefaultConfig() *Config {
	return &Config{
		DeckConfig: deck.DeckConfig{
			Cards: map[deck.CardData]int{
				{deck.RED, deck.NUMBER, 0}:   1,
				{deck.RED, deck.NUMBER, 1}:   2,
				{deck.RED, deck.NUMBER, 2}:   2,
				{deck.RED, deck.NUMBER, 3}:   2,
				{deck.RED, deck.NUMBER, 4}:   2,
				{deck.RED, deck.NUMBER, 5}:   2,
				{deck.RED, deck.NUMBER, 6}:   2,
				{deck.RED, deck.NUMBER, 7}:   2,
				{deck.RED, deck.NUMBER, 8}:   2,
				{deck.RED, deck.NUMBER, 9}:   2,
				{deck.RED, deck.NUMBER, 9}:   2,
				{deck.RED, deck.SKIP, -1}:    2,
				{deck.RED, deck.REVERSE, -1}: 2,
				{deck.RED, deck.DRAW, 2}:     2,

				{deck.BLUE, deck.NUMBER, 0}:   1,
				{deck.BLUE, deck.NUMBER, 1}:   2,
				{deck.BLUE, deck.NUMBER, 2}:   2,
				{deck.BLUE, deck.NUMBER, 3}:   2,
				{deck.BLUE, deck.NUMBER, 4}:   2,
				{deck.BLUE, deck.NUMBER, 5}:   2,
				{deck.BLUE, deck.NUMBER, 6}:   2,
				{deck.BLUE, deck.NUMBER, 7}:   2,
				{deck.BLUE, deck.NUMBER, 8}:   2,
				{deck.BLUE, deck.NUMBER, 9}:   2,
				{deck.BLUE, deck.NUMBER, 9}:   2,
				{deck.BLUE, deck.SKIP, -1}:    2,
				{deck.BLUE, deck.REVERSE, -1}: 2,
				{deck.BLUE, deck.DRAW, 2}:     2,

				{deck.GREEN, deck.NUMBER, 0}:   1,
				{deck.GREEN, deck.NUMBER, 1}:   2,
				{deck.GREEN, deck.NUMBER, 2}:   2,
				{deck.GREEN, deck.NUMBER, 3}:   2,
				{deck.GREEN, deck.NUMBER, 4}:   2,
				{deck.GREEN, deck.NUMBER, 5}:   2,
				{deck.GREEN, deck.NUMBER, 6}:   2,
				{deck.GREEN, deck.NUMBER, 7}:   2,
				{deck.GREEN, deck.NUMBER, 8}:   2,
				{deck.GREEN, deck.NUMBER, 9}:   2,
				{deck.GREEN, deck.NUMBER, 9}:   2,
				{deck.GREEN, deck.SKIP, -1}:    2,
				{deck.GREEN, deck.REVERSE, -1}: 2,
				{deck.GREEN, deck.DRAW, 2}:     2,

				{deck.YELLOW, deck.NUMBER, 0}:   1,
				{deck.YELLOW, deck.NUMBER, 1}:   2,
				{deck.YELLOW, deck.NUMBER, 2}:   2,
				{deck.YELLOW, deck.NUMBER, 3}:   2,
				{deck.YELLOW, deck.NUMBER, 4}:   2,
				{deck.YELLOW, deck.NUMBER, 5}:   2,
				{deck.YELLOW, deck.NUMBER, 6}:   2,
				{deck.YELLOW, deck.NUMBER, 7}:   2,
				{deck.YELLOW, deck.NUMBER, 8}:   2,
				{deck.YELLOW, deck.NUMBER, 9}:   2,
				{deck.YELLOW, deck.NUMBER, 9}:   2,
				{deck.YELLOW, deck.SKIP, -1}:    2,
				{deck.YELLOW, deck.REVERSE, -1}: 2,
				{deck.YELLOW, deck.DRAW, 2}:     2,

				{deck.BLACK, deck.WILD, -1}:            4,
				{deck.BLACK, deck.WILD | deck.DRAW, 4}: 4,
			},
		},
		StackConfig: deck.StackConfig{
			CanStackDraws:  false,
			CanStackWild:   false,
			CanStackBigger: false,
		},
	}
}

func (g *Game) SetConfig(config *Config) {
	g.Config = config
}
