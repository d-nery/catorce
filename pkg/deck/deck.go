package deck

import (
	"fmt"
	"math/rand"
)

// Deck is a group of cards and discarded cards
type Deck struct {
	Cards     []*Card
	Graveyard []*Card
}

// New creates a new filled deck
func New() *Deck {
	// Thre are 108 cards in the official deck -> 25 each color + 8 black
	deck := Deck{
		Cards:     make([]*Card, 0, 108),
		Graveyard: make([]*Card, 0, 108),
	}

	for _, color := range Colors {
		if color == BLACK {
			continue
		}

		for _, value := range CardValues {
			card := NewCard(color, value, SINVALID)

			deck.Cards = append(deck.Cards, &card)

			if value != ZERO {
				card2 := NewCard(color, value, SINVALID)
				deck.Cards = append(deck.Cards, &card2)
			}
		}
	}

	for _, special := range SpecialCards {
		for i := 0; i < 4; i++ {
			card := NewCard(BLACK, VINVALID, special)
			deck.Cards = append(deck.Cards, &card)
		}
	}

	return &deck
}

// Shuffle shuffles all the cards in the deck
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// Print prints the deck to the terminal
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

// Available returns the number of cards available in the deck
func (d *Deck) Available() int {
	return len(d.Cards)
}

// Discarded returns the number of cards in the graveyard (discarded)
func (d *Deck) Discarded() int {
	return len(d.Graveyard)
}

// FillFromGraveyard removes all cards from the graveyard and puts the back on the deck
// The deck is shuffled afterwards
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

// Discard adds card c to the graveyard
func (d *Deck) Discard(c *Card) {
	// Return card to black when discarded
	if c.IsSpecial() {
		c.SetColor(BLACK)
	}

	d.Graveyard = append(d.Graveyard, c)
}

// Draw removes a card from the deck and returns it
// If the deck is empty, it tries to fill itself from the graveyard
// If the deck is still empty, return nil
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
