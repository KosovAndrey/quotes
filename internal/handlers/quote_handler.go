package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"quotes/internal/domain/models"
	"quotes/internal/storage"

	"github.com/gorilla/mux"
)

type QuoteService interface {
	CreateQuote(quote *models.Quote) error
	GetAllQuotes() ([]models.Quote, error)
	GetRandomQuote() (*models.Quote, error)
	GetQuotesByAuthor(author string) ([]models.Quote, error)
	DeleteQuote(id int64) error
}

type QuoteHandler struct {
	service QuoteService
}

func NewQuoteHandler(service QuoteService) *QuoteHandler {
	return &QuoteHandler{
		service: service,
	}
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.quote.CreateQuote"

	var quote models.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		log.Printf("%s: failed to decode request body: %v", op, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateQuote(&quote); err != nil {
		log.Printf("%s: failed to create quote: %v", op, err)
		switch {
		case errors.Is(err, storage.ErrEmptyAuthor):
			http.Error(w, "Author cannot be empty", http.StatusBadRequest)
		case errors.Is(err, storage.ErrEmptyText):
			http.Error(w, "Quote text cannot be empty", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		log.Printf("%s: failed to encode response: %v", op, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) GetAllQuotes(w http.ResponseWriter, _ *http.Request) {
	const op = "handlers.quote.GetAllQuotes"

	quotes, err := h.service.GetAllQuotes()
	if err != nil {
		log.Printf("%s: failed to get quotes: %v", op, err)
		http.Error(w, "Failed to get quotes", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		log.Printf("%s: failed to encode response: %v", op, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, _ *http.Request) {
	const op = "handlers.quote.GetRandomQuote"

	quote, err := h.service.GetRandomQuote()
	if err != nil {
		log.Printf("%s: failed to get random quote: %v", op, err)
		switch {
		case errors.Is(err, storage.ErrNoQuotesAvailable):
			http.Error(w, "No quotes available", http.StatusNotFound)
		default:
			http.Error(w, "Failed to get random quote", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(quote); err != nil {
		log.Printf("%s: failed to encode response: %v", op, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) GetQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.quote.GetQuotesByAuthor"

	author := r.URL.Query().Get("author")
	if author == "" {
		log.Printf("%s: author parameter is missing", op)
		http.Error(w, "Author parameter is required", http.StatusBadRequest)
		return
	}

	quotes, err := h.service.GetQuotesByAuthor(author)
	if err != nil {
		log.Printf("%s: failed to get quotes by author: %v", op, err)
		switch {
		case errors.Is(err, storage.ErrEmptyAuthor):
			http.Error(w, "Author cannot be empty", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to get quotes by author", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		log.Printf("%s: failed to encode response: %v", op, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.quote.DeleteQuote"

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("%s: invalid quote ID: %v", op, err)
		http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteQuote(id); err != nil {
		log.Printf("%s: failed to delete quote: %v", op, err)
		switch {
		case errors.Is(err, storage.ErrQuoteNotFound):
			http.Error(w, "Quote not found", http.StatusNotFound)
		case errors.Is(err, storage.ErrInvalidID):
			http.Error(w, "Invalid quote ID", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to delete quote", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
