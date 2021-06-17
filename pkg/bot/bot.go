package bot

import (
	"time"

	"github.com/d-nery/catorce/pkg/game"
	"github.com/rs/zerolog"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	tb      *tb.Bot
	games   map[int64]*game.Game // Maps chats to games
	players map[int]int64        // Maps players to chats

	catorceBtnMarkup *tb.ReplyMarkup
	logger           zerolog.Logger
}

func New(token string, logger zerolog.Logger) (*Bot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return &Bot{
		tb:      b,
		games:   make(map[int64]*game.Game),
		players: make(map[int]int64),

		logger: logger,
	}, nil
}

func (b *Bot) SetupHandlers() {
	b.catorceBtnMarkup = &tb.ReplyMarkup{}
	btnCatorce := b.catorceBtnMarkup.Data("CATORCE!", "catorce")
	b.catorceBtnMarkup.Inline(b.catorceBtnMarkup.Row(btnCatorce))

	b.tb.Handle("/new", b.HandleNew)
	b.tb.Handle("/join", b.HandleJoin)
	b.tb.Handle("/start", b.HandleStart)
	b.tb.Handle(tb.OnChosenInlineResult, b.HandleResult)
	b.tb.Handle(tb.OnQuery, b.HandleQuery)
	b.tb.Handle(&btnCatorce, b.HandleCatorce)

	// b.tb.Handle(tb.OnSticker, func(m *tb.Message) {
	// 	b.logger.Printf("STICKER %+v", m.Sticker)
	// 	b.tb.Send(m.Chat, m.Sticker.FileID)
	// })
}

func (b *Bot) Start() {
	b.tb.Start()
}
