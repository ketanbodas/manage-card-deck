package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ketanbodas/manage-card-deck/deck"
)

/*
This file contains the code to start the http server and
map the various endpoints to appropriate methods

endpoints:
1. create new deck
2. open deck
3. draw cards
*/

/*

Some error codes:
1 => query parameter "shuffle" has incorrect value (api: create new deck)
2 => wrong card code provided (api: create new deck)
3 => query parameter "deck_id" not provided (api: open deck or draw cards)
4 => query parameter "deck_id" has invalid value (api: open deck)
5 => query parameter "count" not provided (api: draw cards)
6 => query parameter "count" has invalid value (api: draw cards)
7 => error while drawing hand

*/

// common types which are used to form rest api responses
type deckMetadata struct {
	Id        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type cardsList struct {
	Cards []deck.Card `json:"cards"`
}

// rest api responses
type newDeckResponse struct {
	deckMetadata
}

type openDeckResponse struct {
	deckMetadata
	cardsList
}

type drawHandResponse struct {
	cardsList
}

// struct to return error response
type errorMessage struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"error"`
}

// route apis and start server
func StartServer() {
	router := setupRouter()
	router.Run("localhost:3000")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/deck", newDeck)
	router.GET("/deck/open", openDeck)
	router.GET("/deck/draw", drawCards)
	return router
}

// create and return new deck
func newDeck(c *gin.Context) {
	shuffleQueryParam := c.DefaultQuery("shuffle", "false")
	shuffle, e := strconv.ParseBool(shuffleQueryParam)
	if e != nil {
		message := fmt.Sprintf("Invalid query param value for 'shuffle': %v", shuffleQueryParam)
		em := errorMessage{Message: message, ErrorCode: 1}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}
	cards := c.Query("cards")
	deck, error := deck.CreateNewDeck(shuffle, cards)
	if error != nil {
		message := fmt.Sprintf("error in deck creation: %v", error)
		em := errorMessage{Message: message, ErrorCode: 2}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	metadata := deckMetadata{
		Id:        deck.DeckId.String(),
		Shuffled:  deck.Shuffled,
		Remaining: len(deck.Cards),
	}
	response := newDeckResponse{metadata}

	c.IndentedJSON(http.StatusOK, response)
}

// open existing deck
func openDeck(c *gin.Context) {
	deckId := c.Query("deck_id")
	if len(deckId) == 0 {
		em := errorMessage{Message: "'deck_id' query param not provided", ErrorCode: 3}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}
	deck, error := deck.OpenDeck(deckId)
	if error != nil {
		message := fmt.Sprintf("Error in opening deck: %v", error)
		em := errorMessage{Message: message, ErrorCode: 4}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	deckMetadata := deckMetadata{
		Id:        deck.DeckId.String(),
		Shuffled:  deck.Shuffled,
		Remaining: len(deck.Cards),
	}

	cardsList := cardsList{
		Cards: deck.Cards,
	}

	response := openDeckResponse{
		deckMetadata: deckMetadata,
		cardsList:    cardsList,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// draw cards from existing deck
func drawCards(c *gin.Context) {
	deckId := c.Query("deck_id")
	if len(deckId) == 0 {
		em := errorMessage{Message: "'deck_id' query param not provided", ErrorCode: 3}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	count := c.Query("count")
	if len(count) == 0 {
		em := errorMessage{Message: "'count' query param not provided", ErrorCode: 5}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	cardCount, error := strconv.ParseInt(count, 10, 32)
	if error != nil {
		message := fmt.Sprintf("count '%v' is not an integer", count)
		em := errorMessage{Message: message, ErrorCode: 6}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	if cardCount <= 0 {
		em := errorMessage{Message: "count must be greater than zero", ErrorCode: 6}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	hand, error := deck.DrawCards(deckId, int(cardCount))
	if error != nil {
		message := fmt.Sprintf("Error in drawing a hand from deck: %v", error)
		em := errorMessage{Message: message, ErrorCode: 7}
		c.IndentedJSON(http.StatusBadRequest, em)
		return
	}

	cardsList := cardsList{
		Cards: hand,
	}
	response := drawHandResponse{cardsList}
	c.IndentedJSON(http.StatusOK, response)
}
