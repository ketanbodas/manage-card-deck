package deck

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// type to represent a single Card
type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

// type to represent a Deck of cards
type Deck struct {
	DeckId   uuid.UUID
	Cards    []Card
	Shuffled bool
}

// list of suits and values
var cardSuits = []string{"SPADES", "DIMONDS", "CLUBS", "HEARTS"}
var cardValues = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

// map to get suit name from initial letter
var suitNames = map[string]string{
	"S": "SPADES",
	"D": "DIMONDS",
	"C": "CLUBS",
	"H": "HEARTS",
}

// map to get card value from initial letter
var valueNames = map[string]string{
	"A":  "ACE",
	"2":  "2",
	"3":  "3",
	"4":  "4",
	"5":  "5",
	"6":  "6",
	"7":  "7",
	"8":  "8",
	"9":  "9",
	"10": "10",
	"J":  "JACK",
	"Q":  "QUEEN",
	"K":  "KING",
}

// stores map of uuid and deck
var generatedDecks = map[uuid.UUID]Deck{}

/*
Creates a new deck of cards based on input arguments
inputs:
	shuffle :  whether deck is shuffled or not
	codes 	:  comma separated list of card codes which can be used to initilize the deck
			   if code is empty, all 52 cards are added
returns:
	a newly created deck
	error if any card code is invalid
*/
func CreateNewDeck(shuffle bool, codes string) (Deck, error) {
	var d Deck
	var e error
	if len(codes) == 0 {
		d = newSequentialDeck()
	} else {
		d, e = newDeckFromCodes(codes)
		if e != nil {
			return d, e
		}
	}
	if shuffle {
		d.shuffle()
	}
	d.DeckId = uuid.New()
	d.Shuffled = shuffle
	generatedDecks[d.DeckId] = d
	return d, e
}

/*
Returns an existing deck of cards based on input card uuid
inputs:
	deckId :  a UUID in string format
returns:
	an existing deck with given UUID
	error if UUID is not valid or deck not found
*/
func OpenDeck(deckId string) (Deck, error) {
	var d Deck
	uuid, error := parseUUID(deckId)
	if error != nil {
		return d, error
	}

	if !checkUUIDExists(uuid) {
		message := fmt.Sprintf("deck not found for the input uuid %v", deckId)
		return d, errors.New(message)
	}

	return generatedDecks[uuid], nil
}

/*
Returns a hand of cards with specified number of elements
inputs:
	deckId :  a UUID in string format
	count  :  number of cards to draw
returns:
	a slice of cards representing a hand
	error if UUID is not valid or deck not found or sufficient cards not available
*/
func DrawCards(deckId string, count int) ([]Card, error) {
	var cards []Card

	if count <= 0 {
		return cards, errors.New("count must be more than zero")
	}

	deck, error := OpenDeck(deckId)
	if error != nil {
		return cards, error
	}

	cards = deck.Cards

	if len(cards) == 0 {
		return cards, errors.New("cannot draw any cards, deck is empty")
	}

	if count > len(cards) {
		message := fmt.Sprintf("cannot draw %d cards, deck has only %d", count, len(cards))
		return cards, errors.New(message)
	}
	hand, cards := deck.Cards[:count], deck.Cards[count:]
	deck.Cards = cards
	generatedDecks[deck.DeckId] = deck
	return hand, error
}

/*
Prints the deck contents using deck as receiver
*/
func (d Deck) Print() {
	fmt.Printf("is_shuffled = %v remaining_cards = %v\n", d.Shuffled, len(d.Cards))
	for index, card := range d.Cards {
		fmt.Println(index, card)
	}
}

/*
Creates and returns a deck of cards initilized with incoming card codes
Returns error if any code is invalid
*/
func newDeckFromCodes(codes string) (Deck, error) {
	var d Deck
	deckCards := []Card{}
	codeList := strings.Split(codes, ",")
	for _, code := range codeList {
		code = strings.TrimSpace(code)
		error := validateCardCode(code)

		if error != nil {
			return d, error
		}
		l := len(code)
		c := Card{Value: valueNames[code[:l-1]], Suit: suitNames[code[l-1:]], Code: code}
		deckCards = append(deckCards, c)

	}
	d.Cards = deckCards
	return d, nil
}

/*
Valdates card code and returns error if any check fails
Code is valid iff none of below fails:
1. has a valid suit
2. has a valid value
*/
func validateCardCode(code string) error {
	code = strings.TrimSpace(code)
	if len(code) < 2 {
		message := fmt.Sprintf("code %v is invalid", code)
		return errors.New(message)
	}

	suit := code[len(code)-1:]
	_, suitExists := suitNames[suit]
	if !suitExists {
		message := fmt.Sprintf("code %v is invalid, should have proper suit name", code)
		return errors.New(message)
	}

	value := code[:len(code)-1]
	_, valueExists := valueNames[value]
	if !valueExists {
		message := fmt.Sprintf("code %v is invalid, should have proper card value", code)
		return errors.New(message)
	}

	return nil
}

/*
Creates and returns a deck of 52 cards in sequential order
*/
func newSequentialDeck() Deck {
	var d Deck
	deckCards := []Card{}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			code := value[:1] + suit[:1]
			deckCards = append(deckCards, Card{valueNames[value], suit, code})
		}
	}
	d.Cards = deckCards
	return d
}

/*
Shuffles the cards in the receiver deck
*/
func (d Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := range d.Cards {
		pos := rand.Intn(len(d.Cards))
		d.Cards[i], d.Cards[pos] = d.Cards[pos], d.Cards[i]
	}
}

/*
Returns true if UUID exists in the map of generated decks, false otherwise
*/
func checkUUIDExists(uuid uuid.UUID) bool {
	_, exists := generatedDecks[uuid]
	return exists
}

/*
Returns a UUID from input string or error if parsing fails
*/
func parseUUID(uuidStr string) (uuid.UUID, error) {
	uuid, error := uuid.Parse(uuidStr)
	if error != nil {
		message := fmt.Sprintf("input UUID '%v' is not a valid UUID4 value", uuidStr)
		return uuid, errors.New(message)
	}
	return uuid, nil
}
