package main

import (
	"log"
	"net/http"

	"quotes/internal/handlers"
	"quotes/internal/services"
	"quotes/internal/storage/quotes/memory"

	"github.com/gorilla/mux"
)

func main() {
	repository := memory.NewQuoteStorage()
	quoteService := services.NewQuoteService(repository)
	quoteHandler := handlers.NewQuoteHandler(quoteService)

	r := mux.NewRouter()

	r.HandleFunc("/quotes", quoteHandler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", quoteHandler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes", quoteHandler.GetQuotesByAuthor).Methods("GET").Queries("author", "{author}")
	r.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.DeleteQuote).Methods("DELETE")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
