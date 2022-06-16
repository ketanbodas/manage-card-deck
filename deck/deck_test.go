package deck

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSequentialFullDeck(t *testing.T) {
	// create new full sequential deck
	deck, error := CreateNewDeck(false, "")

	// assert no errors and deck is returned
	assert.Nil(t, error)
	assert.NotNil(t, deck)

	// assert 52 cards
	assert.Equal(t, 52, len(deck.Cards))

	// assert first and last cards (for sequence)
	assert.Equal(t, "AS", deck.Cards[0].Code)
	assert.Equal(t, "KH", deck.Cards[len(deck.Cards)-1].Code)

}

func TestNewShuffledFullDeck(t *testing.T) {
	// create new full sequential deck
	deck, error := CreateNewDeck(true, "")

	// assert no errors and deck is returned
	assert.Nil(t, error)
	assert.NotNil(t, deck)

	// assert 52 cards
	assert.Equal(t, 52, len(deck.Cards))
}

func TestNewSequentialPartialDeckValid(t *testing.T) {
	// create new deck from codes
	codes := "AS,KD,AC,2C,KH,10H"
	deck, error := CreateNewDeck(false, codes)

	// assert no errors and deck is returned
	assert.Nil(t, error)
	assert.NotNil(t, deck)

	// assert 6 cards
	assert.Equal(t, 6, len(deck.Cards))

	// assert first and last cards (for sequence)
	assert.Equal(t, "AS", deck.Cards[0].Code)
	assert.Equal(t, "10H", deck.Cards[len(deck.Cards)-1].Code)

}

func TestNewPartialDeckVInvalid(t *testing.T) {
	// some code combinations to test valid and invalid codes
	codesList := []string{
		"AS,11D,AC",
		"2C KH 10H",
		"D1",
		"  ",
		"aS",
		"As",
		"as",
	}

	for _, codes := range codesList {
		_, error := CreateNewDeck(true, codes)
		assert.NotNil(t, error)
		_, error = CreateNewDeck(false, codes)
		assert.NotNil(t, error)
	}

}

func TestOpenDeckValid(t *testing.T) {
	// create new full sequential deck
	d, _ := CreateNewDeck(false, "")
	deck_id := d.DeckId.String()

	// open deck
	deck, error := OpenDeck(deck_id)

	// assert no errors and deck is returned
	assert.Nil(t, error)
	assert.NotNil(t, deck)

	// assert 52 cards
	assert.Equal(t, 52, len(deck.Cards))

	// assert first and last cards (for sequence)
	assert.Equal(t, "AS", deck.Cards[0].Code)
	assert.Equal(t, "KH", deck.Cards[len(deck.Cards)-1].Code)

}

func TestNOpenPartialDeckValid(t *testing.T) {
	// create new deck from codes
	codes := "AS,KD,AC,2C,KH,10H"
	d, _ := CreateNewDeck(false, codes)
	deck_id := d.DeckId.String()

	// open deck
	deck, error := OpenDeck(deck_id)

	// assert no errors and deck is returned
	assert.Nil(t, error)
	assert.NotNil(t, deck)

	// assert 6 cards
	assert.Equal(t, 6, len(deck.Cards))

	// assert first and last cards (for sequence)
	assert.Equal(t, "AS", deck.Cards[0].Code)
	assert.Equal(t, "10H", deck.Cards[len(deck.Cards)-1].Code)

}

func TestOpenDeckInvalidUUID(t *testing.T) {
	_, error := OpenDeck("1234")
	assert.NotNil(t, error)
}

func TestOpenDeckUnknownUUID(t *testing.T) {
	_, error := OpenDeck(uuid.New().String())
	assert.NotNil(t, error)
}

func TestDrawCardInvalidUUID(t *testing.T) {
	_, error := DrawCards("1234", 4)
	assert.NotNil(t, error)

	_, error = DrawCards(uuid.New().String(), 6)
	assert.NotNil(t, error)
}

