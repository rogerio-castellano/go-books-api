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
	r.HandleFunc("/books/", bookHandler).Methods("GET")
	r.HandleFunc("/books", bookHandler).Methods("GET", "POST", "PUT")
	r.HandleFunc("/books/{id}", bookHandler).Methods("GET", "DELETE")
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch r.Method {
	case "GET":
		handleGetBooks(w, r)

	case "POST":
		handlePostBook(w, r)

	case "PUT":
		handlePutBook(w, r)

	case "DELETE":
		handleDeleteBook(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetBooks(w http.ResponseWriter, r *http.Request) {
	//GET http://localhost:8080/books
	//GET http://localhost:8080/books/1
	var responseJSON []byte
	var err error
	var book Book
	var id int

	idParam := mux.Vars(r)["id"]

	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "The provided id is invalid. Please ensure it is a positive integer.", http.StatusBadRequest)
			return
		}
		book, _, err = filterById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
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
}

func filterById(id int) (Book, int, error) {
	for i, book := range books {
		if book.Id == id {
			return book, i, nil
		}
	}
	return Book{}, -1, fmt.Errorf("Book with id %d not found", id)
}

func handlePostBook(w http.ResponseWriter, r *http.Request) {
	/* POST http://localhost:8080/books
	body	{
			"Title" : "Head First Go",
			"Author" : "Jay McGavren",
			"Pages" : 556
		}
	*/
	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateBook(newBook) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newBook.Id = len(books) + 1
	books = append(books, newBook)

	responseJSON, err := json.Marshal(newBook)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusBadRequest)
		return
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response
	_, _ = w.Write(responseJSON)

}

func validateBook(book Book) bool {
	if book.Title == "" || book.Author == "" || book.Pages <= 0 {
		return false
	}
	return true
}

func handlePutBook(w http.ResponseWriter, r *http.Request) {
	/*PUT http://localhost:8080/books
	body	{
				"Id" : 2,
				"Title" : "The Go Programming Language",
				"Author" : "Alan Donovan,Brian Kernighan",
				"Pages" : 400
			}
	*/
	//TODO validate input
	var bookToUpdate Book
	err := json.NewDecoder(r.Body).Decode(&bookToUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if bookToUpdate.Id == 0 {
		http.Error(w, "The id was not provided.", http.StatusBadRequest)
		return
	}

	book, i, err := filterById(bookToUpdate.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

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

	if !validateBook(book) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
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
}

func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	//DELETE http://localhost:8080/books/4
	idParam := mux.Vars(r)["id"]

	if idParam == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "The provided id is invalid. Please ensure it is a positive integer.", http.StatusBadRequest)
		return
	}

	_, i, err := filterById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//Delete book keeping original order
	books = append(books[:i], books[i+1:]...)
	w.WriteHeader(http.StatusNoContent)
}
