package bot

import (
	"fmt"

	"github.com/d-nery/catorce/pkg/deck"
	"github.com/d-nery/catorce/pkg/game"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Result builder for query responses
type ResultBuilder struct {
	results tb.Results
}

// Action stickers
const (
	DRAW_STICKER = "CAACAgEAAxkBAAICpmDKXqpoPbRhwByJkmbxq0bNWNx7AAJDAQACx2NQRmEvrW3ks82BHwQ"
	PASS_STICKER = "CAACAgEAAxkBAAICqGDKXqzVf0tPGtCn6Uk0FwHwut1AAALVAQACfLpIRuYgeD5SDQ2BHwQ"
)

// Result creates a new result builder
func Results() *ResultBuilder {
	return &ResultBuilder{
		results: make(tb.Results, 0),
	}
}

// Result creates a new result builder with the specified capacity
func ResultsWithCap(cap int) *ResultBuilder {
	return &ResultBuilder{
		results: make(tb.Results, 0, cap),
	}
}

// AddGameNotStarted adds an ArticleResult stating that a game wasn't started
func (rb *ResultBuilder) AddGameNotStarted() *ResultBuilder {
	rb.results = append(rb.results, &tb.ArticleResult{
		ResultBase: tb.ResultBase{ID: "nogame"},

		Title: "O jogo ainda não começou",
		Text:  "O jogo ainda não começou.\n/start para começar.",
	})

	return rb
}

// AddNotPlaying adds an ArticleResult stating that the user is not playing at the moment
func (rb *ResultBuilder) AddNotPlaying() *ResultBuilder {
	rb.results = append(rb.results, &tb.ArticleResult{
		ResultBase: tb.ResultBase{ID: "nogame"},

		Title: "Você não está jogando",
		Text:  "Você não está jogando no momento.\n/new para começar um nesse canal ou /join caso já tenha um jogo nesse grupo",
	})

	return rb
}

// AddGameInfo adds an StickerResult that shows current game info
func (rb *ResultBuilder) AddGameInfo(g *game.Game) *ResultBuilder {
	res := &tb.StickerResult{
		ResultBase: tb.ResultBase{ID: "gameinfo"},
	}

	res.SetContent(&tb.InputTextMessageContent{
		Text:      g.GameInfo(),
		ParseMode: tb.ModeMarkdown,
	})

	rb.results = append(rb.results, res)

	return rb
}

// AddDraw adds an StickerResult with the Draw action
func (rb *ResultBuilder) AddDraw(amount int) *ResultBuilder {
	if amount == 0 {
		amount = 1
	}

	res := &tb.StickerResult{}
	res.Cache = DRAW_STICKER
	res.ID = "draw"
	res.SetContent(&tb.InputTextMessageContent{
		Text: fmt.Sprintf("Puxando %d carta(s)", amount),
	})

	rb.results = append(rb.results, res)
	return rb
}

// AddPass adds an StickerResult with the Pass action
func (rb *ResultBuilder) AddPass() *ResultBuilder {
	res := &tb.StickerResult{}
	res.Cache = PASS_STICKER
	res.ID = "pass"
	res.SetContent(&tb.InputTextMessageContent{
		Text: "Passando a vez",
	})

	rb.results = append(rb.results, res)
	return rb
}

// AddCard adds an StickerResult with a card
func (rb *ResultBuilder) AddCard(g *game.Game, c *deck.Card, can_play bool) *ResultBuilder {
	res := &tb.StickerResult{}

	if can_play {
		res.Cache = c.Sticker()
		res.ID = c.UID()
	} else {
		res.Cache = c.StickerNotAvailable()
		res.ID = fmt.Sprintf("cantplay:%s", c.UID())
		res.SetContent(&tb.InputTextMessageContent{
			Text:      g.GameInfo(),
			ParseMode: tb.ModeMarkdown,
		})
	}

	rb.results = append(rb.results, res)

	return rb
}

// AddCard adds an ArticleResult with a list of cards on the player's hand
func (rb *ResultBuilder) AddCurrentPlayerHand(g *game.Game) *ResultBuilder {
	res := &tb.ArticleResult{}
	res.Title = "Mão atual"
	res.ID = "hand"

	desc := ""
	if len(g.CurrentPlayer().Hand) == 0 {
		desc = "Vazia :)"
	} else {
		for _, c := range g.CurrentPlayer().Hand {
			desc += fmt.Sprintf("%s, ", c.StringPretty())
		}

		desc = desc[:len(desc)-2]
	}

	res.Description = desc
	res.SetContent(&tb.InputTextMessageContent{
		Text:      g.GameInfo(),
		ParseMode: tb.ModeMarkdown,
	})

	rb.results = append(rb.results, res)

	return rb
}

// AddColors adds a a list of ArticleResults with the possible colors
func (rb *ResultBuilder) AddColors() *ResultBuilder {
	for k, c := range deck.Colors {
		if c == deck.BLACK {
			continue
		}

		res := &tb.ArticleResult{}
		res.ID = fmt.Sprintf("color:%s", k)
		res.Title = "Escolha uma cor!"
		res.Description = deck.COLOR_ICONS[c]
		res.SetContent(&tb.InputTextMessageContent{
			Text: deck.COLOR_ICONS[c],
		})

		rb.results = append(rb.results, res)
	}

	return rb
}

// Results returns the underlying results
func (rb *ResultBuilder) Results() tb.Results {
	return rb.results
}
