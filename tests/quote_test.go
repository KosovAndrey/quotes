package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"quotes/internal/domain/models"
	"quotes/internal/handlers"
	"quotes/internal/services"
	"quotes/internal/storage/quotes/memory"

	"github.com/gorilla/mux"
)

// setupTestServer создает тестовый сервер с настроенными маршрутами
func setupTestServer() *mux.Router {
	storage := memory.NewQuoteStorage()
	quoteService := services.NewQuoteService(storage)
	quoteHandler := handlers.NewQuoteHandler(quoteService)

	r := mux.NewRouter()
	r.HandleFunc("/quotes", quoteHandler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", quoteHandler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes", quoteHandler.GetQuotesByAuthor).Methods("GET").Queries("author", "{author}")
	r.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.DeleteQuote).Methods("DELETE")

	return r
}

// TestCreateQuote проверяет создание новой цитаты
func TestCreateQuote(t *testing.T) {
	router := setupTestServer()

	quote := models.Quote{
		Author: "Test Author",
		Text:   "Test Quote",
	}
	body, _ := json.Marshal(quote)

	req, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var response models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if response.Author != quote.Author {
		t.Errorf("handler returned unexpected author: got %v want %v",
			response.Author, quote.Author)
	}
}

// TestCreateQuoteEmptyFields проверяет создание цитаты с пустыми полями
func TestCreateQuoteEmptyFields(t *testing.T) {
	router := setupTestServer()

	testCases := []struct {
		name     string
		quote    models.Quote
		wantCode int
	}{
		{
			name:     "Empty Author",
			quote:    models.Quote{Text: "Test Quote"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty Text",
			quote:    models.Quote{Author: "Test Author"},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.quote)
			req, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.wantCode)
			}
		})
	}
}

// TestGetAllQuotes проверяет получение всех цитат
func TestGetAllQuotes(t *testing.T) {
	router := setupTestServer()

	quote := models.Quote{
		Author: "Test Author",
		Text:   "Test Quote",
	}
	body, _ := json.Marshal(quote)

	req, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	req, _ = http.NewRequest("GET", "/quotes", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var quotes []models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &quotes)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if len(quotes) == 0 {
		t.Error("handler returned empty quotes list")
	}
}

// TestGetRandomQuote проверяет получение случайной цитаты
func TestGetRandomQuote(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("GET", "/quotes/random", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for empty quotes: got %v want %v",
			status, http.StatusNotFound)
	}

	quote := models.Quote{
		Author: "Test Author",
		Text:   "Test Quote",
	}
	body, _ := json.Marshal(quote)

	req, _ = http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	req, _ = http.NewRequest("GET", "/quotes/random", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if response.Author != quote.Author {
		t.Errorf("handler returned unexpected author: got %v want %v",
			response.Author, quote.Author)
	}
}

// TestGetQuotesByAuthor проверяет получение цитат по автору
func TestGetQuotesByAuthor(t *testing.T) {
	router := setupTestServer()

	quote := models.Quote{
		Author: "Test Author",
		Text:   "Test Quote",
	}
	body, _ := json.Marshal(quote)

	req, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	req, _ = http.NewRequest("GET", "/quotes?author=Test+Author", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var quotes []models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &quotes)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if len(quotes) == 0 {
		t.Error("handler returned empty quotes list")
	}

	if quotes[0].Author != "Test Author" {
		t.Errorf("handler returned unexpected author: got %v want %v",
			quotes[0].Author, "Test Author")
	}
}

// TestGetQuotesByAuthorEmpty проверяет получение цитат с пустым автором
func TestGetQuotesByAuthorEmpty(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("GET", "/quotes?author=", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for empty author: got %v want %v",
			status, http.StatusOK)
	}
}

// TestDeleteQuote проверяет удаление цитаты
func TestDeleteQuote(t *testing.T) {
	router := setupTestServer()

	quote := models.Quote{
		Author: "Test Author",
		Text:   "Test Quote",
	}
	body, _ := json.Marshal(quote)

	req, _ := http.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	req, _ = http.NewRequest("DELETE", "/quotes/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	req, _ = http.NewRequest("GET", "/quotes", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var quotes []models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &quotes)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if len(quotes) != 0 {
		t.Error("quote was not deleted")
	}
}

// TestDeleteNonExistentQuote проверяет удаление несуществующей цитаты
func TestDeleteNonExistentQuote(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("DELETE", "/quotes/999", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent quote: got %v want %v",
			status, http.StatusNotFound)
	}
}

// TestDeleteInvalidQuoteID проверяет удаление цитаты с некорректным ID
func TestDeleteInvalidQuoteID(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("DELETE", "/quotes/invalid", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for invalid ID: got %v want %v",
			status, http.StatusNotFound)
	}
}

// TestCreateQuoteInvalidJSON проверяет создание цитаты с некорректным JSON
func TestCreateQuoteInvalidJSON(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("POST", "/quotes", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestGetQuotesByAuthorNotFound проверяет получение цитат несуществующего автора
func TestGetQuotesByAuthorNotFound(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("GET", "/quotes?author=NonexistentAuthor", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var quotes []models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &quotes)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if len(quotes) != 0 {
		t.Error("handler returned non-empty quotes list for nonexistent author")
	}
}

// TestGetAllQuotesEmpty проверяет получение всех цитат при пустом хранилище
func TestGetAllQuotesEmpty(t *testing.T) {
	router := setupTestServer()

	req, _ := http.NewRequest("GET", "/quotes", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var quotes []models.Quote
	err := json.Unmarshal(rr.Body.Bytes(), &quotes)
	if err != nil {
		t.Errorf("failed to unmarshal: %v", err)
	}

	if len(quotes) != 0 {
		t.Error("handler returned non-empty quotes list for empty storage")
	}
}
