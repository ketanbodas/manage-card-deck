package api

import (
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ----------- Tests: New Deck  --------------

func TestNewDeckApiFullDeckNoShuffleSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck")
	assert.Equal(t, http.StatusOK, w.Code)
	body := extractNewDeckResponse(w)
	assertValidUUID(t, body.Id)
	assert.Equal(t, 52, body.Remaining)
	assert.Equal(t, false, body.Shuffled)
}

func TestNewDeckApiFullDeckWithShuffleSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=true")
	assert.Equal(t, http.StatusOK, w.Code)
	body := extractNewDeckResponse(w)
	assertValidUUID(t, body.Id)
	assert.Equal(t, 52, body.Remaining)
	assert.Equal(t, true, body.Shuffled)
}

func TestNewDeckApiFullDeckWithShuffleFalseSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=false")
	assert.Equal(t, http.StatusOK, w.Code)
	body := extractNewDeckResponse(w)
	assertValidUUID(t, body.Id)
	assert.Equal(t, 52, body.Remaining)
	assert.Equal(t, false, body.Shuffled)
}

func TestNewDeckApiFullDeckInvalidShuffleValueFailure(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=true1")
	assertBadRequestErrorCode(t, w, 1)
}

func TestNewDeckApiPartialDeckNoShuffleSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?cards=AS,10C")
	assert.Equal(t, http.StatusOK, w.Code)
	body := extractNewDeckResponse(w)
	assertValidUUID(t, body.Id)
	assert.Equal(t, 2, body.Remaining)
	assert.Equal(t, false, body.Shuffled)
}

func TestNewDeckApiPartialDeckWithShuffleSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=true&cards=AS,10C")
	assert.Equal(t, http.StatusOK, w.Code)
	body := extractNewDeckResponse(w)
	assertValidUUID(t, body.Id)
	assert.Equal(t, 2, body.Remaining)
	assert.Equal(t, true, body.Shuffled)
}

func TestNewDeckApiPartialDeckWithShuffleFalseSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=false&cards=AS,10C")
	assert.Equal(t, http.StatusOK, w.Code)
	body := extractNewDeckResponse(w)
	assertValidUUID(t, body.Id)
	assert.Equal(t, 2, body.Remaining)
	assert.Equal(t, false, body.Shuffled)
}

func TestNewDeckApiPartialDeckInvalidShuffleValueFailure(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=true1&cards=AS,10C")
	assertBadRequestErrorCode(t, w, 1)
}

func TestNewDeckApiPartialDeckWrongCodeFailure(t *testing.T) {
	w := runApi(http.MethodPost, "/deck?shuffle=false&cards=AS,10")
	assertBadRequestErrorCode(t, w, 2)
}

// ----------- Tests: Open Deck  --------------

func TestOpenDeckApiFullDeckSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck")
	assert.Equal(t, http.StatusOK, w.Code)
	uuid := extractNewDeckResponse(w).Id

	w = runApi(http.MethodGet, "/deck/open?deck_id="+uuid)
	body := extractOpenDeckResponse(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 52, body.Remaining)
	assert.Equal(t, false, body.Shuffled)
	assert.Equal(t, "AS", body.Cards[0].Code)
}

func TestOpenDeckApiNoDeckIdProvidedError(t *testing.T) {
	w := runApi(http.MethodGet, "/deck/open")
	assertBadRequestErrorCode(t, w, 3)
}

func TestOpenDeckApiDeckIdNotUUIDError(t *testing.T) {
	w := runApi(http.MethodGet, "/deck/open?deck_id="+"1234")
	assertBadRequestErrorCode(t, w, 4)
}

func TestOpenDeckApiDeckIdNotFoundError(t *testing.T) {
	w := runApi(http.MethodGet, "/deck/open?deck_id="+uuid.New().String())
	assertBadRequestErrorCode(t, w, 4)
}

func TestOpenDeckApiPartialDeckSuccess(t *testing.T) {
	// create new partial deck
	w := runApi(http.MethodPost, "/deck?cards=10H,KD,AS,JC,9C,2H")
	assert.Equal(t, http.StatusOK, w.Code)
	uuid := extractNewDeckResponse(w).Id

	// open deck
	w = runApi(http.MethodGet, "/deck/open?deck_id="+uuid)
	body := extractOpenDeckResponse(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 6, body.Remaining)
	assert.Equal(t, false, body.Shuffled)
	assert.Equal(t, "10H", body.Cards[0].Code)
}

// ----------- Tests: Draw cards  --------------

func TestDrawCardsApiFullDeckSuccess(t *testing.T) {
	w := runApi(http.MethodPost, "/deck")
	assert.Equal(t, http.StatusOK, w.Code)
	uuid := extractNewDeckResponse(w).Id
	path := "/deck/draw?deck_id=" + uuid + "&count=10"
	w = runApi(http.MethodGet, path)
	body := extractDrawCardsResponse(w)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 10, len(body.Cards))
	assert.Equal(t, "AS", body.Cards[0].Code)
}

