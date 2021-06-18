package bot

import (
	"fmt"
	"strings"

	"github.com/d-nery/catorce/pkg/deck"
	"github.com/d-nery/catorce/pkg/game"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) HandleHelp(m *tb.Message) {
	b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("New help request received")
	helpMsg := `CATORCE!

2. No grupo, comece uma nova partida com /new
C. Para se juntar a uma partida use /join
4ª O jogo deve ter pelo menos 2 e no máximo 10 jogadores antes de começar
* Para jogar. Digite @catorce_uno_bot na caixa de mensagens ou clique no via @catorce_uno_bot ao lado das mensagens. Aguarde um pouco e você verá suas cartas. Cartas cinzas não podem ser jogadas. Se você não estiver na sua vez, todas as cartas serão cinzas.
110 Selecionar uma carta cinza irá mostrar a atual situação do jogo.
7- Ao ficar com uma unica carta sobrando, lembre-se de apertar no CATORCE!

Jogadores não podem entrar após a partida começar. Caso um jogador demore mais que 6 horas pra jogar ele é um babaca.
Caso o bot entre em colapso, não se preocupe, o estado do jogo é salvo e ao reiniciar, o bot recupera esse savepoint ;)

Outros comandos:
/stats - Mostra dados sobre os jogos do grupo interessantes
/statsself - Mostra seus dados apenas
/kill - F bot (adm only) (NYI)`

	b.tb.Send(m.Chat, helpMsg)
}

func (b *Bot) HandleNew(m *tb.Message) {
	b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("New game request received")

	if !m.FromGroup() {
		b.tb.Send(m.Sender, "Esse comando só funciona em grupos!")
		return
	}

	if _, ok := b.Games[m.Chat.ID]; ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("Game already exists")
		b.tb.Send(m.Chat, "Já tem um jogo rolando nesse chat!")
		return
	}

	b.Games[m.Chat.ID] = game.New(m.Chat.ID, b.logger)

	b.logger.Info().Int64("chat_id", m.Chat.ID).Msg("New game created")
	b.logger.Trace().Int("games_len", len(b.Games)).Send()

	if _, ok := b.stats[m.Chat.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Msg("No stats for current chat, creating")
		b.stats[m.Chat.ID] = &ChatStats{
			Group:   GroupStats{},
			Players: make(map[int]*PlayerStats),
		}
	}

	b.tb.Send(m.Chat, "Jogo criado com sucesso!\n/join para entrar.")
}

func (b *Bot) HandleJoin(m *tb.Message) {
	b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("Join request received")

	if !m.FromGroup() {
		b.tb.Send(m.Sender, "Esse comando só funciona em grupos!")
		return
	}

	if _, ok := b.Games[m.Chat.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Msg("No game running on this chat")
		b.tb.Send(m.Chat, "Não há nenhum jogo nesse chat! /new para criar um")
		return
	}

	if _, ok := b.Players[m.Sender.ID]; ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("Player already in a game")
		b.tb.Send(m.Chat, "Você já está participando de algum jogo!")
		return
	}

	g := b.Games[m.Chat.ID]
	if err := g.FireEvent(&game.EvtAddPlayer{Player: game.NewPlayer(m.Sender.ID, m.Sender)}); err != nil {
		b.logger.Error().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Err(err).Send()
		switch err {
		case game.ErrMaxPlayers:
			b.tb.Send(m.Chat, "Máximo de jogadores atingido!")
			return
		default:
			b.tb.Send(m.Chat, "Erro :(")
		}
		return
	}

	b.Players[m.Sender.ID] = m.Chat.ID

	var out strings.Builder
	out.WriteString("Entrando no jogo... Jogadores atuais:\n")

	for _, p := range g.PlayerList() {
		fmt.Fprintf(&out, " • %s\n", p.Name)
	}

	b.tb.Send(m.Chat, out.String())
}

