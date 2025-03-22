package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

var books []Book

func main() {
	books = []Book{}
	r := mux.NewRouter()
	r.HandleFunc("/books", bookHandler).Methods("GET", "POST", "PUT")
	r.HandleFunc("/books/{id:[0-9]+}", bookHandler).Methods("GET", "DELETE")
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	//GET http://localhost:8080/books
	//GET http://localhost:8080/books/1
	case "GET":
		var responseJSON []byte
		var err error
		var book Book
		var id int

		idParam := mux.Vars(r)["id"]

		if idParam != "" {
			id, err = strconv.Atoi(idParam)
			if err != nil {
				http.Error(w, "Error converting id to int", http.StatusInternalServerError)
				return
			}
			book, _ = filterById(id)
			if book.Id == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			responseJSON, err = json.Marshal(book)

		} else {
			responseJSON, err = json.Marshal(books)
		}

		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(responseJSON)

	/* POST http://localhost:8080/books
	body	{
			"Title" : "Head First Go",
			"Author" : "Jay McGavren",
			"Pages" : 556
		}
	*/
	case "POST":
		//TODO validate input
		var newBook Book
		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newBook.Id = len(books) + 1
		books = append(books, newBook)

		responseJSON, err := json.Marshal(newBook)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the JSON response
		_, _ = w.Write(responseJSON)

	/*PUT http://localhost:8080/books
	body	{
				"Id" : 2,
				"Title" : "The Go Programming Language",
				"Author" : "Alan Donovan,Brian Kernighan",
				"Pages" : 400
			}
	*/
	case "PUT":
		//TODO validate input
		var bookToUpdate Book
		err := json.NewDecoder(r.Body).Decode(&bookToUpdate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if bookToUpdate.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		book, i := filterById(bookToUpdate.Id)
		if book.Id == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if bookToUpdate.Title != "" {
			book.Title = bookToUpdate.Title
		}

		if bookToUpdate.Author != "" {
			book.Author = bookToUpdate.Author
		}

		if bookToUpdate.Pages != 0 {
			book.Pages = bookToUpdate.Pages
		}

		books[i] = book
		responseJSON, err := json.Marshal(book)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		// Set the appropriate headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the JSON response
		_, _ = w.Write(responseJSON)

		//DELETE http://localhost:8080/books/4
	case "DELETE":
		idParam := mux.Vars(r)["id"]

		if idParam == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			w.Write([]byte("Id conversion issue" + err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		book, i := filterById(id)
		if book.Id == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		//Delete book keeping original order
		books = append(books[:i], books[i+1:]...)
		w.WriteHeader(http.StatusNoContent)
	}
}

func filterById(id int) (Book, int) {
	var i int
	var book Book
	for i, book = range books {
		if book.Id == id {
			return book, i
		}
	}
	return Book{}, i
}
