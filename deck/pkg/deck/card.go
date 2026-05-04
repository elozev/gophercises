package deck

import (
	"fmt"
	"math/rand/v2"
	"sort"
)

type Suit int

const (
	_ Suit = iota
	Spade
	Diamond
	Club
	Heart
	Joker
)

type Rank int

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Rank Rank
	Suit Suit
}

func New(opts ...func([]Card) []Card) []Card {
	ret := make([]Card, 52)

	const cardsPerSuit = 13
	// without the joker
	const suits = 4

	cardIndex := 0
	for r := range cardsPerSuit {
		for s := range suits {
			ret[cardIndex] = Card{Rank: Rank(r + 1), Suit: Suit(s + 1)}
			cardIndex++
		}
	}

	for _, opt := range opts {
		ret = opt(ret)
	}

	return ret
}

func filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(s))

	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func FilterOut(deck []Card, ranks []Rank, suits []Suit) []Card {
	result := deck

	for _, r := range ranks {
		result = filter(result, func(c Card) bool { return c.Rank != r })
	}

	for _, s := range suits {
		result = filter(result, func(c Card) bool { return c.Suit != s })
	}

	return result
}

func (c *Card) Equals(o Card) bool {
	return c.Rank == o.Rank && c.Suit == o.Suit
}

func (c *Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

func Combine(decks ...[]Card) []Card {
	var combined []Card
	for _, deck := range decks {
		combined = append(combined, deck...)
	}

	return combined
}

func Deck(decks ...[]Card) func(cards []Card) []Card {
	return func(cards []Card) []Card {
		var combined = cards
		for _, deck := range decks {
			combined = append(combined, deck...)
		}

		return combined
	}
}

func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		result := make([]Card, 0, len(cards))
		for _, c := range cards {
			if !f(c) {
				result = append(result, c)
			}
		}
		return result
	}
}

func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := range n {
			cards = append(cards, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}
		return cards
	}
}

func Shuffle(d *[]Card) {
	for i := range len(*d) {
		swapIndex := rand.IntN(len(*d))
		(*d)[i], (*d)[swapIndex] = (*d)[swapIndex], (*d)[i]
	}
}

type LessFunc func(c1, c2 *Card) bool

func Sort(d []Card, less LessFunc) {
	sort.Slice(d, func(i, j int) bool {
		return less(&d[i], &d[j])
	})
}
