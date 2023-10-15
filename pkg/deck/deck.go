package deck

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// Deck is a group of cards and discarded cards
type Deck struct {
	Cards     []*Card
	Graveyard []*Card

	Config DeckConfig
}

type CardData struct {
	Color
	CardType
	Value int
}

type DeckConfig struct {
	Cards map[CardData]int
}

// New creates a new filled deck
func New(config DeckConfig, half_deck bool) *Deck {
	divider := 1

	if half_deck {
		divider = 2
	}

	deck := Deck{
		Cards:     make([]*Card, 0),
		Graveyard: make([]*Card, 0),

		Config: config,
	}

	for card, amount := range config.Cards {
		for i := 0; i < amount/divider; i++ {
			deck.Cards = append(deck.Cards, NewCard(card.Color, card.CardType, card.Value))
		}
	}

	return &deck
}

// Merge adds other deck's cards to this deck
func (d *Deck) Merge(other *Deck) {
	d.Cards = append(d.Cards, other.Cards...)
	other.Cards = nil
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

// FillFromGraveyard removes all cards from the graveyard and puts them back on the deck
// The deck is shuffled afterwards
func (d *Deck) FillFromGraveyard() {
	if d.Discarded() == 0 {
		return
	}

	d.Cards = append(d.Cards, d.Graveyard...)
	d.Graveyard = []*Card{}

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
// If the deck is still empty, it refills itself with a half deck, increase the total amount of cards in game
func (d *Deck) Draw() *Card {
	if d.Available() == 0 {
		d.FillFromGraveyard()
	}

	if d.Available() == 0 {
		d.Merge(New(d.Config, true))
		d.Shuffle()
	}

	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

func (d *DeckConfig) MarshalJSON() ([]byte, error) {
	var cardData = []struct {
		Color
		CardType
		Value  int
		Amount int
	}{}

	for k, v := range d.Cards {
		cardData = append(cardData, struct {
			Color
			CardType
			Value  int
			Amount int
		}{
			k.Color, k.CardType, k.Value, v,
		})
	}

	return json.Marshal(&cardData)
}

func (d *DeckConfig) UnmarshalJSON(data []byte) error {
	var aux = []struct {
		Color
		CardType
		Value  int
		Amount int
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	d.Cards = map[CardData]int{}

	for _, e := range aux {
		d.Cards[CardData{e.Color, e.CardType, e.Value}] = e.Amount
	}

	return nil
}
