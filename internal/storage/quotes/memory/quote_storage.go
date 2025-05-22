package memory

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"quotes/internal/domain/models"
	"quotes/internal/services"
	"quotes/internal/storage"
)

type QuoteStorage struct {
	quotes []models.Quote
	mu     sync.RWMutex
	nextID int64
}

func NewQuoteStorage() services.QuoteRepository {
	return &QuoteStorage{
		quotes: make([]models.Quote, 0),
		nextID: 1,
	}
}

func (s *QuoteStorage) Create(quote *models.Quote) error {
	const op = "storage.quotes.memory.Create"

	if quote.Author == "" {
		return fmt.Errorf("%s: %w", op, storage.ErrEmptyAuthor)
	}
	if quote.Text == "" {
		return fmt.Errorf("%s: %w", op, storage.ErrEmptyText)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	quote.ID = s.nextID
	quote.CreatedAt = time.Now()
	s.quotes = append(s.quotes, *quote)
	s.nextID++
	return nil
}

func (s *QuoteStorage) GetAll() ([]models.Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	quotes := make([]models.Quote, len(s.quotes))
	copy(quotes, s.quotes)
	return quotes, nil
}

func (s *QuoteStorage) GetRandom() (*models.Quote, error) {
	const op = "storage.quotes.memory.GetRandom"

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.quotes) == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrNoQuotesAvailable)
	}
	randomIndex := rand.IntN(len(s.quotes))
	quote := s.quotes[randomIndex]
	return &quote, nil
}

func (s *QuoteStorage) GetByAuthor(author string) ([]models.Quote, error) {
	const op = "storage.quotes.memory.GetByAuthor"

	if author == "" {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrEmptyAuthor)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Quote
	for _, quote := range s.quotes {
		if quote.Author == author {
			result = append(result, quote)
		}
	}
	return result, nil
}

func (s *QuoteStorage) Delete(id int64) error {
	const op = "storage.quotes.memory.Delete"

	if id <= 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, quote := range s.quotes {
		if quote.ID == id {
			s.quotes[i] = s.quotes[len(s.quotes)-1]
			s.quotes = s.quotes[:len(s.quotes)-1]
			return nil
		}
	}
	return fmt.Errorf("%s: %w", op, storage.ErrQuoteNotFound)
}
