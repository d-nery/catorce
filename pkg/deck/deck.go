package deck

import (
	"fmt"
	"math/rand"
)

type Deck struct {
	cards     []*Card
	graveyard []*Card
}

func New() *Deck {
	// Thre are 108 cards in the official deck -> 25 each color + 8 black
	deck := Deck{
		cards:     make([]*Card, 0, 108),
		graveyard: make([]*Card, 0, 108),
	}

	for _, color := range Colors {
		if color == &BLACK {
			continue
		}

		for _, value := range CardValues {
			card := NewCard(color, value, nil)

			deck.cards = append(deck.cards, &card)

			if value != &ZERO {
				card2 := NewCard(color, value, nil)
				deck.cards = append(deck.cards, &card2)
			}
		}
	}

	for _, special := range SpecialCards {
		for i := 0; i < 4; i++ {
			card := NewCard(&BLACK, nil, special)
			deck.cards = append(deck.cards, &card)
		}
	}

	return &deck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Print() {
	fmt.Println("Cards:")

	for _, c := range d.cards {
		fmt.Printf("%v\n", c)
	}

	fmt.Println("Graveyard:")

	for _, c := range d.graveyard {
		fmt.Printf("%v\n", c)
	}
}

func (d *Deck) Available() int {
	return len(d.cards)
}

func (d *Deck) Discarded() int {
	return len(d.graveyard)
}

func (d *Deck) FillFromGraveyard() {
	if d.Discarded() == 0 {
		return
	}

	for len(d.graveyard) > 0 {
		card := d.graveyard[0]
		d.cards = append(d.cards, card)
		d.graveyard = d.graveyard[1:]
	}

	d.Shuffle()
}

func (d *Deck) Discard(c *Card) {
	// Return card to black when discarded
	if c.IsSpecial() {
		c.SetColor(&BLACK)
	}

	d.graveyard = append(d.graveyard, c)
}

func (d *Deck) Draw() *Card {
	if d.Available() == 0 {
		d.FillFromGraveyard()
	}

	if d.Available() == 0 {
		return nil
	}

	card := d.cards[0]
	d.cards = d.cards[1:]
	return card
}