func (b *Bot) HandleStart(m *tb.Message) {
	b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("Start request received")

	if !m.FromGroup() {
		b.tb.Send(m.Sender, "Esse comando só funciona em grupos!")
		return
	}

	if _, ok := b.Games[m.Chat.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Msg("No game running on this chat")
		b.tb.Send(m.Chat, "Não há nenhum jogo nesse chat! /new para criar um")
		return
	}

	g := b.Games[m.Chat.ID]
	if err := g.FireEvent(&game.EvtStartGame{}); err != nil {
		b.logger.Error().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Err(err).Send()
		switch err {
		case game.ErrNotEnoughPlayers:
			b.tb.Send(m.Chat, "Não há jogadores suficientes! /join para entrar.")
			return
		case game.ErrEventNotCovered:
			b.tb.Send(m.Chat, "Opa, acho que o jogo já começou!")
			return
		default:
			b.tb.Send(m.Chat, "Erro :(")
		}
		return
	}

	b.tb.Send(m.Chat, "Começando!")
	b.tb.Send(m.Chat, g.CurrentCardSticker())
	b.tb.Send(m.Chat, fmt.Sprintf("Jogador(a) Atual: %s", g.CurrentPlayer().NameWithMention()), tb.ModeMarkdown)

	b.Persist()
}

