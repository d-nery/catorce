package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/d-nery/catorce/pkg/game"
	"github.com/rs/zerolog"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	tb      *tb.Bot
	Games   map[int64]*game.Game // Maps chats to games
	Players map[int]int64        // Maps players to chats
	stats   OverallStats

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
		Games:   make(map[int64]*game.Game),
		Players: make(map[int]int64),
		stats:   make(OverallStats),

		logger: logger,
	}, nil
}

func (b *Bot) SetupHandlers() {
	b.catorceBtnMarkup = &tb.ReplyMarkup{}
	btnCatorce := b.catorceBtnMarkup.Data("CATORCE!", "catorce")
	b.catorceBtnMarkup.Inline(b.catorceBtnMarkup.Row(btnCatorce))

	b.tb.Handle("/new", b.HandleNew)
	b.tb.Handle("/help", b.HandleHelp)
	b.tb.Handle("/join", b.HandleJoin)
	b.tb.Handle("/start", b.HandleStart)
	b.tb.Handle("/stats", b.HandleStats)
	b.tb.Handle("/statsself", b.HandleSelfStats)
	b.tb.Handle(tb.OnChosenInlineResult, b.HandleResult)
	b.tb.Handle(tb.OnQuery, b.HandleQuery)
	b.tb.Handle(&btnCatorce, b.HandleCatorce)

	// b.tb.Handle(tb.OnSticker, func(m *tb.Message) {
	// 	b.logger.Printf("STICKER %+v", m.Sticker)
	// 	b.tb.Send(m.Chat, m.Sticker.FileID)
	// })
}

func (b *Bot) Load() {
	body, err := ioutil.ReadFile("data/data.json")

	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to load")
		return
	}

	err = json.Unmarshal(body, b)

	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to load")
		return
	}

	b.stats, err = ReadStatsFromFile("data/stats.json")

	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to load")
	}

	for _, g := range b.Games {
		g.SetLogger(b.logger)
	}
}

func (b *Bot) Persist() {
	body, err := json.Marshal(b)

	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to persist")
		return
	}

	err = ioutil.WriteFile("data/data.json", body, 0644)

	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to persist")
	}

	if b.stats.Persist("data/stats.json") != nil {
		b.logger.Error().Err(err).Msg("Failed to persist stats")
	}
}

func (b *Bot) Dump() {
	body, err := json.MarshalIndent(b, "", "  ")

	if err != nil {
		b.logger.Error().Err(err).Msg("Failed to persist")
		return
	}

	fmt.Println(string(body))

	if b.stats.Dump() != nil {
		b.logger.Error().Err(err).Msg("Failed to persist")
		return
	}
}

func (b *Bot) Start() {
	b.tb.Start()
}
