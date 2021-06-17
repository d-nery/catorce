package deck

import (
	"fmt"
	"math/rand"
)

type Deck struct {
	Cards     []*Card
	Graveyard []*Card
}

func New() *Deck {
	// Thre are 108 cards in the official deck -> 25 each color + 8 black
	deck := Deck{
		Cards:     make([]*Card, 0, 108),
		Graveyard: make([]*Card, 0, 108),
	}

	for _, color := range Colors {
		if color == &BLACK {
			continue
		}

		for _, value := range CardValues {
			card := NewCard(color, value, nil)

			deck.Cards = append(deck.Cards, &card)

			if value != &ZERO {
				card2 := NewCard(color, value, nil)
				deck.Cards = append(deck.Cards, &card2)
			}
		}
	}

	for _, special := range SpecialCards {
		for i := 0; i < 4; i++ {
			card := NewCard(&BLACK, nil, special)
			deck.Cards = append(deck.Cards, &card)
		}
	}

	return &deck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Print() {
	fmt.Println("Cards:")

	for _, c := range d.Cards {
		fmt.Printf("%v\n", c)
	}

	fmt.Println("Graveyard:")

	for _, c := range d.Graveyard {
		fmt.Printf("%v\n", c)
	}
}

func (d *Deck) Available() int {
	return len(d.Cards)
}

func (d *Deck) Discarded() int {
	return len(d.Graveyard)
}

func (d *Deck) FillFromGraveyard() {
	if d.Discarded() == 0 {
		return
	}

	for len(d.Graveyard) > 0 {
		card := d.Graveyard[0]
		d.Cards = append(d.Cards, card)
		d.Graveyard = d.Graveyard[1:]
	}

	d.Shuffle()
}

func (d *Deck) Discard(c *Card) {
	// Return card to black when discarded
	if c.IsSpecial() {
		c.SetColor(&BLACK)
	}

	d.Graveyard = append(d.Graveyard, c)
}

func (d *Deck) Draw() *Card {
	if d.Available() == 0 {
		d.FillFromGraveyard()
	}

	if d.Available() == 0 {
		return nil
	}

	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}