func (b *Bot) HandleResult(c *tb.ChosenInlineResult) {
	b.logger.Info().Int("user_id", c.From.ID).Msg("New Inline Result received")
	b.logger.Trace().Msgf("CHOSE INLINE => %+v", c)

	if _, ok := b.Players[c.From.ID]; !ok {
		return
	}

	chat := b.Players[c.From.ID]
	if _, ok := b.Games[chat]; !ok {
		return
	}

	g := b.Games[chat]
	res_id := c.ResultID

	b.logger.Info().Int("user_id", c.From.ID).Int64("chat_id", chat).Str("chosen", res_id).Msg("Chosen result")
	if strings.HasPrefix(res_id, "cantplay") ||
		strings.HasPrefix(res_id, "nogame") ||
		strings.HasPrefix(res_id, "gameinfo") ||
		strings.HasPrefix(res_id, "hand") {
		return
	}

	player := g.GetPlayer(c.From.ID)

	if res_id == "draw" {
		catorce := g.CatorcePlayer()
		if err := g.FireEvent(&game.EvtDrawCard{Player: player}); err != nil {
			b.logger.Error().Err(err).Int64("chat_id", chat).Send()
			switch err {
			case game.ErrEventNotCovered:
			case game.ErrWrongPlayer:
				return
			default:
				b.tb.Send(&c.From, "Erro :(")
			}
			return
		}

		// If there was a catorce player and the cards were succesfully drawn
		// then the catorce'd player received four cards, we need to warn them
		if catorce != nil {
			b.tb.Send(&tb.Chat{ID: chat},
				fmt.Sprintf(
					"Oh no! 😱\n%s não chamou CATORCE! a tempo e pegou 4 cartas!",
					catorce.NameWithMention(),
				),
				tb.ModeMarkdown,
			)
		}
	} else if res_id == "pass" {
		if err := g.FireEvent(&game.EvtPass{Player: player}); err != nil {
			b.logger.Error().Err(err).Int64("chat_id", chat).Send()
			switch err {
			case game.ErrEventNotCovered:
			case game.ErrWrongPlayer:
				return
			default:
				b.tb.Send(&c.From, "Erro :(")
			}
			return
		}
	} else if strings.HasPrefix(res_id, "color:") {
		colorCode := strings.Split(res_id, ":")[1]
		color := deck.Colors[colorCode]

		if err := g.FireEvent(&game.EvtColorChosen{Color: color, Player: player}); err != nil {
			b.logger.Error().Err(err).Int64("chat_id", chat).Send()
			switch err {
			case game.ErrEventNotCovered:
			case game.ErrWrongPlayer:
				return
			default:
				b.tb.Send(&c.From, "Erro :(")
			}
			return
		}

		if g.HasPendingCatorce() {
			b.tb.Send(&tb.Chat{ID: chat}, "Última carta!", b.catorceBtnMarkup)
		}
	} else {
		var card *deck.Card

		for _, c := range player.Hand {
			if c.UID() == res_id {
				card = c
				break
			}
		}

		if card == nil {
			b.logger.Error().Msg("Couldn't find card on player hand")
			b.tb.Send(&tb.Chat{ID: chat}, "Não encontrei essa carta na sua mão!")
			return
		}

		catorce := g.CatorcePlayer()
		if err := g.FireEvent(&game.EvtCardPlayed{Card: card, Player: player}); err != nil {
			b.logger.Error().Err(err).Int64("chat_id", chat).Send()
			switch err {
			case game.ErrEventNotCovered:
			case game.ErrWrongPlayer:
				return
			case game.ErrCantPlayCard:
				b.tb.Send(&tb.Chat{ID: chat}, "Essa carta é inválida!")
				return
			default:
				b.tb.Send(&c.From, "Erro :(")
			}
			return
		}

		if g.HasPendingCatorce() {
			b.tb.Send(&tb.Chat{ID: chat}, "Última carta!", b.catorceBtnMarkup)
		}

		// If there was a catorce player and the card was succesfully played
		// then the catorce'd player received four cards, we need to warn them
		if catorce != nil {
			b.tb.Send(&tb.Chat{ID: chat},
				fmt.Sprintf(
					"Oh no! 😱\n%s não chamou CATORCE! a tempo e pegou 4 cartas!",
					catorce.NameWithMention(),
				),
				tb.ModeMarkdown,
			)
		}
	}

	// If we returned to lobby, then game is over
	if g.GetState() == game.LOBBY {
		b.tb.Send(&tb.Chat{ID: chat},
			fmt.Sprintf(
				"Jogo finalizado após %d rounds!!\nVitória de %s",
				g.Rounds,
				g.CurrentPlayer().NameWithMention(),
			),
			tb.ModeMarkdown,
		)

		b.logger.Trace().Msg("game returned to lobby, deleting")
		b.SaveGameStats(g)
		delete(b.Games, chat)
		for k := range b.Players {
			if b.Players[k] == chat {
				delete(b.Players, k)
			}
		}
		b.logger.Trace().Int("games_len", len(b.Games)).Send()

		b.Persist()
		return
	}

	b.tb.Send(&tb.Chat{ID: chat}, fmt.Sprintf("Próximo(a) jogador(a): %s", g.CurrentPlayer().NameWithMention()), tb.ModeMarkdown)

	if g.GetState() == game.CHOOSE_COLOR {
		b.tb.Send(&tb.Chat{ID: chat}, "Escolha uma cor!")
	}

	b.Persist()
}

func (b *Bot) HandleQuery(q *tb.Query) {
	b.logger.Info().Int("user_id", q.From.ID).Msg("New Query received")
	b.logger.Trace().Msgf("QUERY => %+v", q)

	results := Results()

	if chat, ok := b.Players[q.From.ID]; !ok {
		results.AddNotPlaying()
	} else if g, ok := b.Games[chat]; !ok {
		results.AddGameNotStarted()
		b.logger.Info().Int64("chat_id", chat).Msg("No game running on this chat")
		return
	} else {
		player := g.GetPlayer(q.From.ID)

		if player == nil {
			b.logger.Error().Int64("chat_id", chat).Int("pid", q.From.ID).Msg("Couldn't get player from game")
			return
		}

		if player.ID != g.CurrentPlayer().ID {
			for _, c := range player.Hand {
				results.AddCard(g, c, false)
			}
		} else if g.GetState() == game.CHOOSE_CARD || g.GetState() == game.DREW {
			if g.GetState() == game.CHOOSE_CARD {
				results.AddDraw(g.DrawCounter())
			} else if g.GetState() == game.DREW {
				results.AddPass()
			}

			for _, c := range player.Hand {
				can_play := c.CanPlayOnTop(g.GetCurrentCard(), g.DrawCounter() > 0)
				results.AddCard(g, c, can_play)
			}
		} else if g.GetState() == game.CHOOSE_COLOR {
			results.AddColors()
			results.AddCurrentPlayerHand(g)
		}
	}

	err := b.tb.Answer(q, &tb.QueryResponse{
		Results:    results.Results(),
		CacheTime:  1,
		IsPersonal: true,
	})

	if err != nil {
		b.logger.Error().Err(err).Send()
	}
}

