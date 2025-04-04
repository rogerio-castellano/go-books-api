package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/cors"
)

var db *sql.DB

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func init() {
	db = openDbConnection()
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(corsMiddleware)

	r.HandleFunc("/api/books", handleGetBooks).Methods("GET")
	r.HandleFunc("/api/books", handlePostBook).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/books", handlePutBook).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/books/{id}", handleGetBookById).Methods("GET")
	r.HandleFunc("/api/books/{id}", handleDeleteBook).Methods("DELETE", "OPTIONS")

	// Ensure database is closed on program exit
	CloseDatabaseOnProgramExit()

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleGetBooks(w http.ResponseWriter, r *http.Request) {
	//GET http://localhost:8080/books
	//GET http://localhost:8080/books/
	var responseJSON []byte
	var err error

	books, err := getBooks()
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}

	responseJSON, err = json.Marshal(books)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func handleGetBookById(w http.ResponseWriter, r *http.Request) {
	//GET http://localhost:8080/books/1
	var responseJSON []byte
	var err error
	var id int

	idParam := mux.Vars(r)["id"]

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

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
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
	var err error

	// Decode the request body into a map to retrieve only the keys that are present in the request body
	var rawBody map[string]json.RawMessage
	err = json.NewDecoder(r.Body).Decode(&rawBody)
	if err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	normalizeKeysToLowercase := func(rawBody *map[string]json.RawMessage) map[string]json.RawMessage {
		var result = make(map[string]json.RawMessage)
		for k, v := range *rawBody {
			result[strings.ToLower(k)] = v
		}
		return result
	}
	rawBody = normalizeKeysToLowercase(&rawBody)

	var bookToUpdate Book
	bodyBytes, _ := json.Marshal(rawBody)
	err = json.Unmarshal(bodyBytes, &bookToUpdate)
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

	// Update only the fields that are present in the request body, keeping the existing values unchanged for the others
	if _, exists := rawBody["title"]; exists {
		book.Title = bookToUpdate.Title
	}

	if _, exists := rawBody["author"]; exists {
		book.Author = bookToUpdate.Author
	}

	if _, exists := rawBody["pages"]; exists {
		book.Pages = bookToUpdate.Pages
	}

	if !validateBook(book) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = updateBook(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(book)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	//DELETE http://localhost:8080/books/4
	idParam := mux.Vars(r)["id"]

	if idParam == "" {
		http.Error(w, "The id was not provided.", http.StatusBadRequest)
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

	err = deleteBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func corsMiddleware(next http.Handler) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},                   // Add allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Specify allowed methods
		AllowedHeaders:   []string{"Content-Type"},                            // Add custom headers
		AllowCredentials: true,                                                // Allow credentials (like cookies)
	})

	return c.Handler(next)
}

func getBooks() ([]Book, error) {
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
	query := `INSERT INTO books (title, author, pages) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, book.Title, book.Author, book.Pages).Scan(&newBookID)
	if err != nil {
		return 0, err
	}
	return newBookID, nil
}

func updateBook(book Book) error {
	query := `UPDATE books SET title = $1, author = $2, pages = $3 WHERE id = $4`
	_, err := db.Exec(query, book.Title, book.Author, book.Pages, book.Id)

	return err
}

func deleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func openDbConnection() *sql.DB {
	if isDBAvailable(db) {
		return db
	}

	var db *sql.DB
	connStr := "postgres://" + os.Getenv("POSTGRES_USERNAME") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DATABASE") + "?sslmode=disable&connect_timeout=10&application_name=books-api"
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

func isDBAvailable(db *sql.DB) bool {
	if db == nil {
		return false
	}

	// Check if the database connection is still alive
	err := db.Ping()
	if err != nil {
		if err == sql.ErrConnDone {
			return false
		}
		fmt.Printf("An error occurred: %v\n", err)
	}
	return true
}

func validateBook(book Book) bool {
	if book.Title == "" || book.Author == "" || book.Pages <= 0 {
		return false
	}
	return true
}

func CloseDatabaseOnProgramExit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		if db != nil {
			db.Close()
			fmt.Println("Closing database connection...")
		}
		os.Exit(0) // Gracefully exit
	}()
}