func TestDrawCardsApiNoDeckIdProvidedError(t *testing.T) {
	path := "/deck/draw?count=10"
	w := runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 3)
}

func TestDrawCardsApiDeckIdNotUUIDError(t *testing.T) {
	path := "/deck/draw?deck_id=1234&count=10"
	w := runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 7)
}

func TestDrawCardsApiDeckIdNotFoundError(t *testing.T) {
	path := "/deck/draw?deck_id=" + uuid.New().String() + "&count=10"
	w := runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 7)
}

func TestDrawCardsApiCountNotProvidedError(t *testing.T) {
	w := runApi(http.MethodPost, "/deck")
	assert.Equal(t, http.StatusOK, w.Code)
	uuid := extractNewDeckResponse(w).Id
	path := "/deck/draw?deck_id=" + uuid
	w = runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 5)
}

func TestDrawCardsApiCountNotNumberError(t *testing.T) {
	w := runApi(http.MethodPost, "/deck")
	assert.Equal(t, http.StatusOK, w.Code)
	uuid := extractNewDeckResponse(w).Id
	path := "/deck/draw?deck_id=" + uuid + "&count=o"
	w = runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 6)
}

func TestDrawCardsApiCountNotPositiveError(t *testing.T) {
	w := runApi(http.MethodPost, "/deck")
	assert.Equal(t, http.StatusOK, w.Code)
	uuid := extractNewDeckResponse(w).Id
	path := "/deck/draw?deck_id=" + uuid + "&count=0"
	w = runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 6)
}

func TestDrawCardsApiEndToEnd(t *testing.T) {
	// create new deck of 6 cards
	w := runApi(http.MethodPost, "/deck?cards=AS,KD,AC,2C,KH,10H")
	assert.Equal(t, http.StatusOK, w.Code)
	newDeckRes := extractNewDeckResponse(w)
	assert.Equal(t, 6, newDeckRes.Remaining)
	uuid := newDeckRes.Id

	// draw 2 cards
	path := "/deck/draw?deck_id=" + uuid + "&count=2"
	w = runApi(http.MethodGet, path)
	body := extractDrawCardsResponse(w)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(body.Cards))
	assert.Equal(t, "AS", body.Cards[0].Code)
	assert.Equal(t, "KD", body.Cards[1].Code)

	// open deck and verify 4 remaining cards
	w = runApi(http.MethodGet, "/deck/open?deck_id="+uuid)
	openDeckRes := extractOpenDeckResponse(w)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 4, openDeckRes.Remaining)
	assert.Equal(t, "AC", openDeckRes.Cards[0].Code)

	// draw 4 cards
	path = "/deck/draw?deck_id=" + uuid + "&count=4"
	w = runApi(http.MethodGet, path)
	body = extractDrawCardsResponse(w)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 4, len(body.Cards))
	assert.Equal(t, "AC", body.Cards[0].Code)
	assert.Equal(t, "2C", body.Cards[1].Code)
	assert.Equal(t, "KH", body.Cards[2].Code)
	assert.Equal(t, "10H", body.Cards[3].Code)

	// open deck and verify no remaining cards
	w = runApi(http.MethodGet, "/deck/open?deck_id="+uuid)
	openDeckRes = extractOpenDeckResponse(w)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 0, openDeckRes.Remaining)

	// verify that cannot draw any more cards
	path = "/deck/draw?deck_id=" + uuid + "&count=1"
	w = runApi(http.MethodGet, path)
	assertBadRequestErrorCode(t, w, 7)
}

// ----------- Helper functions --------------

func runApi(method string, path string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	router.ServeHTTP(w, req)
	return w
}

func assertValidUUID(t *testing.T, uuidStr string) {
	assert.NotNil(t, uuidStr)
	_, e := uuid.Parse(uuidStr)
	assert.Nil(t, e)

}

func assertBadRequestErrorCode(t *testing.T, w *httptest.ResponseRecorder, errorCode int) {
	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := extractErrorResponse(w)
	assert.Equal(t, errorCode, body.ErrorCode)

}

func extractErrorResponse(w *httptest.ResponseRecorder) errorMessage {
	res := w.Body.String()
	body := errorMessage{}
	json.Unmarshal([]byte(res), &body)
	return body
}

func extractNewDeckResponse(w *httptest.ResponseRecorder) newDeckResponse {
	res := w.Body.String()
	body := newDeckResponse{}
	json.Unmarshal([]byte(res), &body)
	return body
}

func extractOpenDeckResponse(w *httptest.ResponseRecorder) openDeckResponse {
	res := w.Body.String()
	body := openDeckResponse{}
	json.Unmarshal([]byte(res), &body)
	return body
}

func extractDrawCardsResponse(w *httptest.ResponseRecorder) drawHandResponse {
	res := w.Body.String()
	body := drawHandResponse{}
	json.Unmarshal([]byte(res), &body)
	return body
}
