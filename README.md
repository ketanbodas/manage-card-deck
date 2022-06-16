# manage-card-deck

This repository provides a module to manage deck of cards via rest endpoints.   
  
Code is divided into two packages -  
1. **deck** - This package contains the types and functions to manage decks
2. **api**  - This package contains the [gin](https://github.com/gin-gonic/gin) based http server which provides endpoints to manage deck of cards 

Test cases (>95% coverage) are written using [testify](https://github.com/stretchr/testify)

  
### Supported endpoints  
1. Create new deck of cards
2. Open an existing deck of cards
3. Draw a hand from an existing deck of cards

#### Create New Deck
Endpoint: `localhost:3000/deck`  
Query Parameters: 
1. shuffle - a boolean indicating whether deck should be shuffled or not. Optional, default value is false
2. cards - comma separated list of card codes. Optional. If not provided all 52 cards would be added to deck

Example:  
Invoking `http://localhost:3000/deck?shuffle=true&cards=AS,KD,AC,2C,KH,10D` returns a json containing deck details:

    {
        "deck_id": "4c0c167a-5ba6-4437-a09d-9dcb7748df44",
        "shuffled": true,
        "remaining": 6
    }


#### Open Deck
Endpoint: `localhost:3000/deck/open`  
Query Parameters: 
1. deck_id - UUID which is returned when new deck is created.  This parameter is mandatory and must represent a valid UUID  

Example:  
Invoking `http://localhost:3000/deck/open?deck_id=4c0c167a-5ba6-4437-a09d-9dcb7748df44` returns following response:

    {
        "deck_id": "4c0c167a-5ba6-4437-a09d-9dcb7748df44",
        "shuffled": true,
        "remaining": 6,
        "cards": [
            {
                "value": "ACE",
                "suit": "CLUBS",
                "code": "AC"
            },
            {
                "value": "KING",
                "suit": "DIMONDS",
                "code": "KD"
            },
            {
                "value": "2",
                "suit": "CLUBS",
                "code": "2C"
            },
            {
                "value": "10",
                "suit": "DIMONDS",
                "code": "10D"
            },
            {
                "value": "ACE",
                "suit": "SPADES",
                "code": "AS"
            },
            {
                "value": "KING",
                "suit": "HEARTS",
                "code": "KH"
            }
        ]
    }

#### Draw Cards
Endpoint: `localhost:3000/deck/draw`  
Query Parameters: 
1. deck_id - UUID which is returned when new deck is created. This parameter is mandatory and must represent a valid UUID  
2. count - number of cards to be drawn. This parameter is mandatory and must be a positive integer with value not more than number of remaining cards in the deck  

Example:  
Invoking `http://localhost:3000/deck/draw?deck_id=4c0c167a-5ba6-4437-a09d-9dcb7748df44&count=4` returns a json containing a hand:

    {
        "cards": [
            {
                "value": "ACE",
                "suit": "CLUBS",
                "code": "AC"
            },
            {
                "value": "KING",
                "suit": "DIMONDS",
                "code": "KD"
            },
            {
                "value": "2",
                "suit": "CLUBS",
                "code": "2C"
            },
            {
                "value": "10",
                "suit": "DIMONDS",
                "code": "10D"
            }
        ]
    }
    
Note that, after above call, if open deck is called, it would return remaining cards as 2.  

#### Error Codes:
Above endpoints will throw error if input parameters are not right or if attempt is made to draw more cards than possible. The error codes are as follows:  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 1 => query parameter *shuffle* has incorrect value   
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 2 => query parameter *codes* has atleast one wrong card code provided  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 3 => query parameter *deck_id* not provided  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 4 => query parameter *deck_id* has invalid value  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 5 => query parameter *count* not provided  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 6 => query parameter *count* has invalid value  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 7 => error while drawing hand    

Some sample error responses:  
  
*deck_id has invalid value:*  
 
      {
          "errorCode": 4,
          "error": "Error in opening deck: deck not found for the input uuid 4c0c167a-5ba6-4437-a09d-9dcb7748df43"
      }

*hand attempts to draw more number of cards that what are remaining:*    

      {
          "errorCode": 7,
          "error": "Error in drawing a hand from deck: cannot draw 8 cards, deck has only 2"
      }


### Execution

#### To get the repository locally:
1. On local machine, clone this repository

#### How to build ?
1. Change directory to repository root
2. Execute `go build`

#### How to test ?
1. Change directory to repository root
2. Execute `go test -v ./...`

To test specific packages, go to specific package and run `go test`  

#### How to run ?
1. Change directory to repository root
2. Execute `go run .`
3. When this happens, server will start listening on port 3000
4. APIs can then be called using curl or postman
  
    
To summarize, to simply download this and start a server:  
`git clone https://github.com/ketanbodas/manage-card-deck`  
`cd manage-card-deck`  
`go run .`  


#### To use this module in another project as dependency:
1. run `go get github.com/ketanbodas/manage-card-deck` to get the module
2. add the package api to the import as `import "github.com/ketanbodas/manage-card-deck/api"`
3. run `go mod tidy`
4. use the package functions in code
5. run `go build` and verify that there are no errors


### Further improvements:
1. Decks are not currently persisted. Once server stops, all decks are lost. This can be added by saving the decks to file and loading the file when server starts back
2. Multithreading is not supported. This can cause problems when hands are drawn simultenously. 
3. Code can be optimized to use a single instance of cards. Currently, for each new deck, a new set of cards is created.
4. For now card codes are case sensitive. This can be improved.
5. No checks done for duplicate card codes while creating a deck of cards from given input. This can be improved with proper use case.