func TestDrawCardsSuccess(t *testing.T) {
	// new deck from codes
	codes := "AS,KD,AC,2C,KH,10H"
	d, _ := CreateNewDeck(false, codes)
	deck_id := d.DeckId.String()

	// open deck
	deck, error := OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 6, len(deck.Cards))
	assert.Equal(t, "AS", deck.Cards[0].Code)

	// draw cards
	cards, error := DrawCards(deck_id, 2)
	assert.Nil(t, error)
	assert.Equal(t, 2, len(cards))
	assert.Equal(t, "AS", cards[0].Code)
	assert.Equal(t, "KD", cards[1].Code)

	// open deck to verify remaining cards
	deck, error = OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 4, len(deck.Cards))
	assert.Equal(t, "AC", deck.Cards[0].Code)

	// draw some more cards
	cards, error = DrawCards(deck_id, 3)
	assert.Nil(t, error)
	assert.Equal(t, 3, len(cards))
	assert.Equal(t, "AC", cards[0].Code)
	assert.Equal(t, "2C", cards[1].Code)
	assert.Equal(t, "KH", cards[2].Code)

	// open deck to verify remaining cards
	deck, error = OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 1, len(deck.Cards))
	assert.Equal(t, "10H", deck.Cards[0].Code)

	// draw remaining cards
	cards, error = DrawCards(deck_id, 1)
	assert.Nil(t, error)
	assert.Equal(t, 1, len(cards))
	assert.Equal(t, "10H", cards[0].Code)

	// open deck to verify no remaining cards
	deck, error = OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 0, len(deck.Cards))
}

func TestDrawCardsCountZero(t *testing.T) {
	// new deck from codes
	codes := "AS,KD,AC,2C,KH,10H"
	d, _ := CreateNewDeck(false, codes)
	deck_id := d.DeckId.String()

	// open deck
	deck, error := OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 6, len(deck.Cards))
	assert.Equal(t, "AS", deck.Cards[0].Code)

	// cannot draw zero cards
	_, error = DrawCards(deck_id, 0)
	assert.NotNil(t, error)
}

func TestDrawCardsNoCardsLeft(t *testing.T) {
	codes := "AS,KD,AC,2C,KH,10H"
	d, _ := CreateNewDeck(false, codes)
	deck_id := d.DeckId.String()

	deck, error := OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 6, len(deck.Cards))
	assert.Equal(t, "AS", deck.Cards[0].Code)

	cards, error := DrawCards(deck_id, 6)
	assert.Nil(t, error)
	assert.Equal(t, 6, len(cards))

	deck, error = OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 0, len(deck.Cards))

	_, error = DrawCards(deck_id, 1)
	assert.NotNil(t, error)
}

func TestDrawCardsCountMoreThanCards(t *testing.T) {
	codes := "AS,KD,AC,2C,KH,10H"
	d, _ := CreateNewDeck(false, codes)
	deck_id := d.DeckId.String()

	deck, error := OpenDeck(deck_id)
	assert.Nil(t, error)
	assert.Equal(t, 6, len(deck.Cards))
	assert.Equal(t, "AS", deck.Cards[0].Code)

	_, error = DrawCards(deck_id, 10)
	assert.NotNil(t, error)
}

func TestValidateCardCodes(t *testing.T) {

	// validate different card codes

	assert.Nil(t, validateCardCode("AH"))
	assert.Nil(t, validateCardCode("KC"))
	assert.Nil(t, validateCardCode("JD"))
	assert.Nil(t, validateCardCode("QS"))
	assert.Nil(t, validateCardCode("10D"))
	assert.Nil(t, validateCardCode("4C"))
	assert.NotNil(t, validateCardCode("4 S"))
	assert.NotNil(t, validateCardCode("11S"))
	assert.NotNil(t, validateCardCode("10"))
	assert.NotNil(t, validateCardCode("4h"))
	assert.NotNil(t, validateCardCode("1S"))
	assert.NotNil(t, validateCardCode(""))
	assert.NotNil(t, validateCardCode("  "))
}
