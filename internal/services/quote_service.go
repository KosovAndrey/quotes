package services

import (
	"fmt"
	"quotes/internal/domain/models"
	"quotes/internal/storage"
)

type QuoteRepository interface {
	Create(quote *models.Quote) error
	GetAll() ([]models.Quote, error)
	GetRandom() (*models.Quote, error)
	GetByAuthor(author string) ([]models.Quote, error)
	Delete(id int64) error
}

type QuoteService struct {
	repo QuoteRepository
}

func NewQuoteService(repo QuoteRepository) *QuoteService {
	return &QuoteService{
		repo: repo,
	}
}

func (s *QuoteService) CreateQuote(quote *models.Quote) error {
	const op = "services.quote.CreateQuote"

	if quote == nil {
		return fmt.Errorf("%s: %w", op, fmt.Errorf("quote cannot be nil"))
	}

	if err := s.repo.Create(quote); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *QuoteService) GetAllQuotes() ([]models.Quote, error) {
	const op = "services.quote.GetAllQuotes"

	quotes, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return quotes, nil
}

func (s *QuoteService) GetRandomQuote() (*models.Quote, error) {
	const op = "services.quote.GetRandomQuote"

	quote, err := s.repo.GetRandom()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return quote, nil
}

func (s *QuoteService) GetQuotesByAuthor(author string) ([]models.Quote, error) {
	const op = "services.quote.GetQuotesByAuthor"

	if author == "" {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrEmptyAuthor)
	}

	quotes, err := s.repo.GetByAuthor(author)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return quotes, nil
}

func (s *QuoteService) DeleteQuote(id int64) error {
	const op = "services.quote.DeleteQuote"

	if id <= 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrInvalidID)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
