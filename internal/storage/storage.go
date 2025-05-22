package storage

import "errors"

var (
	ErrQuoteNotFound     = errors.New("quote not found")
	ErrEmptyAuthor       = errors.New("author cannot be empty")
	ErrEmptyText         = errors.New("quote text cannot be empty")
	ErrNoQuotesAvailable = errors.New("no quotes available")
	ErrInvalidID         = errors.New("invalid quote ID")
)