func (b *Bot) HandleCatorce(c *tb.Callback) {
	m := c.Message
	b.logger.Info().Int("user_id", c.Sender.ID).Int64("chat", m.Chat.ID).Msg("New Handle Catorce")

	if _, ok := b.Games[m.Chat.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Msg("No game running on this chat")
		return
	}

	g := b.Games[m.Chat.ID]
	player := g.GetPlayer(c.Sender.ID)

	if err := g.FireEvent(&game.EvtCatorce{Player: player}); err != nil {
		b.logger.Error().Err(err).Int64("chat_id", m.Chat.ID).Send()
		switch err {
		case game.ErrWrongPlayer:
			return
		default:
			b.tb.Edit(m, "Última carta!\nNão chamou catorce a tempo :(")
		}
		return
	}

	b.tb.Respond(c, &tb.CallbackResponse{Text: "CATORCE!"})
	b.tb.Edit(m, fmt.Sprintf("Última carta!\n%s chamou CATORCE!", player.Name))
	player.CatorcesCalled += 1
	b.Persist()
}

func (b *Bot) HandleStats(m *tb.Message) {
	b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("New stats request received")

	if !m.FromGroup() {
		b.tb.Send(m.Sender, "Esse comando só funciona em grupos!")
		return
	}

	if _, ok := b.stats[m.Chat.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("No stats for this chat")
		b.tb.Send(m.Chat, "Não encontrei estatísticas para esse chat! Tente terminar um jogo primeiro")
		return
	}

	cs := b.stats[m.Chat.ID]
	gs := cs.Group
	// ps := b.stats[m.Chat.ID].Players

	_, err := b.tb.Send(m.Chat, gs.Report(), tb.ModeMarkdown)

	if err != nil {
		b.logger.Error().Err(err).Int64("chat_id", m.Chat.ID).Send()
	}

	_, err = b.tb.Send(m.Chat, cs.Ranking(), tb.ModeMarkdown)

	if err != nil {
		b.logger.Error().Err(err).Int64("chat_id", m.Chat.ID).Send()
	}
}

func (b *Bot) HandleSelfStats(m *tb.Message) {
	b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("New self stats request received")

	if !m.FromGroup() {
		b.tb.Send(m.Sender, "Esse comando só funciona em grupos!")
		return
	}

	if _, ok := b.stats[m.Chat.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("No stats for this chat")
		b.tb.Send(m.Chat, "Não encontrei estatísticas para esse chat! Tente terminar um jogo primeiro")
		return
	}

	ps := b.stats[m.Chat.ID].Players

	if _, ok := ps[m.Sender.ID]; !ok {
		b.logger.Info().Int64("chat_id", m.Chat.ID).Int("user_id", m.Sender.ID).Msg("No stats for this chat")
		b.tb.Send(m.Chat, "Não encontrei estatísticas para você esse chat! Tente terminar um jogo primeiro")
		return
	}

	_, err := b.tb.Send(m.Chat, ps[m.Sender.ID].Report(), tb.ModeMarkdown)

	if err != nil {
		b.logger.Error().Err(err).Int64("chat_id", m.Chat.ID).Send()
	}
}
