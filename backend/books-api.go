package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func main() {

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/books", bookHandler).Methods("GET", "POST", "PUT", "OPTIONS")
	r.HandleFunc("/books/", bookHandler).Methods("GET", "POST", "PUT", "OPTIONS")
	r.HandleFunc("/books/{id}", bookHandler).Methods("GET", "DELETE", "PUT", "OPTIONS")
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
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
	var id int

	idParam := mux.Vars(r)["id"]

	if idParam != "" {
		id, err = strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "The provided id is invalid. Please ensure it is a positive integer.", http.StatusBadRequest)
			return
		}
		book, ok := getBookById(id)
		if !ok {
			http.Error(w, "The requested book id could not be found.", http.StatusNotFound)
			return
		}

		responseJSON, err = json.Marshal(book)

	} else {
		books, _ := getBooks()
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

	newBookId, err := insertBook(newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newBook.Id = newBookId

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

	book, ok := getBookById(bookToUpdate.Id)
	if !ok {
		http.Error(w, "The requested book id could not be found.", http.StatusNotFound)
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

	// books[i] = book
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

	_, ok := getBookById(id)
	if !ok {
		http.Error(w, "The requested book id could not be found.", http.StatusNotFound)
		return
	}

	//Delete book keeping original order
	// books = append(books[:i], books[i+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the origin of your frontend
		// allowedOrigins := map[string]bool{
		//     "http://example.com":       true,
		//     "http://anotherexample.com": true,
		// }
		origin := r.Header.Get("Origin")
		// if allowedOrigins[origin] {
		//     w.Header().Set("Access-Control-Allow-Origin", origin)
		// }
		fmt.Println(origin)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getBooks() ([]Book, error) {
	db := getDbConnection()
	defer db.Close()

	query := `SELECT id, title, author, pages FROM books`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books = make([]Book, 0)

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Pages)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func getBookById(id int) (Book, bool) {

	db := getDbConnection()
	defer db.Close()

	var book Book
	query := `SELECT id, title, author, pages FROM books WHERE id = $1`
	row := db.QueryRow(query, id)
	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.Pages)

	switch err {
	case nil:
		return book, true
	case sql.ErrNoRows:
		return Book{}, false
	default:
		panic(err)
	}
}

func insertBook(book Book) (int, error) {
	var newBookID int

	db := getDbConnection()
	defer db.Close()

	query := `INSERT INTO books (title, author, pages) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, book.Title, book.Author, book.Pages).Scan(&newBookID)
	if err != nil {
		return 0, err
	}
	fmt.Println(newBookID)
	return newBookID, nil
}

func getDbConnection() *sql.DB {
	connStr := "postgres://postgres:example@books-db:5432/postgres?sslmode=disable"
	// connStr := "host=books-db dbname=postgres user=postgres password=example sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Test the database connection
	for retries := 5; retries > 0; retries-- {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	return db
}

func validateBook(book Book) bool {
	if book.Title == "" || book.Author == "" || book.Pages <= 0 {
		return false
	}
	return true
}
