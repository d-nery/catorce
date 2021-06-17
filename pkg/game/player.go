package game

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/d-nery/catorce/pkg/deck"
)

type Player struct {
	ID       int
	Name     string
	Username string
	Hand     []*deck.Card
}

func NewPlayer(id int, user *tb.User) *Player {
	return &Player{
		ID:       id,
		Name:     user.FirstName,
		Username: user.Username,
		Hand:     make([]*deck.Card, 0, 7),
	}
}

func (p *Player) AddCard(c *deck.Card) {
	p.Hand = append(p.Hand, c)
}

func (p *Player) RemoveCard(c *deck.Card) {
	for i, crd := range p.Hand {
		if c.UID() == crd.UID() {
			p.Hand[i] = p.Hand[len(p.Hand)-1]
			p.Hand = p.Hand[:len(p.Hand)-1]
			return
		}
	}
}

func (p *Player) NameWithMention() string {
	if p.Username == "" {
		return fmt.Sprintf("*%s*", p.Name)
	}

	return fmt.Sprintf("*%s* (@%s)", p.Name, p.Username)
}

func (p *Player) PrintHand() {
	fmt.Printf("[%s] Hand: ", p.Name)
	for _, c := range p.Hand {
		fmt.Printf("%+v ", c)
	}
	fmt.Println()
}
