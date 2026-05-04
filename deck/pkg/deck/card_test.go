package deck

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ranks = []Rank{
	Ace,
	Two,
	Three,
	Four,
	Five,
	Six,
	Seven,
	Eight,
	Nine,
	Ten,
	Jack,
	Queen,
	King,
}

var suits = []Suit{
	Club,
	Diamond,
	Spade,
	Heart,
}

func TestCardTypeAssign(t *testing.T) {
	ace := Card{Rank: Ace, Suit: Club}

	assert.Equal(t, ace.Rank, Ace)
	assert.Equal(t, ace.Suit, Club)
}

func TestCardEqualsFunc(t *testing.T) {
	aceDiamonds := Card{Rank: Ace, Suit: Diamond}
	secondAceDiamonds := Card{Rank: Ace, Suit: Diamond}

	jackClubs := Card{Rank: Jack, Suit: Club}

	assert.True(t, aceDiamonds.Equals(aceDiamonds))
	assert.True(t, aceDiamonds.Equals(secondAceDiamonds))
	assert.False(t, aceDiamonds.Equals(jackClubs))
}

type suitCount struct {
	Diamonds int
	Clubs    int
	Spades   int
	Hearts   int
}

func countBySuit(deck []Card) []int {
	var suits []int = make([]int, 4)
	for _, card := range deck {
		suits[card.Suit-1] = suits[card.Suit-1] + 1
	}
	return suits
}

func countByRank(deck []Card) []int {
	var suits []int = make([]int, 13)
	for _, card := range deck {
		suits[card.Rank-1] = suits[card.Rank-1] + 1
	}
	return suits
}

func TestNewDeckFunc(t *testing.T) {
	deck := New()

	assert.Len(t, deck, 52)

	deckSuits := countBySuit(deck[:])

	for _, n := range deckSuits {
		assert.Equal(t, 13, n)
	}

	deckRanks := countByRank(deck[:])

	for _, n := range deckRanks {
		assert.Equal(t, 4, n)
	}
}

func TestEquals(t *testing.T) {
	aceOfSpades := Card{Suit: Spade, Rank: Ace}
	secondAceOfSpades := Card{Suit: Spade, Rank: Ace}
	jackOfHearths := Card{Suit: Heart, Rank: Jack}

	assert.True(t, aceOfSpades.Equals(secondAceOfSpades))
	assert.False(t, jackOfHearths.Equals(aceOfSpades))
}

func TestString(t *testing.T) {
	joker := Card{Suit: Joker}

	assert.Equal(t, "Joker", joker.String())

	aceOfHearts := Card{Suit: Heart, Rank: Ace}

	assert.Equal(t, "Ace of Hearts", aceOfHearts.String())
}

func TestFilterByRank(t *testing.T) {

	twoOfHearts := Card{Rank: Two, Suit: Heart}
	threeOfSpades := Card{Rank: Three, Suit: Spade}
	deckWithoutTwos := New(Filter(func(c Card) bool { return c.Rank == Two }))

	assert.NotContains(t, deckWithoutTwos, twoOfHearts)
	assert.Contains(t, deckWithoutTwos, threeOfSpades)
	assert.Len(t, deckWithoutTwos, 48)
}

func TestFilterBySuit(t *testing.T) {
	deckWithoutHearts := New(Filter(func(c Card) bool { return c.Suit == Heart }))

	kingOfHearts := Card{Rank: King, Suit: Heart}
	queenOfSpades := Card{Rank: Queen, Suit: Spade}

	assert.NotContains(t, deckWithoutHearts, kingOfHearts)
	assert.Contains(t, deckWithoutHearts, queenOfSpades)
	assert.Len(t, deckWithoutHearts, 39)
}

func TestFilterByRankAndSuit(t *testing.T) {
	filteredOut := New(Filter(func(c Card) bool {
		return slices.Contains([]Suit{Spade, Heart}, c.Suit) || slices.Contains([]Rank{Four, Five}, c.Rank)
	}))

	assert.NotContains(t, filteredOut, Card{Rank: Four, Suit: Spade})
	assert.NotContains(t, filteredOut, Card{Rank: Four, Suit: Heart})
	assert.NotContains(t, filteredOut, Card{Rank: Four, Suit: Diamond})
	assert.NotContains(t, filteredOut, Card{Rank: Five, Suit: Club})
	assert.Contains(t, filteredOut, Card{Rank: Six, Suit: Club})

	assert.Len(t, filteredOut, 22)
}

func TestCombineDecks(t *testing.T) {
	// len: 39
	deckWithoutSpades := New(Filter(func(c Card) bool {
		return c.Suit == Spade
	}))
	fmt.Printf("len(deckWithoutSpades): %d\n", len(deckWithoutSpades))

	combinedDeck := New(Deck(deckWithoutSpades))
	assert.Len(t, combinedDeck, 91)
}

func TestInsertJoker(t *testing.T) {
	deckWithThreeJokers := New(Jokers(3))

	assert.Len(t, deckWithThreeJokers, 55)
	assert.Contains(t, deckWithThreeJokers, Card{Suit: Joker})
}

// This is flaky
func TestShuffle(t *testing.T) {
	fullDeck := New()
	firstElement := fullDeck[0]
	Shuffle(&fullDeck)
	assert.NotEqual(t, firstElement, fullDeck[0])
}

func TestSortBySuit(t *testing.T) {

	fullDeck := New()

	// sort by suit
	Sort(fullDeck, func(c1, c2 *Card) bool {
		return c1.Suit < c2.Suit
	})

	assert.Equal(t, Spade, fullDeck[0].Suit)
	assert.Equal(t, Ace, fullDeck[0].Rank)
}
