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

// func openMongoDbConnection() {
// 		uri := "mongodb://books-mongo:27017"

// 	client, err := mongo.Connect(options.Client().ApplyURI(uri))
// 	if err != nil {
// 		panic(err)
// 	}

// 	coll := client.Database("test").Collection("books")
// 	filter := bson.D{{Key: "id", Value: 3}}
// 	// Retrieves the first matching document

// 	var result Book
// 	err = coll.FindOne(context.TODO(), filter).Decode(&result)
// 	// err = coll.FindOne(context.TODO(), bson.M{}).Decode(&result)
// 	// Prints a message if no documents are matched or if any
// 	// other errors occur during the operation
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return
// 		}
// 		panic(err)
// 	}
// fmt.Println(result)
// }

// func sendAndReceiveMessageRabbitMq() {
//    	retries := 5               // Number of retry attempts
//     delay := 2 * time.Second   // Delay between attempts
//     conn := connectToRabbitMQ(retries, delay)
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"hello", // name
// 		false,   // durable
// 		false,   // delete when unused
// 		false,   // exclusive
// 		false,   // no-wait
// 		nil,     // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	body := "Hello World! " + time.Now().Format("2006-01-02 15:04:05")
// 	err = ch.PublishWithContext(ctx,
// 		"",     // exchange
// 		q.Name, // routing key
// 		false,  // mandatory
// 		false,  // immediate
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte(body),
// 		})
// 	failOnError(err, "Failed to publish a message")
// 	log.Printf(" [x] Sent %s\n", body)

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	go func() {
// 		for d := range msgs {
// 			log.Printf("Received a message: %s", d.Body)
// 		}
// 	}()
// }

// func connectToRabbitMQ(retries int, delay time.Duration) *amqp.Connection {
//     var conn *amqp.Connection
//     var err error

//     for i := 1; i <= retries; i++ {
//         conn, err = amqp.Dial("amqp://guest:guest@rabbitmq/")
//         if err == nil {
//             fmt.Printf("Successfully connected to RabbitMQ on attempt %d\n", i)
//             return conn
//         }

//         log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v\n", i, retries, err)
//         time.Sleep(delay) // Wait before retrying
//     }

//     log.Fatalf("Exceeded maximum retries (%d). RabbitMQ connection failed.", retries)
//     return nil // This will not be reached due to log.Fatalf, but is included for completeness.
// }

// func setAndGetRedisKey() {
// 	retries := 5             // Number of retry attempts
// 	delay := 2 * time.Second // Delay between attempts

// 	client := connectToRedis(retries, delay)
// 	defer client.Close()

// 	ctx := context.Background()

// 	err := client.Set(ctx, "foo", "bar", 0).Err()
// 	if err != nil {
// 		fmt.Println(10)
// 		panic(err)
// 	}

// 	val, err := client.Get(ctx, "foo").Result()
// 	if err != nil {
// 		fmt.Println(20)
// 		panic(err)
// 	}
// 	fmt.Println("foo", val)
// }

// func connectToRedis(retries int, delay time.Duration) *redis.Client {
//     var client *redis.Client
//     var err error

//     for i := 1; i <= retries; i++ {
//         // Initialize Redis client
//         client = redis.NewClient(&redis.Options{
//             Addr:     "books-redis:6379",
//             Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", // No password set
//             DB:       0,                                  // Use default DB
//             Protocol: 2,                                  // Connection protocol
//         })

//         // Test connection using PING
//         ctx := context.Background()
//         _, err = client.Ping(ctx).Result()

//         if err == nil {
//             fmt.Printf("Successfully connected to Redis on attempt %d\n", i)
//             return client
//         }

//         log.Printf("Failed to connect to Redis (attempt %d/%d): %v\n", i, retries, err)
//         time.Sleep(delay) // Wait before retrying
//     }

//     log.Fatalf("Exceeded maximum retries (%d). Redis connection failed.", retries)
//     return nil // Will not be reached due to log.Fatalf
// }

func main() {
	r := mux.NewRouter()
	r.Use(corsMiddleware)

	r.HandleFunc("/books", handleGetBooks).Methods("GET")
	r.HandleFunc("/books", handlePostBook).Methods("POST", "OPTIONS")
	r.HandleFunc("/books", handlePutBook).Methods("PUT", "OPTIONS")
	r.HandleFunc("/books/", handleGetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", handleGetBookById).Methods("GET")
	r.HandleFunc("/books/{id}", handleDeleteBook).Methods("DELETE", "OPTIONS")

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

	books, _ := getBooks()
	responseJSON, err = json.Marshal(books)

	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseJSON)
}

func handleGetBookById(w http.ResponseWriter, r *http.Request) {
	//GET http://localhost:8080/books/1
	var responseJSON []byte
	var err error
	var id int

	idParam := mux.Vars(r)["id"]

	if idParam == "" {
		http.Error(w, "The id was not provided.", http.StatusBadRequest)
		return
	} else {
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
	w.WriteHeader(http.StatusCreated)

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

	// Wrap the next handler with CORS
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
	if db != nil && !isDBClosed(db) {
		return db
	}

	var db *sql.DB
	connStr := "postgres://postgres:example@books-db:5432/postgres?sslmode=disable&connect_timeout=10&application_name=books-api"
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

func isDBClosed(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		if err == sql.ErrConnDone {
			return true
		}
		fmt.Printf("An error occurred: %v\n", err)
	}
	return false
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

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }
